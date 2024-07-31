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

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	for _, freq := range freq {
		totalMinFreq += freq.MinFrequency
		totalPertFreq += freq.MostLikelyFrequency
		totalMaxFreq += freq.MaxFrequency
	}

	for _, loss := range loss {
		totalMinLoss += loss.MinimumLoss
		totalPertLoss += loss.MostLikelyLoss
		totalMaxLoss += loss.MaximumLoss
	}

	finalResponse := FrontEndResponse{
		FrequencyMax:      totalMaxFreq,
		FrequencyMin:      totalMinFreq,
		FrequencyEstimate: totalPertFreq,
		LossMax:           totalMaxLoss,
		LossMin:           totalMinLoss,
		LossEstimate:      totalPertLoss,
	}

	c.JSON(http.StatusOK, finalResponse)
}
