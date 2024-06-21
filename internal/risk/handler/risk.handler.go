package risk

import (
	"net/http"
	"qira/internal/interfaces"
	risk "qira/internal/risk/service"
	erros "qira/middleware/interfaces/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create Risk
// @Description Create new Risk
// @Tags Risk
// @Accept json
// @Produce json
// @Param request body interfaces.InputRiskCalculator true "Data for create new Risk"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.RiskCalculator "Risk Create"
// @Router /api/create-Risk [post]
func CreateRisk(c *gin.Context) {
	var riskInput interfaces.InputRiskCalculator

	if err := c.ShouldBindJSON(&riskInput); err != nil {
		c.JSON(erros.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if err := risk.CreateRiskService(c, riskInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("Response", "Risk created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Risks
// @Description Retrieve all Risks
// @Tags Risk
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} []interfaces.RiskCalculator "List of All Risks"
// @Router /api/Risk [get]
func PullAllRisk(c *gin.Context) {
	risk.PullAllRisk(c)
}

// @Summary Retrieve Risk by ID
// @Description Retrieve an Risk by its ID
// @Tags Risk
// @Accept json
// @Produce json
// @Param id path int true "Risk ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.RiskCalculator "Risk Details"
// @Router /api/Risk/{id} [get]
func PullRiskId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Risk ID"})
		return
	}
	risk.PullRiskId(c, id)
}
