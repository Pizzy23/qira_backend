package environment

import (
	"net/http"
	riskService "qira/internal/environment/service"
	"qira/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create Risk Assessment
// @Description Create new Risk Assessment
// @Tags 0 - Environment
// @Accept json
// @Produce json
// @Param request body interfaces.InputRiskAssessment true "Data for create new Risk Assessment"
// @Success 200 {object} db.RiskAssessment "Risk Assessment Created"
// @Router /api/create-risk-assessment [post]
func CreateRiskAssessment(c *gin.Context) {
	var risk interfaces.InputRiskAssessment

	if err := c.ShouldBindJSON(&risk); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := riskService.CreateRiskAssessmentService(c, risk); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Risk Assessment created successfully")
	c.Status(http.StatusOK)
}

// @Summary Retrieve All Risk Assessments
// @Description Retrieve all Risk Assessments
// @Tags 0 - Environment
// @Accept json
// @Produce json
// @Success 200 {object} []db.RiskAssessment "List of All Risk Assessments"
// @Router /api/risk-assessments [get]
func PullAllRiskAssessments(c *gin.Context) {
	riskService.PullAllRiskAssessments(c)
}

// @Summary Retrieve Risk Assessment by ID
// @Description Retrieve a Risk Assessment by its ID
// @Tags 0 - Environment
// @Accept json
// @Produce json
// @Param id path int true "Risk Assessment ID"
// @Success 200 {object} db.RiskAssessment "Risk Assessment Details"
// @Router /api/risk-assessment/{id} [get]
func PullRiskAssessmentById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}
	riskService.PullRiskAssessmentById(c, id)
}

// @Summary Update Risk Assessment
// @Description Update an existing Risk Assessment
// @Tags 0 - Environment
// @Accept json
// @Produce json
// @Param id path int true "Risk Assessment ID"
// @Param request body interfaces.InputRiskAssessment true "Data to update Risk Assessment"
// @Success 200 {object} db.RiskAssessment "Risk Assessment Updated"
// @Router /api/risk-assessment/{id} [put]
func UpdateRiskAssessment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	var risk interfaces.InputRiskAssessment
	if err := c.ShouldBindJSON(&risk); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := riskService.UpdateRiskAssessmentService(c, id, risk); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Risk Assessment updated successfully")
	c.Status(http.StatusOK)
}

// @Summary Delete Risk Assessment
// @Description Delete an existing Risk Assessment
// @Tags 0 - Environment
// @Accept json
// @Produce json
// @Param id path int true "Risk Assessment ID"
// @Success 200 {object} db.RiskAssessment "Risk Assessment Deleted"
// @Router /api/risk-assessment/{id} [delete]
func DeleteRiskAssessment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := riskService.DeleteRiskAssessment(c, id); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Risk Assessment deleted successfully")
	c.Status(http.StatusOK)
}
