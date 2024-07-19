package simulation

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"qira/db"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
	"xorm.io/xorm"
)

const sims = 10000

type EventData struct {
	Event                                     string
	MinFrequency, PertFrequency, MaxFrequency float64
	MinLoss, PertLoss, MaxLoss                float64
}

type RiskAnalysisResults struct {
	Event        string
	AverageRisk  float64
	P95Risk      float64
	ValueAtRisk  float64
	Error        float64
	FrequencyMap map[int]int
}

func PERTLogNormal(min, pert, max float64) float64 {

	if min <= 0 || pert <= 0 || max <= 0 || min >= max || pert >= max {
		return 0
	}

	rand.Seed(time.Now().UnixNano())
	logPert := math.Log(pert)
	sigma := (math.Log(max) - math.Log(min)) / 6
	mu := logPert - (sigma * sigma / 2)

	if math.IsNaN(mu) || math.IsNaN(sigma) {
		log.Printf("Invalid LogNormal parameters: mu=%f, sigma=%f", mu, sigma)
		return 0
	}

	dist := distuv.LogNormal{
		Mu:    mu,
		Sigma: sigma,
	}
	return dist.Rand()
}

func MonteCarloSimulation(c *gin.Context, threatEvent string) {
	engine := c.MustGet("db").(*xorm.Engine)

	// Get the threat event from the query parameter

	var riskCalculations []db.RiskCalculation
	err := engine.Where("threat_event = ?", threatEvent).Find(&riskCalculations)
	if err != nil {
		log.Println("Failed to fetch risk calculations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch risk calculations"})
		return
	}

	if len(riskCalculations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Threat event not found"})
		return
	}

	results := make(map[string]RiskAnalysisResults)
	for _, calc := range riskCalculations {
		eventData := EventData{
			Event:         calc.ThreatEvent,
			MinFrequency:  calc.Min,
			PertFrequency: calc.Mode,
			MaxFrequency:  calc.Max,
			MinLoss:       calc.Min,
			PertLoss:      calc.Mode,
			MaxLoss:       calc.Max,
		}
		results[calc.ThreatEvent] = generateRiskData(eventData, sims)
		break
	}

	c.JSON(http.StatusOK, results)
}

func generateRiskData(event EventData, iterations int) RiskAnalysisResults {
	riskSamples := make([]float64, iterations)
	frequencyMap := make(map[int]int)

	for i := range riskSamples {
		freqSample := PERTLogNormal(event.MinFrequency, event.PertFrequency, event.MaxFrequency)
		lossSample := PERTLogNormal(event.MinLoss, event.PertLoss, event.MaxLoss)

		if math.IsNaN(freqSample) || math.IsNaN(lossSample) {
			log.Printf("NaN value detected: freqSample=%f, lossSample=%f", freqSample, lossSample)
			riskSamples[i] = 0
		} else {
			riskResult := freqSample * lossSample
			riskSamples[i] = riskResult
			roundedRiskResult := int(math.Round(riskResult/100000)) * 100000
			frequencyMap[roundedRiskResult]++
		}
	}

	sort.Float64s(riskSamples)

	meanRisk := stat.Mean(riskSamples, nil)
	p95Risk := stat.Quantile(0.95, stat.Empirical, riskSamples, nil)
	varRisk := stat.Quantile(0.99, stat.Empirical, riskSamples, nil)

	if math.IsNaN(meanRisk) {
		meanRisk = 0
	}
	if math.IsNaN(p95Risk) {
		p95Risk = 0
	}
	if math.IsNaN(varRisk) {
		varRisk = 0
	}

	return RiskAnalysisResults{
		Event:        event.Event,
		AverageRisk:  meanRisk,
		P95Risk:      p95Risk,
		ValueAtRisk:  varRisk,
		Error:        p95Risk - meanRisk,
		FrequencyMap: frequencyMap,
	}
}
