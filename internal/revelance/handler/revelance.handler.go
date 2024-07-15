package revelance

import (
	"net/http"
	"qira/db"
	revelance "qira/internal/revelance/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve All Revelances
// @Description Retrieve all Revelances
// @Tags Revelance
// @Accept json
// @Produce json
// @Success 200 {object} []db.RelevanceDinamic "List of All Revelances"
// @Router /api/revelance [get]
func PullAllRevelance(c *gin.Context) {
	revelance.PullAllRevelance(c)
}

// @Summary Retrieve Revelance by ID
// @Description Retrieve an Revelance by its ID
// @Tags Revelance
// @Accept json
// @Produce json
// @Param id path int true "Revelance ID"
// @Success 200 {object} db.RelevanceDinamic "Revelance Details"
// @Router /api/revelance/{id} [get]
func PullRevelanceId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Revelance ID"})
		return
	}
	revelance.PullRevelanceId(c, id)
}

// @Summary Create Relevance
// @Description Create new Relevance
// @Tags Revelance
// @Accept json
// @Produce json
// @Param request body db.RelevanceDinamicInput true "Data for create new Relevance"
// @Success 200 {object} db.LossHigh "LossHigh Create"
// @Router /api/create-revelance [post]
func CreateRelevance(c *gin.Context) {
	var RelevanceInput db.RelevanceDinamicInput

	if err := c.ShouldBindJSON(&RelevanceInput); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if err := revelance.CreateRelevanceService(c, RelevanceInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("Response", "Relevance created successfully")
	c.Status(http.StatusOK)

}
