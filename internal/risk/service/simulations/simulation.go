package simulation

import (
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
	rand.Seed(time.Now().UnixNano())
	logPert := math.Log(pert)
	sigma := (math.Log(max) - math.Log(min)) / 6
	mu := logPert - (sigma * sigma / 2)

	dist := distuv.LogNormal{
		Mu:    mu,
		Sigma: sigma,
	}
	return dist.Rand()
}

func MonteCarloSimulation(c *gin.Context) {
	engine := c.MustGet("db").(*xorm.Engine)

	var riskCalculations []db.RiskCalculation
	err := engine.Find(&riskCalculations)
	if err != nil {
		c.Set("Response", "Failed to fetch risk calculations")
		c.Status(http.StatusInternalServerError)
		return
	}

	results := make(map[string]RiskAnalysisResults)
	for _, calc := range riskCalculations {
		eventData := EventData{
			Event:        calc.ThreatEvent,
			MinFrequency: calc.Min, PertFrequency: calc.Mode, MaxFrequency: calc.Max,
			MinLoss: calc.Min, PertLoss: calc.Mode, MaxLoss: calc.Max,
		}
		results[calc.ThreatEvent] = generateRiskData(eventData, sims)
	}

	c.Set("Response", results)
	c.Status(http.StatusOK)
}

func generateRiskData(event EventData, iterations int) RiskAnalysisResults {
	riskSamples := make([]float64, iterations)
	frequencyMap := make(map[int]int)

	for i := range riskSamples {
		freqSample := PERTLogNormal(event.MinFrequency, event.PertFrequency, event.MaxFrequency)
		lossSample := PERTLogNormal(event.MinLoss, event.PertLoss, event.MaxLoss)
		riskResult := freqSample * lossSample
		riskSamples[i] = riskResult
		// Arredondar para a centena de milhar mais próxima para menos granularidade
		roundedRiskResult := int(math.Round(riskResult/100000)) * 100000
		frequencyMap[roundedRiskResult]++
	}

	sort.Float64s(riskSamples)

	meanRisk := stat.Mean(riskSamples, nil)
	p95Risk := stat.Quantile(0.95, stat.Empirical, riskSamples, nil)
	varRisk := stat.Quantile(0.99, stat.Empirical, riskSamples, nil)

	return RiskAnalysisResults{
		Event:        event.Event,
		AverageRisk:  meanRisk,
		P95Risk:      p95Risk,
		ValueAtRisk:  varRisk,
		Error:        p95Risk - meanRisk,
		FrequencyMap: frequencyMap,
	}
}
