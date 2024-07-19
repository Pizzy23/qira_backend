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

const sims = 10000
const stdev = 3.29

type RiskData struct {
	EventName     string  `json:"event_name"`
	MinFreq       float64 `json:"min_freq"`
	PertFreq      float64 `json:"pert_freq"`
	MaxFreq       float64 `json:"max_freq"`
	MinLoss       float64 `json:"min_loss"`
	PertLoss      float64 `json:"pert_loss"`
	MaxLoss       float64 `json:"max_loss"`
	MeanRisk      float64 `json:"mean_risk"`
	Percentile95  float64 `json:"percentile_95"`
	ValueAtRisk   float64 `json:"value_at_risk"`
	Error         float64 `json:"error"`
	MeanFrequency float64 `json:"mean_frequency"`
	StdFrequency  float64 `json:"std_frequency"`
	MeanLoss      float64 `json:"mean_loss"`
	StdLoss       float64 `json:"std_loss"`
}

func MonteCarloSimulation(c *gin.Context) {
	var risk []db.RiskCalculation
	var event []db.ThreatEventCatalog
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

	manualData := processRiskData(risk, event)
	iterations := sims // Número de Simulações Monte Carlo

	updatedData, riskArray, frequencyTrack := generateData(manualData, iterations)

	c.JSON(http.StatusOK, gin.H{
		"data":            updatedData,
		"risk_array":      riskArray,
		"frequency_track": frequencyTrack,
	})
}

func processRiskData(risks []db.RiskCalculation, events []db.ThreatEventCatalog) []RiskData {
	eventSize := len(events)
	if eventSize == 0 {
		return nil
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64
	var totalEstimativeLoss, totalEstimativeFreq float64

	for _, risk := range risks {
		if risk.RiskType == "Frequency" {
			totalMinFreq += risk.Min
			totalPertFreq += risk.Mode
			totalMaxFreq += risk.Max
			totalEstimativeFreq += risk.Estimate
		} else if risk.RiskType == "Loss" {
			totalMinLoss += risk.Min
			totalPertLoss += risk.Mode
			totalMaxLoss += risk.Max
			totalEstimativeLoss += risk.Estimate
		}
	}
	riskSize := len(risks)
	meanMinFreq := totalMinFreq / float64(riskSize)
	meanPertFreq := totalPertFreq / float64(riskSize)
	meanMaxFreq := totalMaxFreq / float64(riskSize)
	meanMinLoss := totalMinLoss / float64(riskSize)
	meanPertLoss := totalPertLoss / float64(riskSize)
	meanMaxLoss := totalMaxLoss / float64(riskSize)

	meanFrequency := (meanMinFreq + meanMaxFreq) / 2
	stdFrequency := calculateStdDevFromMinMax(meanMinFreq, meanMaxFreq)
	meanLoss := (meanMinLoss + meanMaxLoss) / 2
	stdLoss := calculateStdDevFromMinMax(meanMinLoss, meanMaxLoss)

	manualData := make([]RiskData, eventSize)
	for i := 0; i < eventSize; i++ {
		manualData[i] = RiskData{
			EventName:     events[i].ThreatEvent,
			MinFreq:       meanMinFreq,
			PertFreq:      meanPertFreq,
			MaxFreq:       meanMaxFreq,
			MinLoss:       meanMinLoss,
			PertLoss:      meanPertLoss,
			MaxLoss:       meanMaxLoss,
			MeanFrequency: meanFrequency,
			StdFrequency:  stdFrequency,
			MeanLoss:      meanLoss,
			StdLoss:       stdLoss,
		}
	}

	return manualData
}

func calculateStdDevFromMinMax(min, max float64) float64 {
	mean := (min + max) / 2
	variance := ((min-mean)*(min-mean) + (max-mean)*(max-mean)) / 2
	return math.Sqrt(variance)
}

func lognorminvpert(min, pert, max float64) float64 {
	mu := math.Log(pert)
	sigma := (math.Log(max) - math.Log(min)) / stdev
	return math.Exp(distuv.LogNormal{Mu: mu, Sigma: sigma}.Rand())
}

func lognormRiskPert(minfreq, pertfreq, maxfreq, minloss, pertloss, maxloss float64) float64 {
	freq := lognorminvpert(minfreq, pertfreq, maxfreq)
	loss := lognorminvpert(minloss, pertloss, maxloss)
	return freq * loss
}

func generateSimData(rdata RiskData, totalTE, teNo int) []float64 {
	simData := make([]float64, sims)
	totalSims := sims * totalTE
	cumulativeTotal := sims * teNo

	for simCtr := 0; simCtr < sims; simCtr++ {
		simData[simCtr] = lognormRiskPert(rdata.MinFreq, rdata.PertFreq, rdata.MaxFreq, rdata.MinLoss, rdata.PertLoss, rdata.MaxLoss)
		if simCtr%1000 == 0 {
			pct := float64(simCtr+cumulativeTotal) / float64(totalSims)
			progressBar(pct)
		}
	}
	return simData
}

func generateData(manualData []RiskData, iterations int) ([]RiskData, [][]float64, map[int]int) {
	rand.Seed(time.Now().UnixNano())
	riskArray := make([][]float64, iterations)
	frequencyTrack := make(map[int]int)

	for i := 0; i < iterations; i++ {
		riskArray[i] = make([]float64, len(manualData))
		for j, data := range manualData {
			riskArray[i][j] = lognormRiskPert(data.MinFreq, data.PertFreq, data.MaxFreq, data.MinLoss, data.PertLoss, data.MaxLoss)
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

func progressBar(pct float64) {
	// Implemente a lógica do progressBar conforme necessário
}
