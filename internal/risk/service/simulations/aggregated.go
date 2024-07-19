package simulation

import (
	"net/http"
	"qira/db"
	"sort"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
	"xorm.io/xorm"
)

const maxBuckets = 100 // Define a maximum number of buckets

func MonteCarloSimulationAggregated(c *gin.Context) {
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
		results[calc.ThreatEvent] = generateRiskDataAggregated(eventData, sims)
	}
	c.JSON(http.StatusOK, results)
}

func generateRiskDataAggregated(event EventData, iterations int) RiskAnalysisResults {
	riskSamples := make([]float64, iterations)
	for i := range riskSamples {
		freqSample := PERTLogNormal(event.MinFrequency, event.PertFrequency, event.MaxFrequency)
		lossSample := PERTLogNormal(event.MinLoss, event.PertLoss, event.MaxLoss)
		riskSamples[i] = freqSample * lossSample
	}

	sort.Float64s(riskSamples)
	frequencyMap := aggregateRiskResults(riskSamples)

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

func aggregateRiskResults(riskSamples []float64) map[int]int {
	frequencyMap := make(map[int]int)
	minRisk := riskSamples[0]
	maxRisk := riskSamples[len(riskSamples)-1]
	bucketSize := (maxRisk - minRisk) / float64(maxBuckets)

	for _, risk := range riskSamples {
		bucket := int(risk/bucketSize) * int(bucketSize)
		frequencyMap[bucket]++
	}
	return frequencyMap
}
