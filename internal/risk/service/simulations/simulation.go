package simulation

import (
	"net/http"
	"qira/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type ThreatEventRequestT struct {
	MinFreq  float64 `json:"minfreq,omitempty"`
	PertFreq float64 `json:"pertfreq,omitempty"`
	MaxFreq  float64 `json:"maxfreq,omitempty"`
	MinLoss  float64 `json:"minloss,omitempty"`
	PertLoss float64 `json:"pertloss,omitempty"`
	MaxLoss  float64 `json:"maxloss,omitempty"`
}

type FrontEndResponseT struct {
	FrequencyMax      float64 `json:"FrequencyMax"`
	FrequencyMin      float64 `json:"FrequencyMin"`
	FrequencyEstimate float64 `json:"FrequencyEstimate"`
	LossMax           float64 `json:"LossMax"`
	LossMin           float64 `json:"LossMin"`
	LossEstimate      float64 `json:"LossEstimate"`
}

func truncateToOneDecimal(value float64) float64 {
	strValue := strconv.FormatFloat(value, 'f', 1, 64)
	truncatedValue, _ := strconv.ParseFloat(strValue, 64)
	return truncatedValue
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

	finalResponse := FrontEndResponseT{
		FrequencyMax:      totalMaxFreq,
		FrequencyMin:      totalMinFreq,
		FrequencyEstimate: totalPertFreq,
		LossMax:           truncateToOneDecimal(totalMaxLoss),
		LossMin:           truncateToOneDecimal(totalMinLoss),
		LossEstimate:      truncateToOneDecimal(totalPertLoss),
	}

	c.JSON(http.StatusOK, finalResponse)
}
