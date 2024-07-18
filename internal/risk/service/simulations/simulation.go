package simulation

import (
	"math/rand"
	"net/http"
	"qira/db"
	"time"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat/distuv"
	"xorm.io/xorm"
)

type MonteCarloResult struct {
	Simulation int     `json:"simulation"`
	Result     float64 `json:"result"`
}

func MonteCarloSimulation(c *gin.Context) {
	var events []db.ThreatEventCatalog
	engine, exists := c.Get("db")

	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &events); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	iterations := 10000
	predefinedValue := 19.55185417
	results := generateMonteCarloRiskTest(iterations, events, predefinedValue)

	c.Set("Response", results)
	c.Status(http.StatusOK)
}

func generateMonteCarloRiskTest(iterations int, events []db.ThreatEventCatalog, predefinedValue float64) []MonteCarloResult {
	rand.Seed(time.Now().UnixNano())
	var results []MonteCarloResult

	for i := 0; i < iterations; i++ {
		var totalAggregatedRisk float64
		for range events {
			riskSample := distuv.LogNormal{Mu: predefinedValue, Sigma: 1}.Rand()
			totalAggregatedRisk += riskSample
		}

		result := MonteCarloResult{
			Simulation: i + 1,
			Result:     totalAggregatedRisk,
		}

		results = append(results, result)
	}

	return results
}
