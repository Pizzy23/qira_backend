package event

import (
	"net/http"
	event "qira/internal/event/service"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// @Summary Event
// @Description Event
// @Tags 4 - Event
// @Accept json
// @Produce json
// @Success 200 {object} db.ThreatEventAssets "Your Frequency is by add"
// @Router /api/all-event [get]
func PullAllForEvent(c *gin.Context) {
	event.PullEventService(c)
}

// @Summary Create Event
// @Description Create Event
// @Tags 4 - Event
// @Accept json
// @Produce json
// @Param request body interfaces.InputThreatEventAssets true "Data for create new Event"
// @Success 200 {object} db.ThreatEventAssets "Event created successfully"
// @Router /api/event/{id} [put]
func CreateEvent(c *gin.Context) {
	var eventCatalogue interfaces.InputThreatEventAssets

	if err := c.ShouldBindJSON(&eventCatalogue); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := event.CreateEventService(c, eventCatalogue); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Event created successfully")
	c.Status(http.StatusOK)
}
