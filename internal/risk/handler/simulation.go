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
// @Success 200 {object} []db.RiskCalculation "List of All Risks"
// @Router /api/simulation [get]
func RiskMount(c *gin.Context) {
	simulation.MonteCarloSimulation(c)
}

// @Summary Test for simulation aggregated
// @Description Test for simulation aggregated
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Success 200 {object} []db.RiskCalculation "List of All Risks"
// @Router /api/simulation-aggregated [get]
func RiskMountAggregated(c *gin.Context) {
	simulation.MonteCarloSimulationAggregated(c)
}
