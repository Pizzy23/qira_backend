package simulation

import (
	"net/http"
	"qira/db"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func MonteCarloSimulationRisk(c *gin.Context, threatEvent string, lossType string) {
	var controlGaps []db.Control

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

	// Buscar entradas de frequência e perda
	frequencies, losses, err := retrieveFrequencyAndLossEntries(dbEngine, threatEvent, lossType)
	if err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// Buscar control gaps
	err = dbEngine.Where("control_i_d = ?", -2).Find(&controlGaps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving control gap data"})
		return
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	// Agregar valores de frequência
	for _, freq := range frequencies {
		totalMinFreq += freq.MinFrequency
		totalPertFreq += freq.MostLikelyFrequency
		totalMaxFreq += freq.MaxFrequency
	}

	// Agregar valores de perda
	for _, loss := range losses {
		totalMinLoss += loss.MinimumLoss
		totalPertLoss += loss.MostLikelyLoss
		totalMaxLoss += loss.MaximumLoss
	}

	// Calcular riscos inerentes usando control gaps
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

	// Calcular riscos propostos usando control gaps
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

	// Buscar dados de loss exceedance
	var lossEc []db.LossExceedance
	if err := dbEngine.Find(&lossEc); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// Preparar a resposta final
	finalResponse := FrontEndResponseAppReport{
		ProposedMin:    proposedMin,
		ProposedMax:    proposedMax,
		ProposedPert:   proposedEstimate,
		LossExceedance: lossEc,
	}

	c.JSON(http.StatusOK, finalResponse)
}
