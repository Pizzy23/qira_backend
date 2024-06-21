package losshigh

import (
	"net/http"
	"qira/internal/interfaces"
	losshigh "qira/internal/loss-high/service"
	erros "qira/middleware/interfaces/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create LossHigh
// @Description Create new LossHigh
// @Tags losshigh
// @Accept json
// @Produce json
// @Param request body interfaces.InputLossHigh true "Data for create new LossHigh"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.LossHigh "LossHigh Create"
// @Router /api/create-losshigh [post]
func CreateLossHigh(c *gin.Context) {
	var LossHigh interfaces.InputLossHigh

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

// @Summary Retrieve All LossHighs
// @Description Retrieve all LossHighs
// @Tags losshigh
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} []interfaces.LossHigh "List of All LossHighs"
// @Router /api/losshigh [get]
func PullAllLoss(c *gin.Context) {
	losshigh.PullAllLossHigh(c)
}

// @Summary Retrieve LossHigh by ID
// @Description Retrieve an LossHigh by its ID
// @Tags losshigh
// @Accept json
// @Produce json
// @Param id path int true "LossHigh ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.LossHigh "LossHigh Details"
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
