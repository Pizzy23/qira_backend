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
	var events []db.ThreatEventCatalog

	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, "Database connection not found")
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &riskCalculations); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.InScope(engine.(*xorm.Engine).NewSession(), &events); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	finalResult, err := calculationRisk(riskCalculations, events, lossType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lossEc []db.LossExceedance
	if err := db.GetAll(engine.(*xorm.Engine), &lossEc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	finalResponse := FrontEndResponseAppLoss{
		FrequencyMax:      finalResult.FrequencyMax,
		FrequencyMin:      finalResult.FrequencyMin,
		FrequencyEstimate: finalResult.FrequencyEstimate,
		LossMax:           finalResult.LossMax,
		LossMin:           finalResult.LossMin,
		LossEstimate:      finalResult.LossEstimate,
		LossExceedance:    lossEc,
	}

	if err := validateSimulationData(finalResponse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, finalResponse)
}

func UploadLossData(c *gin.Context, lossData []interfaces.LossExceedance) {
	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not found"})
		return
	}

	for _, ld := range lossData {
		var existing db.LossExceedance
		has, err := engine.(*xorm.Engine).Where("risk = ?", ld.Risk).Get(&existing)
		if err != nil {
			fmt.Printf("Error fetching existing loss data: %v\n", err)
			continue
		}

		if has && existing.Risk == 0 || has && existing.Loss != ld.Loss {
			existing.Risk = ld.Risk
			existing.Loss = ld.Loss
			if _, err := engine.(*xorm.Engine).ID(existing.ID).Update(&existing); err != nil {
				fmt.Printf("Error updating loss data: %v\n", err)
			}
		} else if !has {
			newLoss := db.LossExceedance{
				Risk: ld.Risk,
				Loss: ld.Loss,
			}
			if _, err := engine.(*xorm.Engine).Insert(&newLoss); err != nil {
				fmt.Printf("Error inserting loss data: %v\n", err)
			}
		}
	}
	c.JSON(http.StatusCreated, "LossExceedance Updated")
}
