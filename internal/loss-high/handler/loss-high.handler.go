package losshigh

import (
	"net/http"
	"qira/db"
	losshigh "qira/internal/loss-high/service"
	erros "qira/middleware/interfaces/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create LossHigh
// @Description Create new LossHigh
// @Tags Losshigh
// @Accept json
// @Produce json
// @Param request body db.LossHigh true "Data for create new LossHigh"
// @Success 200 {object} db.LossHigh "LossHigh Create"
// @Router /api/create-losshigh [post]
func CreateLossHigh(c *gin.Context) {
	var LossHigh db.LossHigh

	if err := c.ShouldBindJSON(&LossHigh); err != nil {
		c.JSON(erros.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if err := losshigh.CreateLossHighService(c, LossHigh); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("Response", "LossHigh created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All LossHigh
// @Description Retrieve and aggregate all LossHigh records
// @Tags LossHigh
// @Accept json
// @Produce json
// @Success 200 {object} []db.LossHigh "List of All LossHigh with Aggregated Data"
// @Router /api/losshigh [get]
func PullAllLossHigh(c *gin.Context) {
	aggregatedLosses, err := losshigh.GetAggregatedLosses(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aggregatedLosses)
}

// @Summary Retrieve LossHigh by ID
// @Description Retrieve an LossHigh by its ID
// @Tags Losshigh
// @Accept json
// @Produce json
// @Param id path int true "LossHigh ID"
// @Success 200 {object} db.LossHigh "LossHigh Details"
// @Router /api/losshigh/{id} [get]
func PullLosstId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid LossHigh ID"})
		return
	}
	losshigh.PullLossHighId(c, id)
}
