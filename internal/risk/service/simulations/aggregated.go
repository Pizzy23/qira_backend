package simulation

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulationAggregated(c *gin.Context, lossType string) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	var riskCalculations []db.RiskCalculation

	dbEngine, ok := engine.(*xorm.Engine)
	if !ok {
		c.Set("Response", "Failed to cast database connection to *xorm.Engine")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Buscar todos os registros de RiskCalculation
	if err := db.GetAll(dbEngine, &riskCalculations); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	// Agregar valores de risco com base no tipo de perda (lossType)
	for _, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			totalMinFreq += risk.Min
			totalPertFreq += risk.Estimate
			totalMaxFreq += risk.Max
		} else if risk.RiskType == lossType {
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
