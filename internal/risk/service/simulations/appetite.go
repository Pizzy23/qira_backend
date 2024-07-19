package simulation

import (
	"math"
	"net/http"
	"qira/db"
	"sort"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
	"xorm.io/xorm"
)

type RiskAppetiteData struct {
	EventID               int64       `json:"event_id"`
	EventName             string      `json:"event_name"`
	CurrentControl        float64     `json:"current_control"`
	ProposedControl       float64     `json:"proposed_control"`
	AcceptableLosses      []LossLevel `json:"acceptable_losses"`
	ProbabilityThresholds map[float64]float64
}

type LossLevel struct {
	Probability    float64 `json:"probability"`
	AcceptableLoss float64 `json:"acceptable_loss"`
}

func FetchRiskAppetite() (RiskAppetiteData, error) {
	return RiskAppetiteData{
		EventID:         1,
		EventName:       "Example Event",
		CurrentControl:  15.0,
		ProposedControl: 25.0,
		AcceptableLosses: []LossLevel{
			{Probability: 100, AcceptableLoss: 5000000},
			{Probability: 75, AcceptableLoss: 3750000},
			{Probability: 50, AcceptableLoss: 2500000},
			{Probability: 25, AcceptableLoss: 1250000},
			{Probability: 0, AcceptableLoss: 625000},
		},
		ProbabilityThresholds: map[float64]float64{
			100: 5000000,
			75:  3750000,
			50:  2500000,
			25:  1250000,
			0:   625000,
		},
	}, nil
}

func MonteCarloSimulationAppetite(c *gin.Context) {
	engine := c.MustGet("db").(*xorm.Engine)

	riskAppetite, err := FetchRiskAppetite()
	if err != nil {
		c.Set("Response", "Failed to fetch risk calculations")
		c.Status(http.StatusInternalServerError)
		return
	}

	var riskCalculations []db.RiskCalculation
	if err := engine.Find(&riskCalculations); err != nil {
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
		result := generateRiskDataAppetite(eventData, sims, riskAppetite)
		results[calc.ThreatEvent] = result
	}

	c.Set("Response", results)
	c.Status(http.StatusOK)
}

func generateRiskDataAppetite(event EventData, iterations int, appetite RiskAppetiteData) RiskAnalysisResults {
	riskSamples := make([]float64, iterations)
	frequencyMap := make(map[int]int)

	for i := 0; i < iterations; i++ {
		freqSample := PERTLogNormal(event.MinFrequency, event.PertFrequency, event.MaxFrequency)
		lossSample := PERTLogNormal(event.MinLoss, event.PertLoss, event.MaxLoss)
		riskResult := freqSample * lossSample

		// Apply risk appetite filter
		if isAcceptableRisk(riskResult, appetite) {
			roundedRiskResult := int(math.Round(riskResult/100000)) * 100000
			frequencyMap[roundedRiskResult]++
		}
	}

	sort.Float64s(riskSamples)
	return RiskAnalysisResults{
		Event:        event.Event,
		AverageRisk:  stat.Mean(riskSamples, nil),
		P95Risk:      stat.Quantile(0.95, stat.Empirical, riskSamples, nil),
		ValueAtRisk:  stat.Quantile(0.99, stat.Empirical, riskSamples, nil),
		FrequencyMap: frequencyMap,
	}
}

func isAcceptableRisk(loss float64, appetite RiskAppetiteData) bool {
	for _, level := range appetite.AcceptableLosses {
		if loss <= level.AcceptableLoss {
			return true
		}
	}
	return false
}
