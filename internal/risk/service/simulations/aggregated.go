package simulation

import (
	"log"
	"math/rand"
	"net/http"
	"qira/db"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat/distuv"
	"xorm.io/xorm"
)

type AggregatedRiskResult struct {
	Simulation     int            `json:"simulation"`
	AggregatedRisk float64        `json:"aggregated_risk"`
	Frequency      map[string]int `json:"frequency"` // Change to map[string]int
	Cumulative     []float64      `json:"cumulative"`
}

func MonteCarloSimulationAggregated(c *gin.Context) {
	var riskCalculations []db.RiskCalculation
	engine, exists := c.Get("db")

	if !exists {
		log.Println("Database connection not found")
		c.JSON(http.StatusInternalServerError, gin.H{"Response": "Database connection not found"})
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &riskCalculations); err != nil {
		log.Println("Error fetching risk calculations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Response": err.Error()})
		return
	}

	// Calcula o valor predefinido somando todos os campos Estimate dos eventos de risco
	predefinedValue := 0.0
	for _, risk := range riskCalculations {
		if risk.RiskType == "Risk" {
			predefinedValue += risk.Estimate
		}
	}
	log.Println("Predefined value:", predefinedValue)

	iterations := 10
	results := generateMonteCarloRisk(iterations, predefinedValue)

	c.JSON(http.StatusOK, gin.H{"Response": results})
}

func generateMonteCarloRisk(iterations int, predefinedValue float64) []AggregatedRiskResult {
	rand.Seed(time.Now().UnixNano())
	var results []AggregatedRiskResult
	frequencyMap := make(map[float64]int)

	for i := 0; i < iterations; i++ {
		var totalAggregatedRisk float64

		// Usando o valor predefinido como a média (Mu) da distribuição log-normal
		riskSample := distuv.LogNormal{Mu: predefinedValue, Sigma: 1}.Rand()
		totalAggregatedRisk += riskSample

		// Armazenar o risco agregado total para cada simulação
		simulationResult := AggregatedRiskResult{
			Simulation:     i + 1,
			AggregatedRisk: totalAggregatedRisk,
		}

		frequencyMap[totalAggregatedRisk]++
		results = append(results, simulationResult)
	}

	log.Println("Simulations completed. Processing frequencies and cumulative data.")

	// Calculate cumulative distribution
	var freqKeys []float64
	for k := range frequencyMap {
		freqKeys = append(freqKeys, k)
	}
	sort.Float64s(freqKeys)

	cumulative := make([]float64, len(freqKeys))
	totalSimulations := float64(len(results))
	for i, key := range freqKeys {
		if i == 0 {
			cumulative[i] = float64(frequencyMap[key]) / totalSimulations * 100
		} else {
			cumulative[i] = cumulative[i-1] + float64(frequencyMap[key])/totalSimulations*100
		}
	}

	log.Println("Frequencies and cumulative data processed. Assigning to results.")

	// Atribuir frequências e cumulativas corretas para cada resultado
	for i := range results {
		results[i].Cumulative = cumulative
		stringFrequencyMap := make(map[string]int)
		for k, v := range frequencyMap {
			stringFrequencyMap[formatFloat(k)] = v
		}
		results[i].Frequency = stringFrequencyMap
	}

	log.Println("Results prepared for response.")
	return results
}

func formatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 3, 64)
}
