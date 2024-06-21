package control

import (
	"net/http"
	control "qira/internal/control/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve All Relevance
// @Description Retrieve all Relevance
// @Tags Control
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Relevance "List of All Event"
// @Router /api/all-relevance [get]
func PullRelevanceControl(c *gin.Context) {
	control.GetControl(c, "Relevance")
}

// @Summary Retrieve Relevance by ID
// @Description Retrieve an Relevance by its ID
// @Tags Control
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Relevance "Event Details"
// @Router /api/relevance/{id} [get]
func PullRelevanceId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Error", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	id64 := int64(id)
	control.GetById(c, "Relevance", id64)
}

// @Summary Retrieve All Implementation
// @Description Retrieve all Implementation
// @Tags Control
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.ControlImplementation "List of All Implementation"
// @Router /api/all-implementation [get]
func PullImplementationControl(c *gin.Context) {
	control.GetControl(c, "Implementation")
}

// @Summary Retrieve Implementation by ID
// @Description Retrieve an Implementation by its ID
// @Tags Control
// @Accept json
// @Produce json
// @Param id path int true "Implementation ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.ControlImplementation "Implementation Details"
// @Router /api/implementation/{id} [get]
func PullImplementationId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Error", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	id64 := int64(id)
	control.GetById(c, "Implementation", id64)
}

// @Summary Retrieve All Propused
// @Description Retrieve all Propused
// @Tags Control
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Propused "List of All Propused"
// @Router /api/all-propused [get]
func PullPropusedControl(c *gin.Context) {
	control.GetControl(c, "Propused")
}

// @Summary Retrieve Propused by ID
// @Description Retrieve an Propused by its ID
// @Tags Control
// @Accept json
// @Produce json
// @Param id path int true "Propused ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Propused "Propused Details"
// @Router /api/propused/{id} [get]
func PullPropusedId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}
	id64 := int64(id)
	control.GetById(c, "Propused", id64)
}

// @Summary Retrieve All library
// @Description Retrieve all library
// @Tags Control
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.ControlLibrary "List of All Event"
// @Router /api/all-library [get]
func PullLibraryControl(c *gin.Context) {
	control.GetControl(c, "Library")
}

// @Summary Retrieve Library by ID
// @Description Retrieve an library by its ID
// @Tags Control
// @Accept json
// @Produce json
// @Param id path int true "library ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.ControlLibrary "Event Details"
// @Router /api/library/{id} [get]
func PullLibraryId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Error", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	id64 := int64(id)
	control.GetById(c, "Library", id64)
}

// @Summary Retrieve All Strength
// @Description Retrieve all Strength
// @Tags Control
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Strength "List of All Event"
// @Router /api/all-strength [get]
func PullStrengthControl(c *gin.Context) {
	control.GetControl(c, "Strength")
}

// @Summary Retrieve Strength by ID
// @Description Retrieve an Event by its ID
// @Tags Control
// @Accept json
// @Produce json
// @Param id path int true "Strength ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Strength "Event Details"
// @Router /api/strength/{id} [get]
func PullStrengthId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Error", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}
	id64 := int64(id)
	control.GetById(c, "Strength", id64)
}
