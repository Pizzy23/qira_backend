package event

import (
	"net/http"
	"qira/db"
	event "qira/internal/event/service"

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
// @Param request body db.ThreatEventAssets true "Data for create new Event"
// @Success 200 {object} db.ThreatEventAssets "Event created successfully"
// @Router /api/event [post]
func CreateEvent(c *gin.Context) {
	var eventCatalogue db.ThreatEventAssets

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
