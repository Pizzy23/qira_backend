package simulation

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulationAggregated(c *gin.Context, lossType string) {
	var riskCalculations []db.RiskCalculation
	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, "Database connection not found")
		return
	}

	dbEngine, ok := engine.(*xorm.Engine)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Failed to cast database connection to *xorm.Engine")
		return
	}

	if err := db.GetAll(dbEngine, &riskCalculations); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	for _, risk := range riskCalculations {
		if risk.RiskType == "Frequency" && risk.Categorie == lossType {
			totalMinFreq += risk.Min
			totalPertFreq += risk.Estimate
			totalMaxFreq += risk.Max
		} else if risk.Categorie == lossType {
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
