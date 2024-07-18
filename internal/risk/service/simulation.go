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
	EventName    string    `json:"event_name"`
	MeanRisk     float64   `json:"mean_risk"`
	Percentile95 float64   `json:"percentile_95"`
	ValueAtRisk  float64   `json:"value_at_risk"`
	RiskSims     []float64 `json:"risk_sims"`
	RiskBins     []float64 `json:"risk_bins"`
	RiskFreqs    []float64 `json:"risk_freqs"`
	RiskLecs     []float64 `json:"risk_lecs"`
	RiskCFreqs   []float64 `json:"risk_cfreqs"`
}

func generateData(iterations int, events []db.ThreatEventCatalog, frequencies []db.Frequency, losses []db.LossHigh) []CombinedRiskLEC {
	rand.Seed(time.Now().UnixNano())
	var combinedRiskLEC []CombinedRiskLEC

	for _, event := range events {
		var frequency db.Frequency
		var loss db.LossHigh

		for _, f := range frequencies {
			if f.ThreatEventID == event.ID {
				frequency = f
				break
			}
		}

		for _, l := range losses {
			if l.ThreatEventID == event.ID {
				loss = l
				break
			}
		}

		meanFreq := (frequency.MinFrequency + frequency.MostLikelyFrequency + frequency.MaxFrequency) / 3
		stdFreq := (frequency.MaxFrequency - frequency.MinFrequency) / 6
		meanLoss := (loss.MinimumLoss + loss.MostLikelyLoss + loss.MaximumLoss) / 3
		stdLoss := (loss.MaximumLoss - loss.MinimumLoss) / 6

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

		riskFreqs := make([]float64, len(eventRisk))
		for i := range eventRisk {
			riskFreqs[i] = eventRisk[i] * probability[i]
		}

		riskBins := make([]float64, len(eventRisk))
		for i := range eventRisk {
			riskBins[i] = eventRisk[i]
		}

		riskLecs := make([]float64, len(eventRisk))
		for i := range eventRisk {
			riskLecs[i] = probability[i]
		}

		riskCFreqs := make([]float64, len(eventRisk))
		for i := range eventRisk {
			riskCFreqs[i] = 1.0 - probability[i]
		}

		combinedRiskLEC = append(combinedRiskLEC, CombinedRiskLEC{
			EventName:    event.ThreatEvent,
			MeanRisk:     meanRisk,
			Percentile95: percentile95,
			ValueAtRisk:  valueAtRisk,
			RiskSims:     eventRisk,
			RiskBins:     riskBins,
			RiskFreqs:    riskFreqs,
			RiskLecs:     riskLecs,
			RiskCFreqs:   riskCFreqs,
		})
	}

	return combinedRiskLEC
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

	iterations := 10000

	combinedRiskLEC := generateData(iterations, events, freq, loss)

	c.Set("Response", combinedRiskLEC)
	c.Status(http.StatusOK)
}
