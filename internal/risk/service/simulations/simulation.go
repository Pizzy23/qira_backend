package simulation

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type ThreatEventRequestT struct {
	MinFreq  int64 `json:"minfreq,omitempty"`
	PertFreq int64 `json:"pertfreq,omitempty"`
	MaxFreq  int64 `json:"maxfreq,omitempty"`
	MinLoss  int64 `json:"minloss,omitempty"`
	PertLoss int64 `json:"pertloss,omitempty"`
	MaxLoss  int64 `json:"maxloss,omitempty"`
}

type FrontEndResponseT struct {
	FrequencyMax      int64 `json:"FrequencyMax"`
	FrequencyMin      int64 `json:"FrequencyMin"`
	FrequencyEstimate int64 `json:"FrequencyEstimate"`
	LossMax           int64 `json:"LossMax"`
	LossMin           int64 `json:"LossMin"`
	LossEstimate      int64 `json:"LossEstimate"`
}

func MonteCarloSimulation(c *gin.Context, threatEvent string) {
	var riskCalculations []db.RiskCalculation

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&riskCalculations)
	if err != nil {
		c.Set("Response", "Error retrieving risk calculations")
		c.Status(http.StatusInternalServerError)
		return
	}

	frequencyRequests := []ThreatEventRequestT{}
	lossRequests := []ThreatEventRequestT{}

	for _, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			minFreq := int64(risk.Min)
			pertFreq := int64(risk.Estimate)
			maxFreq := int64(risk.Max)
			frequencyRequests = append(frequencyRequests, ThreatEventRequestT{
				MinFreq:  minFreq,
				PertFreq: pertFreq,
				MaxFreq:  maxFreq,
			})
		} else if risk.RiskType == "Loss" {
			minLoss := int64(risk.Min)
			pertLoss := int64(risk.Estimate)
			maxLoss := int64(risk.Max)
			lossRequests = append(lossRequests, ThreatEventRequestT{
				MinLoss:  minLoss,
				PertLoss: pertLoss,
				MaxLoss:  maxLoss,
			})
		}
	}

	threatEventRequests := []ThreatEventRequestT{}
	for i := range frequencyRequests {
		te := ThreatEventRequestT{
			MinFreq:  frequencyRequests[i].MinFreq,
			PertFreq: frequencyRequests[i].PertFreq,
			MaxFreq:  frequencyRequests[i].MaxFreq,
		}
		if i < len(lossRequests) {
			te.MinLoss = lossRequests[i].MinLoss
			te.PertLoss = lossRequests[i].PertLoss
			te.MaxLoss = lossRequests[i].MaxLoss
		}
		threatEventRequests = append(threatEventRequests, te)
	}

	finalResponse := FrontEndResponseT{
		FrequencyMax:      threatEventRequests[0].MaxFreq,
		FrequencyMin:      threatEventRequests[0].MinFreq,
		FrequencyEstimate: threatEventRequests[0].PertFreq,
		LossMax:           threatEventRequests[0].MaxLoss,
		LossMin:           threatEventRequests[0].MinLoss,
		LossEstimate:      threatEventRequests[0].PertLoss,
	}

	c.JSON(http.StatusOK, finalResponse)
}
