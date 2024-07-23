package simulation

import (
	"net/http"
	"qira/db"
	"strconv"

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

	for _, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			totalMinFreq += risk.Min
			totalPertFreq += risk.Estimate
			totalMaxFreq += risk.Max
		} else if risk.RiskType == "Loss" {
			totalMinLoss += risk.Min
			totalPertLoss += risk.Estimate
			totalMaxLoss += risk.Max
		}
	}

	var lossEc []db.LossExceedance
	if err := db.GetAll(engine.(*xorm.Engine), &lossEc); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	finalResponse := FrontEndResponseApp{
		FrequencyMax:      strconv.FormatFloat(totalMaxFreq, 'f', -1, 64),
		FrequencyMin:      strconv.FormatFloat(totalMinFreq, 'f', -1, 64),
		FrequencyEstimate: strconv.FormatFloat(totalPertFreq, 'f', -1, 64),
		LossMax:           strconv.FormatFloat(totalMaxLoss, 'f', -1, 64),
		LossMin:           strconv.FormatFloat(totalMinLoss, 'f', -1, 64),
		LossEstimate:      strconv.FormatFloat(totalPertLoss, 'f', -1, 64),
		LossExceedance:    lossEc,
	}

	c.JSON(http.StatusOK, finalResponse)
}
