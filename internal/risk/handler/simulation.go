package risk

import (
	simulation "qira/internal/risk/service/simulations"

	"github.com/gin-gonic/gin"
)

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Router /api/simulation [get]
func RiskMount(c *gin.Context) {
	simulation.MonteCarloSimulation(c)
}

// @Summary Test for simulation aggregated
// @Description Test for simulation aggregated
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Router /api/simulation-aggregated [get]
func RiskMountAggregated(c *gin.Context) {
	simulation.MonteCarloSimulationAggregated(c)
}

// @Summary Test for simulation appetite
// @Description Test for simulation appetite
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Router /api/simulation-appetite [get]
func RiskMountAppetite(c *gin.Context) {
	simulation.MonteCarloSimulationAppetite(c)
}
