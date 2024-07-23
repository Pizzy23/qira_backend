package simulation

import (
	"fmt"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type AcceptableLoss struct {
	Risk string  `json:"risk"`
	Loss float64 `json:"loss"`
}

type FrontEndResponseApp struct {
	FrequencyMax      string              `json:"FrequencyMax"`
	FrequencyMin      string              `json:"FrequencyMin"`
	FrequencyEstimate string              `json:"FrequencyEstimate"`
	LossMax           string              `json:"LossMax"`
	LossMin           string              `json:"LossMin"`
	LossEstimate      string              `json:"LossEstimate"`
	LossExceedance    []db.LossExceedance `json:"LossExceedance"`
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
			totalPertFreq += risk.Estimate
			totalMaxFreq += risk.Max
			frequencyRequests[i] = ThreatEventRequest{
				MinFreq:  strconv.FormatFloat(risk.Min, 'f', -1, 64),
				PertFreq: strconv.FormatFloat(risk.Estimate, 'f', -1, 64),
				MaxFreq:  strconv.FormatFloat(risk.Max, 'f', -1, 64),
			}
		} else if risk.RiskType == "Loss" {
			totalMinLoss += risk.Min
			totalPertLoss += risk.Estimate
			totalMaxLoss += risk.Max
			lossRequests[i] = ThreatEventRequest{
				MinLoss:  strconv.FormatFloat(risk.Min, 'f', -1, 64),
				PertLoss: strconv.FormatFloat(risk.Estimate, 'f', -1, 64),
				MaxLoss:  strconv.FormatFloat(risk.Max, 'f', -1, 64),
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
		FrequencyMax:      strconv.FormatFloat(totalMaxFreq, 'f', -1, 64),
		FrequencyMin:      strconv.FormatFloat(totalMinFreq, 'f', -1, 64),
		FrequencyEstimate: strconv.FormatFloat(totalPertFreq, 'f', -1, 64),
		LossMax:           strconv.FormatFloat(totalMaxLoss, 'f', -1, 64),
		LossMin:           strconv.FormatFloat(totalMinLoss, 'f', -1, 64),
		LossEstimate:      strconv.FormatFloat(totalPertLoss, 'f', -1, 64),
		LossExceedance:    lossEc,
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
