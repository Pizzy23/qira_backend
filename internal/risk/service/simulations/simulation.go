package simulation

import (
	"net/http"
	"qira/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type ThreatEventRequest struct {
	MinFreq  string `json:"minfreq,omitempty"`
	PertFreq string `json:"pertfreq,omitempty"`
	MaxFreq  string `json:"maxfreq,omitempty"`
	MinLoss  string `json:"minloss,omitempty"`
	PertLoss string `json:"pertloss,omitempty"`
	MaxLoss  string `json:"maxloss,omitempty"`
}

type FrontEndResponse struct {
	FrequencyMax      string `json:"FrequencyMax"`
	FrequencyMin      string `json:"FrequencyMin"`
	FrequencyEstimate string `json:"FrequencyEstimate"`
	LossMax           string `json:"LossMax"`
	LossMin           string `json:"LossMin"`
	LossEstimate      string `json:"LossEstimate"`
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

	frequencyRequests := []ThreatEventRequest{}
	lossRequests := []ThreatEventRequest{}

	for _, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			frequencyRequests = append(frequencyRequests, ThreatEventRequest{
				MinFreq:  strconv.FormatFloat(risk.Min, 'f', -1, 64),
				PertFreq: strconv.FormatFloat(risk.Estimate, 'f', -1, 64),
				MaxFreq:  strconv.FormatFloat(risk.Max, 'f', -1, 64),
			})
		} else if risk.RiskType == "Loss" {
			lossRequests = append(lossRequests, ThreatEventRequest{
				MinLoss:  strconv.FormatFloat(risk.Min, 'f', -1, 64),
				PertLoss: strconv.FormatFloat(risk.Estimate, 'f', -1, 64),
				MaxLoss:  strconv.FormatFloat(risk.Max, 'f', -1, 64),
			})
		}
	}

	threatEventRequests := []ThreatEventRequest{}
	for i := range frequencyRequests {
		te := ThreatEventRequest{
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

	finalResponse := FrontEndResponse{
		FrequencyMax:      threatEventRequests[0].MaxFreq,
		FrequencyMin:      threatEventRequests[0].MinFreq,
		FrequencyEstimate: threatEventRequests[0].PertFreq,
		LossMax:           threatEventRequests[0].MaxLoss,
		LossMin:           threatEventRequests[0].MinLoss,
		LossEstimate:      threatEventRequests[0].PertLoss,
	}

	c.JSON(http.StatusOK, finalResponse)
}
