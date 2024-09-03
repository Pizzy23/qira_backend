package simulation

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulationAggregated(c *gin.Context, lossType string) {
	var riskCalculations []db.RiskCalculation
	var events []db.ThreatEventCatalog

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

	if err := db.InScope(dbEngine.NewSession(), &events); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	finalResult, err := calculationRisk(riskCalculations, events, lossType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateSimulationData(finalResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, finalResult)
}
