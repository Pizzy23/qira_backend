package frequency

import (
	"net/http"
	frequency "qira/internal/frequency/service"
	"qira/internal/interfaces"
	erros "qira/middleware/interfaces/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Edit Frequency
// @Description Edit Frequency
// @Tags Frequency
// @Accept json
// @Produce json
// @Param request body interfaces.InputFrequency true "Edit Frequency"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Frequency "Your Frequency is by add"
// @Router /api/frequency [put]
func EditFrequency(c *gin.Context) {
	var frequencyInput interfaces.InputFrequency

	if err := c.ShouldBindJSON(&frequencyInput); err != nil {
		c.JSON(erros.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if err := frequency.EditFrequencyService(c, frequencyInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("Response", "Event created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Frequency
// @Description Retrieve all Event
// @Tags Frequency
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Frequency "List of All Frequency"
// @Router /api/all-frequency [get]
func PullAllFrequency(c *gin.Context) {
	frequency.PullAllEventService(c)
}

// @Summary Retrieve one Frequency
// @Description Retrieve one Frequency
// @Tags Frequency
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Frequency "List of One Frequency"
// @Router /api/frequency/{id} [get]
func PullFrequencyById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}
	frequency.PullEventIdService(c, id)
}
