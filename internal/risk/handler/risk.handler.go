package risk

import (
	"net/http"
	risk "qira/internal/risk/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve All Risks
// @Description Retrieve all Risks
// @Tags 6 - Risk
// @Accept json
// @Produce json
// @Success 200 {object} []db.RiskCalculation "List of All Risks"
// @Router /api/risk [get]
func PullAllRisk(c *gin.Context) {
	risk.PullAllRisk(c)
}

// @Summary Test for simulation
// @Description Test for simulation
// @Tags 13 - Simulation
// @Accept json
// @Produce json
// @Success 200 {object} []db.RiskCalculation "List of All Risks"
// @Router /api/simulation [get]
func RiskMount(c *gin.Context) {
	risk.MainSimulation(c)
}

// @Summary Retrieve Risk by ID
// @Description Retrieve an Risk by its ID
// @Tags 6 - Risk
// @Accept json
// @Produce json
// @Param id path int true "Threat event ID"
// @Success 200 {object} db.RiskCalculation "Risk Details"
// @Router /api/risk/{id} [get]
func PullRiskId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}
	risk.PullRiskId(c, id)
}
