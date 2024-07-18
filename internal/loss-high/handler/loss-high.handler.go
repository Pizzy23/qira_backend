package losshigh

import (
	"net/http"
	"qira/internal/interfaces"
	losshigh "qira/internal/loss-high/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create LossHigh
// @Description Create new LossHigh
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Param id header int true "Threat Event ID"
// @Param request body interfaces.InputLossHigh true "Data for create new LossHigh"
// @Success 200 {object} db.LossHigh "LossHigh Create"
// @Router /api/create-losshigh/{id} [put]
func CreateLossHigh(c *gin.Context) {
	var LossHigh interfaces.InputLossHigh

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		c.Set("Response", "Parameters are invalid, need a Id")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&LossHigh); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := losshigh.CreateLossHighService(c, LossHigh, id); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "LossHigh created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All LossHigh
// @Description Retrieve and aggregate all LossHigh records
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Success 200 {object} []db.LossHigh "List of All LossHigh with Aggregated Data"
// @Router /api/losshigh [get]
func PullAllLossHigh(c *gin.Context) {
	aggregatedLosses, err := losshigh.GetAggregatedLosses(c)
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", aggregatedLosses)
	c.Status(http.StatusOK)
}

// @Summary Retrieve LossHigh by ID
// @Description Retrieve an LossHigh by its ID
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Param id path int true "LossHigh ID"
// @Success 200 {object} db.LossHigh "LossHigh Details"
// @Router /api/losshigh/{id} [get]
func PullLosstId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Response": "Invalid LossHigh ID"})
		return
	}
	losshigh.PullLossHighId(c, id)
}
