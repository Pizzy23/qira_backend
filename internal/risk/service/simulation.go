package risk

import (
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

type RiskData struct {
	EventName     string
	MeanFrequency float64
	StdFrequency  float64
	MeanLoss      float64
	StdLoss       float64
	MeanRisk      float64
	Percentile95  float64
	ValueAtRisk   float64
	Error         float64
}

type CombinedRiskLEC struct {
	EventName    string         `json:"event_name"`
	MeanRisk     float64        `json:"mean_risk"`
	Percentile95 float64        `json:"percentile_95"`
	ValueAtRisk  float64        `json:"value_at_risk"`
	RiskLEC      []RiskLECPoint `json:"risk_lec"`
}

type RiskLECPoint struct {
	Risk float64 `json:"risk"`
	LEC  float64 `json:"lec"`
}

func generateData(iterations int, events []db.ThreatEventCatalog, freqMap map[int64]db.Frequency, lossMap map[int64]db.LossHigh) ([]CombinedRiskLEC, [][]float64) {
	rand.Seed(time.Now().UnixNano())
	var combinedRiskLEC []CombinedRiskLEC

	for _, event := range events {
		frequency := freqMap[event.ID]
		loss := lossMap[event.ID]

		meanFreq := rand.NormFloat64()*0.5 + frequency.MostLikelyFrequency
		stdFreq := (frequency.MaxFrequency - frequency.MinFrequency) / 6 // Assuming standard deviation is 1/6th of the range
		meanLoss := rand.NormFloat64()*1 + loss.MostLikelyLoss
		stdLoss := (loss.MaximumLoss - loss.MinimumLoss) / 6 // Assuming standard deviation is 1/6th of the range

		var eventRisk []float64
		for i := 0; i < iterations; i++ {
			freqSample := distuv.LogNormal{Mu: meanFreq, Sigma: stdFreq}.Rand()
			lossSample := distuv.LogNormal{Mu: meanLoss, Sigma: stdLoss}.Rand()
			eventRisk = append(eventRisk, freqSample*lossSample)
		}

		sort.Float64s(eventRisk)

		meanRisk := stat.Mean(eventRisk, nil)
		percentile95 := stat.Quantile(0.95, stat.Empirical, eventRisk, nil)
		valueAtRisk := stat.Quantile(0.99, stat.Empirical, eventRisk, nil)

		probability := make([]float64, len(eventRisk))
		for i := range eventRisk {
			probability[i] = 1.0 - float64(i)/float64(len(eventRisk))
		}

		var riskLEC []RiskLECPoint
		for i := range eventRisk {
			riskLEC = append(riskLEC, RiskLECPoint{
				Risk: eventRisk[i],
				LEC:  probability[i],
			})
		}

		combinedRiskLEC = append(combinedRiskLEC, CombinedRiskLEC{
			EventName:    event.ThreatEvent,
			MeanRisk:     meanRisk,
			Percentile95: percentile95,
			ValueAtRisk:  valueAtRisk,
			RiskLEC:      riskLEC,
		})
	}

	return combinedRiskLEC, nil
}

func MainSimulation(c *gin.Context) {
	var events []db.ThreatEventCatalog
	var loss []db.LossHigh
	var freq []db.Frequency
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

	if err := db.GetAll(engine.(*xorm.Engine), &loss); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &freq); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	iterations := 50

	freqMap := make(map[int64]db.Frequency)
	for _, f := range freq {
		freqMap[f.ThreatEventID] = f
	}

	lossMap := make(map[int64]db.LossHigh)
	for _, l := range loss {
		lossMap[l.ThreatEventID] = l
	}

	combinedRiskLEC, _ := generateData(iterations, events, freqMap, lossMap)

	c.Set("Response", combinedRiskLEC)
	c.Status(http.StatusOK)

}
