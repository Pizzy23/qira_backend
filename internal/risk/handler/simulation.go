package risk

import (
	"net/http"
	"qira/internal/interfaces"
	simulation "qira/internal/risk/service/simulations"

	"github.com/gin-gonic/gin"
)

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param threatEvent header string true "Threat Event "
// @Param email header string false "Email for recive "
// @Router /simulation [get]
func RiskMount(c *gin.Context) {
	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	email := c.GetHeader("Email")
	simulation.MonteCarloSimulation(c, threatEvent, email)
}

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param threatEvent header string true "Threat Event "
// @Param email header string false "Email for recive"
// @Router /simulation-report [get]
func RiskMountReport(c *gin.Context) {
	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	email := c.GetHeader("Email")
	simulation.MonteCarloSimulationReport(c, threatEvent, email)
}

// @Summary Test for simulation aggregated
// @Description Test for simulation aggregated
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param threatEvent header string true "Threat Event "
// @Param email header string false "Email for recive "
// @Router /simulation-aggregated [get]
func RiskMountAggregated(c *gin.Context) {
	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	email := c.GetHeader("Email")
	simulation.MonteCarloSimulationAggregated(c, threatEvent, email)
}

// @Summary Test for simulation appetite
// @Description Test for simulation appetite
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param request body []interfaces.LossExceedance true "Loss Exceedance Graph"
// @Param threatEvent header string true "Threat Event "
// @Param email header string false "Email for recive "
// @Router /simulation-appetite [get]
func RiskMountAppetite(c *gin.Context) {
	var lossData []interfaces.LossExceedance
	if err := c.ShouldBindJSON(&lossData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters are invalid, need a JSON array of LossExceedance"})
		return
	}

	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	email := c.GetHeader("Email")
	simulation.MonteCarloSimulationAppetite(c, threatEvent, email, lossData)
}
