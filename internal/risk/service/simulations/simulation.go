package simulation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulation(c *gin.Context, threatEvent string, lossType string) {
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

	freq, loss, err := retrieveFrequencyAndLossEntries(dbEngine, threatEvent, lossType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	final, err := calculationLossAndFreq(freq, loss)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := validateSimulationData(final); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, final)
}
