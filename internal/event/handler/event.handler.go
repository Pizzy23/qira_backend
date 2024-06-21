package event

import (
	"net/http"
	"qira/db"
	event "qira/internal/event/service"

	"github.com/gin-gonic/gin"
)

// @Summary Event
// @Description Event
// @Tags Event
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.ThreatEventAssets "Your Frequency is by add"
// @Router /api/all-event [get]
func PullAllForEvent(c *gin.Context) {
	event.PullEventService(c)
}

// @Summary Create Event
// @Description Create Event
// @Tags Event
// @Accept json
// @Produce json
// @Param request body interfaces.db.ThreatEventAssets true "Data for create new Event"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} db.ThreatEventAssets "Event created successfully"
// @Router /api/event [post]
func CreateEvent(c *gin.Context) {
	var eventCatalogue db.ThreatEventAssets

	if err := c.ShouldBindJSON(&eventCatalogue); err != nil {
		c.Set("Error", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := event.CreateEventService(c, eventCatalogue); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Event created successfully")
	c.Status(http.StatusOK)
}
