package simulation

import (
	"net/http"
	"qira/db"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type FrontEndResponseAppReport struct {
	ProposedMin    float64             `json:"ProposedMin"`
	ProposedMax    float64             `json:"ProposedMax"`
	ProposedPert   float64             `json:"ProposedPert"`
	LossExceedance []db.LossExceedance `json:"LossExceedance"`
}

func MonteCarloSimulationRisk(c *gin.Context, threatEvent string) {
	var riskCalculations []db.RiskCalculation
	var controlGaps []db.Control

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Fetching risk calculations
	err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&riskCalculations)
	if err != nil {
		c.Set("Response", "Error retrieving risk calculations")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Fetching control gaps
	err = engine.(*xorm.Engine).Where("control_id = ?", -2).Find(&controlGaps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving control gap data"})
		return
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	// Aggregating risk values
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

	// Calculating inherent risks using control gaps
	var ihRiskMin, ihRiskMax, ihRiskEstimate float64
	for _, gap := range controlGaps {
		gapStr := strings.TrimSuffix(gap.ControlGap, "%")
		gapValue, err := strconv.ParseFloat(gapStr, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing control gap value"})
			return
		}
		if gap.TypeOfAttack == "Frequency" {
			ihRiskMin += totalMinFreq / (gapValue / 100)
			ihRiskMax += totalMaxFreq / (gapValue / 100)
			ihRiskEstimate += totalPertFreq / (gapValue / 100)
		} else if gap.TypeOfAttack == "Loss" {
			ihRiskMin += totalMinLoss / (gapValue / 100)
			ihRiskMax += totalMaxLoss / (gapValue / 100)
			ihRiskEstimate += totalPertLoss / (gapValue / 100)
		}
	}

	// Calculating proposed risks using control gaps
	var proposedMin, proposedMax, proposedEstimate float64
	for _, gap := range controlGaps {
		gapStr := strings.TrimSuffix(gap.ControlGap, "%")
		gapValue, err := strconv.ParseFloat(gapStr, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing control gap value"})
			return
		}
		if gap.TypeOfAttack == "Proposed" {
			proposedMin += ihRiskMin * (gapValue / 100)
			proposedMax += ihRiskMax * (gapValue / 100)
			proposedEstimate += ihRiskEstimate * (gapValue / 100)
		}
	}

	// Fetching loss exceedance data
	var lossEc []db.LossExceedance
	if err := db.GetAll(engine.(*xorm.Engine), &lossEc); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Format float values to strings, removing last digit for losses
	finalResponse := FrontEndResponseAppReport{
		ProposedMin:    proposedMin,
		ProposedMax:    proposedMax,
		ProposedPert:   proposedEstimate,
		LossExceedance: lossEc,
	}

	c.JSON(http.StatusOK, finalResponse)
}
