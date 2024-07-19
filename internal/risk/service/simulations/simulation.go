package simulation

import (
	"math"
	"math/rand"
	"net/http"
	"qira/db"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat/distuv"
	"xorm.io/xorm"
)

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
	var risk []db.RiskCalculation
	var catalog []db.ThreatEventCatalog
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &risk); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := db.GetAll(engine.(*xorm.Engine), &catalog); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	manualData := processRiskData(risk, len(catalog))

	iterations := 10000 // Número de Simulações Monte Carlo

	updatedData, riskArray, frequencyTrack := generateData(manualData, iterations)

	c.JSON(http.StatusOK, gin.H{
		"data":            updatedData,
		"risk_array":      riskArray,
		"frequency_track": frequencyTrack,
	})
}

func processRiskData(risk []db.RiskCalculation, catalogSize int) []RiskData {
	var frequencyEstimates, lossEstimates []float64

	for _, r := range risk {
		if r.RiskType == "Frequency" {
			frequencyEstimates = append(frequencyEstimates, r.Min, r.Max, r.Mode)
		} else if r.RiskType == "Loss" {
			lossEstimates = append(lossEstimates, r.Min, r.Max, r.Mode)
		}
	}

	// Generate slices of log-normal and uniform values for simulation
	meanFrequency := GenerateLogNormalSlice(2, 0.5, catalogSize)
	stdFrequency := GenerateUniformSlice(0.1, 1.0, catalogSize)
	meanLoss := GenerateLogNormalSlice(3, 1, catalogSize)
	stdLoss := GenerateUniformSlice(0.5, 2.0, catalogSize)

	// Aggregate generated slices into single values (simple average used here for demonstration)
	aggregatedMeanFrequency := average(meanFrequency)
	aggregatedStdFrequency := average(stdFrequency)
	aggregatedMeanLoss := average(meanLoss)
	aggregatedStdLoss := average(stdLoss)

	manualData := []RiskData{
		{
			EventName:     "Risk Event",
			MeanFrequency: aggregatedMeanFrequency,
			StdFrequency:  aggregatedStdFrequency,
			MeanLoss:      aggregatedMeanLoss,
			StdLoss:       aggregatedStdLoss,
		},
	}

	return manualData
}

func average(values []float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total / float64(len(values))
}

func generateData(manualData []RiskData, iterations int) ([]RiskData, [][]float64, map[int]int) {
	rand.Seed(time.Now().UnixNano())
	riskArray := make([][]float64, iterations)
	frequencyTrack := make(map[int]int)

	for i := 0; i < iterations; i++ {
		riskArray[i] = make([]float64, len(manualData))
		for j, data := range manualData {
			freqDist := distuv.LogNormal{Mu: math.Log(data.MeanFrequency), Sigma: data.StdFrequency}
			lossDist := distuv.LogNormal{Mu: math.Log(data.MeanLoss), Sigma: data.StdLoss}
			frequencySample := freqDist.Rand()
			lossSample := lossDist.Rand()
			riskArray[i][j] = frequencySample * lossSample
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

		for _, risk := range risks {
			roundedRisk := int(math.Round(risk))
			frequencyTrack[roundedRisk]++
		}
	}

	return manualData, riskArray, frequencyTrack
}

func mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range data {
		total += value
	}
	return total / float64(len(data))
}

func percentile(data []float64, percent float64) float64 {
	if len(data) == 0 {
		return 0
	}
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
