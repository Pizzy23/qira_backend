package frequency

import (
	"net/http"
	frequency "qira/internal/frequency/service"
	"qira/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Edit Frequency
// @Description Edit Frequency
// @Tags 3 - Frequency
// @Accept json
// @Produce json
// @Param id path int true "Threat Event ID"
// @Param request body interfaces.InputFrequency true "Edit Frequency"
// @Success 200 {object} db.Frequency "Your Frequency is by add"
// @Router /api/frequency/{id} [put]
func EditFrequency(c *gin.Context) {
	var frequencyInput interfaces.InputFrequency

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Parameters are invalid, need a Id")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&frequencyInput); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := frequency.EditFrequencyService(c, frequencyInput, id); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Event update successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Frequency
// @Description Retrieve all Event
// @Tags 3 - Frequency
// @Accept json
// @Produce json
// @Success 200 {object} []db.Frequency "List of All Frequency"
// @Router /api/all-frequency [get]
func PullAllFrequency(c *gin.Context) {
	frequency.PullAllEventService(c)
}

// @Summary Retrieve one Frequency
// @Description Retrieve one Frequency
// @Tags 3 - Frequency
// @Accept json
// @Produce json
// @Success 200 {object} db.Frequency "List of One Frequency"
// @Router /api/frequency/{id} [get]
func PullFrequencyById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", "Parameters are invalid, need a Id")
		c.Status(http.StatusInternalServerError)
		return
	}
	frequency.PullEventIdService(c, id)
}
