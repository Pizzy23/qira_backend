package risk

import (
	"net/http"
	risk "qira/internal/risk/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary {WIP} Create Risk
// @Description Create new Risk
// @Tags 6 - Risk
// @Accept json
// @Produce json
// @Success 200 {object} db.RiskCalculation "Risk Create"
// @Router /api/create-Risk [post]
func CreateRisk(c *gin.Context) {

	if risk, err := risk.CreateRiskService(c); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	} else if risk != nil {
		c.Set("Response", risk)
		c.Status(http.StatusOK)
	}
	c.Set("Response", "Risk created successfully")
	c.Status(http.StatusOK)

}

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
