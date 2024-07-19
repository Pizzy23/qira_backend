package simulation

import (
	"math"
	"math/rand"
	"net/http"
	"qira/db"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

const sims = 10000

type RiskData struct {
	EventName     string  `json:"event_name"`
	MeanFrequency float64 `json:"mean_frequency"`
	StdFrequency  float64 `json:"std_frequency"`
	MeanLoss      float64 `json:"mean_loss"`
	StdLoss       float64 `json:"std_loss"`
	MeanRisk      float64 `json:"mean_risk"`
	Percentile95  float64 `json:"percentile_95"`
	ValueAtRisk   float64 `json:"value_at_risk"`
	Error         float64 `json:"error"`
}

func MonteCarloSimulation(c *gin.Context) {
	var riskCalculations []db.RiskCalculation
	var events []db.ThreatEventCatalog
	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &riskCalculations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.GetAll(engine.(*xorm.Engine), &events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Processar os dados de risco para calcular médias e desvios padrão
	manualData := processRiskCalculations(riskCalculations)
	iterations := 10000 // Número de Simulações Monte Carlo

	updatedData, riskArray, frequencyTrack := generateData(manualData, iterations)

	c.JSON(http.StatusOK, gin.H{
		"data":            updatedData,
		"risk_array":      riskArray,
		"frequency_track": frequencyTrack,
	})
}

func processRiskCalculations(riskCalculations []db.RiskCalculation) []RiskData {
	eventMap := make(map[int64]*RiskData)
	for _, calc := range riskCalculations {
		if _, exists := eventMap[calc.ThreatEventID]; !exists {
			eventMap[calc.ThreatEventID] = &RiskData{
				EventName: calc.ThreatEvent,
			}
		}
		event := eventMap[calc.ThreatEventID]
		if calc.RiskType == "Frequency" {
			event.MeanFrequency = (calc.Min + 4*calc.Mode + calc.Max) / 6
			event.StdFrequency = (calc.Max - calc.Min) / 6
		} else if calc.RiskType == "Loss" {
			event.MeanLoss = (calc.Min + 4*calc.Mode + calc.Max) / 6
			event.StdLoss = (calc.Max - calc.Min) / 6
		}
	}
	var manualData []RiskData
	for _, data := range eventMap {
		manualData = append(manualData, *data)
	}
	return manualData
}

func generateData(manualData []RiskData, iterations int) ([]RiskData, [][]float64, map[int]int) {
	rand.Seed(time.Now().UnixNano())
	riskArray := make([][]float64, iterations)
	frequencyTrack := make(map[int]int)

	for i := 0; i < iterations; i++ {
		riskArray[i] = make([]float64, len(manualData))
		for j, data := range manualData {
			freqSamples := GenerateLogNormalSlice(math.Log(data.MeanFrequency), data.StdFrequency, 1)
			lossSamples := GenerateLogNormalSlice(math.Log(data.MeanLoss), data.StdLoss, 1)
			riskArray[i][j] = freqSamples[0] * lossSamples[0]
		}
	}

	for j := range manualData {
		risks := make([]float64, iterations)
		for i := 0; i < iterations; i++ {
			risks[i] = riskArray[i][j]
		}
		meanRisk := mean(risks)
		percentile95 := percentile(risks, 95)
		valueAtRisk := percentile(risks, 99)
		error := percentile95 - meanRisk

		manualData[j].MeanRisk = meanRisk
		manualData[j].Percentile95 = percentile95
		manualData[j].ValueAtRisk = valueAtRisk
		manualData[j].Error = error

		// Rastreamento de frequência
		for _, risk := range risks {
			roundedRisk := int(math.Round(risk))
			frequencyTrack[roundedRisk]++
		}
	}

	return manualData, riskArray, frequencyTrack
}

func mean(data []float64) float64 {
	total := 0.0
	for _, value := range data {
		total += value
	}
	return total / float64(len(data))
}

func percentile(data []float64, percent float64) float64 {
	sort.Float64s(data)
	k := int(math.Ceil(percent/100*float64(len(data)))) - 1
	if k < 0 {
		k = 0
	}
	return data[k]
}

func GenerateUniformSlice(min, max float64, size int) []float64 {
	slice := make([]float64, size)
	for i := range slice {
		slice[i] = min + rand.Float64()*(max-min)
	}
	return slice
}

func GenerateLogNormalSlice(mean, sigma float64, size int) []float64 {
	slice := make([]float64, size)
	for i := range slice {
		normal := rand.NormFloat64()*sigma + mean
		slice[i] = math.Exp(normal)
	}
	return slice
}
