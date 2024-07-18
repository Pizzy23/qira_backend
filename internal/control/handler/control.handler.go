package control

import (
	"net/http"
	control "qira/internal/control/service"
	"qira/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

// @Summary Retrieve All Control
// @Description Retrieve all Event
// @Tags 7 - Control
// @Accept json
// @Produce json
// @Success 200 {object} db.ControlLibrary "List of All Event"
// @Router /api/all-control [get]
func PullAllControl(c *gin.Context) {
	control.PullAllControl(c)
}

// @Summary Retrieve Control by ID
// @Description Retrieve an Event by its ID
// @Tags 7 - Control
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} db.ControlLibrary "Event Details"
// @Router /api/control/{id} [get]
func PullControlId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	control.PullControlId(c, id)
}

// @Summary Create New Implementation
// @Description Create New Event Implementation
// @Tags 8 - Implementation
// @Accept json
// @Produce json
// @Param id path int true "Implementation Id"
// @Param request body interfaces.ImplementsInputNoID true "Data for create new Event"
// @Success 200 {object} db.ControlLibrary "List of All Assets"
// @Router /api/implementation [put]
func EditControlImplementation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		c.Set("Response", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	var implement interfaces.ImplementsInputNoID

	if err := c.ShouldBindJSON(&implement); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := control.CreateImplementService(c, implement, id); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Event created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Implementation
// @Description Retrieve all Implementation
// @Tags 8 - Implementation
// @Accept json
// @Produce json
// @Success 200 {object} db.ControlLibrary "List of All Implementation"
// @Router /api/all-implementation [get]
func PullAllControlImplementation(c *gin.Context) {
	control.PullAllImplements(c)
}

// @Summary Retrieve Implementation by ID
// @Description Retrieve an Implementation by its ID
// @Tags 8 - Implementation
// @Accept json
// @Produce json
// @Param id path int true "Implementation ID"
// @Success 200 {object} db.ControlLibrary "Implementation Details"
// @Router /api/implementation/{id} [get]
func PullControlImplementationId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", "Invalid asset ID")
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
// @Success 200 {object} []db.Control "List of All Control Strength"
// @Router /api/all-strength [get]
func PullAllControlStrength(c *gin.Context) {
	control.PullAllControlStrength(c)
}

// @Summary {WIP} Retrieve All Control Proposed
// @Description Retrieve all Control Proposed
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} []db.Propused "List of All Control Strength"
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
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	finalResults, err := control.CalculateAggregatedControlStrength(engine.(*xorm.Engine))
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", finalResults)
	c.Status(http.StatusOK)
}

// @Summary Update Controll
// @Description Create Controll
// @Tags 8 - Implementation
// @Accept json
// @Produce json
// @Param id path int true "Implementation Id"
// @Param request body interfaces.ImplementsInputNoID true "Data for create new Event"
// @Success 200 {object} db.ControlLibrary "List of All Assets"
// @Router /api/control/{id} [put]
func UpdateControl(c *gin.Context) {
	var controlInput interfaces.InputControlLibrary
	if err := c.ShouldBindJSON(&controlInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controlIDStr := c.Param("controlID")
	controlID, err := strconv.ParseInt(controlIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid control ID"})
		return
	}

	if err := control.UpdateControlService(c, controlID, controlInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "control updated successfully"})
}
