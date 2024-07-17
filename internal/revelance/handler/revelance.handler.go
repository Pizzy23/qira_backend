package relevance

import (
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	revelance "qira/internal/revelance/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve All Relevance
// @Description Retrieve all Relevance
// @Tags Revelance
// @Accept json
// @Produce json
// @Success 200 {object} []db.Relevance "List of All Relevance"
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
// @Success 200 {object} db.Relevance "Revelance Details"
// @Router /api/revelance/{id} [get]
func PullRevelanceId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}
	revelance.PullRevelanceId(c, id)
}

// @Summary Create Relevance
// @Description Create new Relevance
// @Tags Revelance
// @Accept json
// @Produce json
// @Param request body db.Relevance true "Data for create new Relevance"
// @Success 200 {object} db.Relevance "Relevance Create"
// @Router /api/create-revelance [post]
func CreateRelevance(c *gin.Context) {
	var RelevanceInput db.Relevance
	if err := c.ShouldBindJSON(&RelevanceInput); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := revelance.CreateRelevanceService(c, RelevanceInput); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Relevance created successfully")
	c.Status(http.StatusOK)

}

// @Summary {WIP} Update Relevance
// @Description Update an existing Relevance
// @Tags Revelance
// @Accept json
// @Produce json
// @Param id path int true "Relevance ID"
// @Param request body interfaces.RelevanceDinamicInput true "Data to update Relevance"
// @Success 200 {object} db.Relevance "Relevance Updated"
// @Router /api/revelance/{id} [put]
func UpdateRelevance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	var relevanceInput interfaces.RelevanceDinamicInput
	if err := c.ShouldBindJSON(&relevanceInput); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := revelance.UpdateRelevanceService(c, id, relevanceInput); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Relevance updated successfully")
	c.Status(http.StatusOK)
}
