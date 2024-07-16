package control

import (
	"net/http"
	control "qira/internal/control/service"
	"qira/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

// @Summary Create New Control
// @Description Create New Event Control
// @Tags Control
// @Accept json
// @Produce json
// @Param request body interfaces.InputControlLibrary true "Data for create new Event"
// @Success 200 {object} db.ControlLibrary "List of All Assets"
// @Router /api/control [post]
func CreateControl(c *gin.Context) {
	var controlInput interfaces.InputControlLibrary

	if err := c.ShouldBindJSON(&controlInput); err != nil {
		c.Set("Error", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := c.ShouldBindJSON(&controlInput); err != nil {

		return
	}
	if err := control.CreateControlService(c, controlInput); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Event created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Control
// @Description Retrieve all Event
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} db.ControlLibrary "List of All Event"
// @Router /api/all-control [get]
func PullAllControl(c *gin.Context) {
	control.PullAllControl(c)
}

// @Summary Retrieve Control by ID
// @Description Retrieve an Event by its ID
// @Tags Control
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} db.ControlLibrary "Event Details"
// @Router /api/control/{id} [get]
func PullControlId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Error", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	control.PullControlId(c, id)
}

// @Summary Create New Implementation
// @Description Create New Event Implementation
// @Tags Control
// @Accept json
// @Produce json
// @Param request body interfaces.ImplementsInput true "Data for create new Event"
// @Success 200 {object} db.ControlLibrary "List of All Assets"
// @Router /api/implementation [post]
func CreateControlImplementation(c *gin.Context) {
	var controlInput interfaces.ImplementsInput

	if err := c.ShouldBindJSON(&controlInput); err != nil {
		c.Set("Error", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := c.ShouldBindJSON(&controlInput); err != nil {

		return
	}
	if err := control.CreateImplementService(c, controlInput); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Event created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Implementation
// @Description Retrieve all Implementation
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} db.ControlLibrary "List of All Implementation"
// @Router /api/all-implementation [get]
func PullAllControlImplementation(c *gin.Context) {
	control.PullAllControl(c)
}

// @Summary Retrieve Implementation by ID
// @Description Retrieve an Implementation by its ID
// @Tags Control
// @Accept json
// @Produce json
// @Param id path int true "Implementation ID"
// @Success 200 {object} db.ControlLibrary "Implementation Details"
// @Router /api/implementation/{id} [get]
func PullControlImplementationId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Error", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	control.PullControlId(c, id)
}

// @Summary {WIP} Retrieve All Control Strength
// @Description Retrieve all Control Strength
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} []db.ControlDinamic "List of All Control Strength"
// @Router /api/all-strength [get]
func PullAllControlStrength(c *gin.Context) {
	control.PullAllControlStrength(c)
}

// @Summary {WIP} Retrieve All Control Proposed
// @Description Retrieve all Control Proposed
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} []db.PropusedDinamic "List of All Control Strength"
// @Router /api/all-proposed [get]
func PullAllControlProposed(c *gin.Context) {
	control.PullAllControlProposed(c)
}

// @Summary {WIP} Retrieve Aggregated Control Strength
// @Description Retrieve aggregated control strength for all threat events
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} []db.AggregatedStrength "List of Aggregated Control Strength"
// @Router /api/aggregated-control-strength [get]
func PullAggregatedControlStrength(c *gin.Context) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	finalResults, err := control.CalculateAggregatedControlStrength(engine.(*xorm.Engine))
	if err != nil {
		c.Set("Error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", finalResults)
	c.Status(http.StatusOK)
}
