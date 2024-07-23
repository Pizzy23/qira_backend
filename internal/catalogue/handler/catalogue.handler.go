package catalogue

import (
	"net/http"
	catalogue "qira/internal/catalogue/service"
	"qira/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary {WIP} Create New Catalogue
// @Description Create New Event Catalogue
// @Tags 2 - Catalogue
// @Accept json
// @Produce json
// @Param request body interfaces.InputThreatEventCatalogue true "Data for create new Event"
// @Success 200 {object} db.ThreatEventCatalog "List of All Events catalogues"
// @Router /api/catalogue [post]
func CreateEvent(c *gin.Context) {
	var eventCatalogue interfaces.InputThreatEventCatalogue

	if err := c.ShouldBindJSON(&eventCatalogue); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := catalogue.CreateEventService(c, eventCatalogue); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Event created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Catalogue
// @Description Retrieve all Event
// @Tags 2 - Catalogue
// @Accept json
// @Produce json
// @Success 200 {object} db.ThreatEventCatalog "List of All Event"
// @Router /api/all-catalogue [get]
func PullAllEvent(c *gin.Context) {
	catalogue.PullAllEventService(c)
}

// @Summary Retrieve Catalogue by ID
// @Description Retrieve an Event by its ID
// @Tags 2 - Catalogue
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} db.ThreatEventCatalog "Event Details"
// @Router /api/catalogue/{id} [get]
func PullEventId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	catalogue.PullEventIdService(c, id)
}

// @Summary Delete Catalogue by ID
// @Description Delete an Event by its ID
// @Tags 2 - Catalogue
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} db.ThreatEventCatalog "Event Details"
// @Router /api/catalogue/{id} [delete]
func DeleteEventId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	catalogue.DeleteEventService(c, id)
}
