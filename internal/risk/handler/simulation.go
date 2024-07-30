package risk

import (
	"net/http"
	"qira/internal/interfaces"
	simulation "qira/internal/risk/service/simulations"
	"qira/util"

	"github.com/gin-gonic/gin"
)

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param threatEvent header string true "Threat Event "
// @Param Loss header string true "Tipo de loss" Enums(Singular,LossHigh,Granular)
// @Router /simulation [get]
func RiskMount(c *gin.Context) {
	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	typeLoss := c.GetHeader("Loss")
	vali := util.EnumLoss(typeLoss)
	if !vali {
		c.Set("Response", "TypeLoss its not valid use: `Singular, LossHigh, Granular`")
		c.Status(http.StatusInternalServerError)
	}
	simulation.MonteCarloSimulation(c, threatEvent, typeLoss)
}

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param threatEvent header string true "Threat Event "
// @Param Loss header string true "Tipo de loss" Enums(Singular,LossHigh,Granular)
// @Router /simulation-report [get]
func RiskMountReport(c *gin.Context) {
	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	typeLoss := c.GetHeader("Loss")
	vali := util.EnumLoss(typeLoss)
	if !vali {
		c.Set("Response", "TypeLoss its not valid use: `Singular, LossHigh, Granular`")
		c.Status(http.StatusInternalServerError)
	}
	simulation.MonteCarloSimulationRisk(c, threatEvent, typeLoss)
}

// @Summary Test for simulation aggregated
// @Description Test for simulation aggregated
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param Loss header string true "Tipo de loss" Enums(Singular,LossHigh,Granular)
// @Router /simulation-aggregated [get]
func RiskMountAggregated(c *gin.Context) {
	typeLoss := c.GetHeader("Loss")
	vali := util.EnumLoss(typeLoss)
	if !vali {
		c.Set("Response", "TypeLoss its not valid use: `Singular, LossHigh, Granular`")
		c.Status(http.StatusInternalServerError)
	}
	simulation.MonteCarloSimulationAggregated(c, typeLoss)
}

// @Summary Test for simulation appetite
// @Description Test for simulation appetite
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param Loss header string true "Tipo de loss" Enums(Singular,LossHigh,Granular)
// @Router /simulation-appetite [get]
func RiskMountAppetite(c *gin.Context) {
	typeLoss := c.GetHeader("Loss")
	vali := util.EnumLoss(typeLoss)
	if !vali {
		c.Set("Response", "TypeLoss its not valid use: `Singular, LossHigh, Granular`")
		c.Status(http.StatusInternalServerError)
	}
	simulation.MonteCarloSimulationAppetite(c, typeLoss)
}

// @Summary Test for simulation appetite
// @Description Test for simulation appetite
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param request body []interfaces.LossExceedance true "Loss Exceedance Graph"
// @Router /api/upload-appetite [put]
func UploadAppetite(c *gin.Context) {
	var lossData []interfaces.LossExceedance
	if err := c.ShouldBindJSON(&lossData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters are invalid, need a JSON array of LossExceedance"})
		return
	}
	simulation.UploadLossData(c, lossData)
}
