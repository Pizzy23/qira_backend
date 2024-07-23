package simulation

import (
	"fmt"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type AcceptableLoss struct {
	Risk string  `json:"risk"`
	Loss float64 `json:"loss"`
}

type FrontEndResponseApp struct {
	FrequencyMax   float64             `json:"FrequencyMax"`
	FrequencyMin   float64             `json:"FrequencyMin"`
	FrequencyMode  float64             `json:"FrequencyMode"`
	LossMax        float64             `json:"LossMax"`
	LossMin        float64             `json:"LossMin"`
	LossMode       float64             `json:"LossMode"`
	lossExceedance []db.LossExceedance `json:"LossExceedance"`
}

func MonteCarloSimulationAppetite(c *gin.Context, threatEvent string) {
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

	frequencyRequests := make([]ThreatEventRequest, len(riskCalculations))
	lossRequests := make([]ThreatEventRequest, len(riskCalculations))

	for i, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			totalMinFreq += risk.Min
			totalPertFreq += risk.Mode
			totalMaxFreq += risk.Max
			frequencyRequests[i] = ThreatEventRequest{
				MinFreq:  risk.Min,
				PertFreq: risk.Mode,
				MaxFreq:  risk.Max,
			}
		} else if risk.RiskType == "Loss" {
			totalMinLoss += risk.Min
			totalPertLoss += risk.Mode
			totalMaxLoss += risk.Max
			lossRequests[i] = ThreatEventRequest{
				MinLoss:  risk.Min,
				PertLoss: risk.Mode,
				MaxLoss:  risk.Max,
			}
		}
	}

	threatEventRequests := make([]ThreatEventRequest, len(frequencyRequests))
	for i := range frequencyRequests {
		threatEventRequests[i] = ThreatEventRequest{
			MinFreq:  frequencyRequests[i].MinFreq,
			PertFreq: frequencyRequests[i].PertFreq,
			MaxFreq:  frequencyRequests[i].MaxFreq,
			MinLoss:  lossRequests[i].MinLoss,
			PertLoss: lossRequests[i].PertLoss,
			MaxLoss:  lossRequests[i].MaxLoss,
		}
	}
	var lossEc []db.LossExceedance
	if err := db.GetAll(engine.(*xorm.Engine), &lossEc); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	finalResponse := FrontEndResponseApp{
		FrequencyMax:   totalMaxFreq,
		FrequencyMin:   totalMinFreq,
		FrequencyMode:  totalPertFreq,
		LossMax:        totalMaxLoss,
		LossMin:        totalMinLoss,
		LossMode:       totalPertLoss,
		lossExceedance: lossEc,
	}

	c.JSON(http.StatusOK, finalResponse)
}

func UploadLossData(c *gin.Context, lossData []interfaces.LossExceedance) {

	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not found",
		})
		return
	}

	for _, ld := range lossData {
		existing := &db.LossExceedance{}
		has, err := engine.(*xorm.Engine).Where("risk = ? AND loss = ?", ld.Risk, ld.Loss).Get(existing)
		if err == nil && !has {
			newLoss := db.LossExceedance{
				Risk: ld.Risk,
				Loss: ld.Loss,
			}
			_, err := engine.(*xorm.Engine).Insert(newLoss)
			if err != nil {
				fmt.Printf("Error inserting loss data: %v\n", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Loss data uploaded successfully",
	})
}
