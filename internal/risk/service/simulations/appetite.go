package simulation

import (
	"fmt"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulationAppetite(c *gin.Context, lossType string) {
	var riskCalculations []db.RiskCalculation
	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, "Database connection not found")
		return
	}
	if err := db.GetAll(engine.(*xorm.Engine), &riskCalculations); err != nil {
		c.JSON(http.StatusInternalServerError, err)
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

	var lossEc []db.LossExceedance
	if err := db.GetAll(engine.(*xorm.Engine), &lossEc); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	finalResponse := FrontEndResponseAppLoss{
		FrequencyMax:      totalMaxFreq,
		FrequencyMin:      totalMinFreq,
		FrequencyEstimate: totalPertFreq,
		LossMax:           totalMaxLoss,
		LossMin:           totalMinLoss,
		LossEstimate:      totalPertLoss,
		LossExceedance:    lossEc,
	}

	c.JSON(http.StatusOK, finalResponse)
}

func UploadLossData(c *gin.Context, lossData []interfaces.LossExceedance) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database dont find")
		c.Status(http.StatusInternalServerError)
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

	c.Set("Response", "LossExceedance Update")
	c.Status(http.StatusCreated)
}
