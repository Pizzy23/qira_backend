package risk

import (
	"net/http"
	"qira/internal/interfaces"
	risk "qira/internal/risk/service"
	erros "qira/middleware/interfaces/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary {WIP} Create Risk
// @Description Create new Risk
// @Tags Risk
// @Accept json
// @Produce json
// @Param request body interfaces.RiskCalc true "Data for create new Risk"
// @Success 200 {object} db.RiskCalculation "Risk Create"
// @Router /api/create-Risk [post]
func CreateRisk(c *gin.Context) {
	var riskInput interfaces.RiskCalc

	if err := c.ShouldBindJSON(&riskInput); err != nil {
		c.JSON(erros.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if risk, err := risk.CreateRiskService(c, riskInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Tags Risk
// @Accept json
// @Produce json
// @Success 200 {object} []db.RiskCalculation "List of All Risks"
// @Router /api/risk [get]
func PullAllRisk(c *gin.Context) {
	risk.PullAllRisk(c)
}

// @Summary Retrieve Risk by ID
// @Description Retrieve an Risk by its ID
// @Tags Risk
// @Accept json
// @Produce json
// @Param id path int true "Risk ID"
// @Success 200 {object} db.RiskCalculation "Risk Details"
// @Router /api/risk/{id} [get]
func PullRiskId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Risk ID"})
		return
	}
	risk.PullRiskId(c, id)
}
