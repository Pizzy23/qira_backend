package simulation

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type FrontEndResponseAgg struct {
	FrequencyMax      float64 `json:"FrequencyMax"`
	FrequencyMin      float64 `json:"FrequencyMin"`
	FrequencyEstimate float64 `json:"FrequencyEstimate"`
	LossMax           float64 `json:"LossMax"`
	LossMin           float64 `json:"LossMin"`
	LossEstimate      float64 `json:"LossEstimate"`
}

func MonteCarloSimulationAggregated(c *gin.Context) {
	var riskCalculations []db.RiskCalculation

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Use GetAll to fetch all records
	if err := db.GetAll(engine.(*xorm.Engine), &riskCalculations); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	// Aggregate the values
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

	finalResponse := FrontEndResponseAgg{
		FrequencyMax:      totalMaxFreq,
		FrequencyMin:      totalMinFreq,
		FrequencyEstimate: totalPertFreq,
		LossMax:           totalMaxLoss,
		LossMin:           totalMinLoss,
		LossEstimate:      totalPertLoss,
	}

	c.JSON(http.StatusOK, finalResponse)
}
