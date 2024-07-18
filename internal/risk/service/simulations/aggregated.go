package simulation

import (
	"math/rand"
	"net/http"
	"qira/db"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat/distuv"
	"xorm.io/xorm"
)

type AggregatedRiskResult struct {
	Simulation     int
	EventRisks     map[string]float64
	AggregatedRisk float64
	Frequency      []int     `json:"frequency"`
	Cumulative     []float64 `json:"cumulative"`
}
type MonteCarlo struct {
	Simulation     int
	EventRisks     map[string]float64
	AggregatedRisk float64
}

func AggregatedRisk(c *gin.Context) {
	var events []db.ThreatEventCatalog
	var riskCalculations []db.RiskCalculation
	engine, exists := c.Get("db")

	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &events); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &riskCalculations); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	iterations := 10000
	results := generateAggregatedRisk(iterations, events, riskCalculations)

	c.Set("Response", results)
	c.Status(http.StatusOK)
}

func generateAggregatedRisk(iterations int, events []db.ThreatEventCatalog, riskCalculations []db.RiskCalculation) []AggregatedRiskResult {
	rand.Seed(time.Now().UnixNano())
	var results []AggregatedRiskResult
	eventRiskMap := make(map[string][]float64)

	for i := 0; i < iterations; i++ {
		simulationResult := AggregatedRiskResult{
			Simulation: i + 1,
			EventRisks: make(map[string]float64),
		}

		var totalAggregatedRisk float64

		for _, event := range events {
			riskCalculation := getRiskCalculation(event.ID, riskCalculations)
			if riskCalculation == nil {
				continue
			}

			riskSample := distuv.LogNormal{Mu: riskCalculation.Estimate, Sigma: 1}.Rand() // Usando a média (estimate) para o cálculo
			simulationResult.EventRisks[event.ThreatEvent] = riskSample
			totalAggregatedRisk += riskSample

			eventRiskMap[event.ThreatEvent] = append(eventRiskMap[event.ThreatEvent], riskSample)
		}

		simulationResult.AggregatedRisk = totalAggregatedRisk
		results = append(results, simulationResult)
	}

	// Calculate frequency and cumulative distribution
	for _, risks := range eventRiskMap {
		sort.Float64s(risks)
		frequency := make([]int, iterations)
		cumulative := make([]float64, iterations)
		for i, risk := range risks {
			frequency[i] = int(risk)
			cumulative[i] = float64(i+1) / float64(iterations) * 100
		}
		for i := range results {
			results[i].Frequency = frequency
			results[i].Cumulative = cumulative
		}
	}

	return results
}

func getRiskCalculation(eventID int64, riskCalculations []db.RiskCalculation) *db.RiskCalculation {
	for _, rc := range riskCalculations {
		if rc.ThreatEventID == eventID && rc.RiskType == "Risk" {
			return &rc
		}
	}
	return nil
}
