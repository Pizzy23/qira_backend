package simulation

import (
	"errors"
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulation(c *gin.Context, threatEvent string, lossType string) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	dbEngine, ok := engine.(*xorm.Engine)
	if !ok {
		c.Set("Response", "Failed to cast database connection to *xorm.Engine")
		c.Status(http.StatusInternalServerError)
		return
	}

	freq, loss, err := retrieveFrequencyAndLossEntries(dbEngine, threatEvent, lossType)
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
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

func retrieveFrequencyAndLossEntries(engine *xorm.Engine, threatEvent, lossType string) ([]db.Frequency, []db.LossHighTotal, error) {
	var frequencyEntries []db.Frequency
	var lossEntries []db.LossHighTotal

	err := engine.Where("threat_event = ?", threatEvent).Find(&frequencyEntries)
	if err != nil {
		return nil, nil, errors.New("error retrieving frequency entries")
	}

	err = engine.Where("threat_event = ? AND type_of_loss = ?", threatEvent, lossType).Find(&lossEntries)
	if err != nil {
		return nil, nil, errors.New("error retrieving loss entries")
	}

	return frequencyEntries, lossEntries, nil
}
