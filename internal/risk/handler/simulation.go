package risk

import (
	"net/http"
	simulation "qira/internal/risk/service/simulations"

	"github.com/gin-gonic/gin"
)

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param threatEvent header int true "Threat Event "
// @Param email header int true "Email for recive "
// @Router /simulation [get]
func RiskMount(c *gin.Context) {
	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	email := c.GetHeader("Email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email header is required"})
		return
	}
	simulation.MonteCarloSimulation(c, threatEvent, email)
}

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Param threatEvent header int true "Threat Event "
// @Param email header int true "Email for recive "
// @Router /simulation-report [get]
func RiskMountReport(c *gin.Context) {
	threatEvent := c.GetHeader("ThreatEvent")
	if threatEvent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ThreatEvent header is required"})
		return
	}
	email := c.GetHeader("Email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email header is required"})
		return
	}
	simulation.MonteCarloSimulationReport(c, threatEvent, email)
}

// @Summary Test for simulation aggregated
// @Description Test for simulation aggregated
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Router /simulation-aggregated [get]
func RiskMountAggregated(c *gin.Context) {
	//simulation.MonteCarloSimulationAggregated(c)
}

// @Summary Test for simulation appetite
// @Description Test for simulation appetite
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Router /simulation-appetite [get]
func RiskMountAppetite(c *gin.Context) {
	//simulation.MonteCarloSimulationAppetite(c)
}
