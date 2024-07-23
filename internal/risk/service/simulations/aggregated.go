package simulation

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulationAggregated(c *gin.Context, threatEvent string) {
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

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	frequencyRequests := make([]ThreatEventRequest, len(riskCalculations))
	lossRequests := make([]ThreatEventRequest, len(riskCalculations))

	for i, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			totalMinFreq += risk.Min
			totalPertFreq += risk.Mode
			totalMaxFreq += risk.Max
			frequencyRequests[i] = ThreatEventRequest{
				MinFreq:  risk.Min,
				PertFreq: risk.Mode,
				MaxFreq:  risk.Max,
			}
		} else if risk.RiskType == "Loss" {
			totalMinLoss += risk.Min
			totalPertLoss += risk.Mode
			totalMaxLoss += risk.Max
			lossRequests[i] = ThreatEventRequest{
				MinLoss:  risk.Min,
				PertLoss: risk.Mode,
				MaxLoss:  risk.Max,
			}
		}
	}

	threatEventRequests := make([]ThreatEventRequest, len(frequencyRequests))
	for i := range frequencyRequests {
		threatEventRequests[i] = ThreatEventRequest{
			MinFreq:  frequencyRequests[i].MinFreq,
			PertFreq: frequencyRequests[i].PertFreq,
			MaxFreq:  frequencyRequests[i].MaxFreq,
			MinLoss:  lossRequests[i].MinLoss,
			PertLoss: lossRequests[i].PertLoss,
			MaxLoss:  lossRequests[i].MaxLoss,
		}
	}
	var lossEc []db.LossExceedance
	if err := db.GetAll(engine.(*xorm.Engine), &lossEc); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	finalResponse := FrontEndResponseApp{
		FrequencyMax:   totalMaxFreq,
		FrequencyMin:   totalMinFreq,
		FrequencyMode:  totalPertFreq,
		LossMax:        totalMaxLoss,
		LossMin:        totalMinLoss,
		LossMode:       totalPertLoss,
		lossExceedance: lossEc,
	}

	c.JSON(http.StatusOK, finalResponse)
}
