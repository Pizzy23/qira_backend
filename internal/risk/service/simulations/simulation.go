package simulation

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type ThreatEventRequest struct {
	MinFreq  float64 `json:"minfreq,omitempty"`
	PertFreq float64 `json:"pertfreq,omitempty"`
	MaxFreq  float64 `json:"maxfreq,omitempty"`
	MinLoss  float64 `json:"minloss,omitempty"`
	PertLoss float64 `json:"pertloss,omitempty"`
	MaxLoss  float64 `json:"maxloss,omitempty"`
}

type FrontEndResponse struct {
	FrequencyMax      float64 `json:"FrequencyMax"`
	FrequencyMin      float64 `json:"FrequencyMin"`
	FrequencyEstimate float64 `json:"FrequencyEstimate"`
	LossMax           float64 `json:"LossMax"`
	LossMin           float64 `json:"LossMin"`
	LossEstimate      float64 `json:"LossEstimate"`
}

func MonteCarloSimulation(c *gin.Context, threatEvent string) {
	var frequencyEntries []db.Frequency
	var lossEntries []db.LossHighTotal

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Retrieve Frequency entries
	err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&frequencyEntries)
	if err != nil {
		c.Set("Response", "Error retrieving frequency entries")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Retrieve Loss entries
	err = engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&lossEntries)
	if err != nil {
		c.Set("Response", "Error retrieving loss entries")
		c.Status(http.StatusInternalServerError)
		return
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	// Sum Frequency values
	for _, freq := range frequencyEntries {
		totalMinFreq += freq.MinFrequency
		totalPertFreq += freq.MostLikelyFrequency
		totalMaxFreq += freq.MaxFrequency
	}

	// Sum Loss values
	for _, loss := range lossEntries {
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
