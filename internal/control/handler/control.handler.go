package control

import (
	"net/http"
	control "qira/internal/control/service"
	"qira/internal/interfaces"
	erros "qira/middleware/interfaces/errors"

	"github.com/gin-gonic/gin"
)

// @Summary Create Control Relevance
// @Description Create new Control Relevance
// @Tags Control
// @Accept json
// @Produce json
// @Param request body interfaces.InputControlls true "Data for create new Risk"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Relevance "Risk Create"
// @Router /api/create-relevance [post]
func CreateRelevance(c *gin.Context) {
	var relevance interfaces.InputControlls
	relevance.ControlType = "Relevance"
	if err := c.ShouldBindJSON(&relevance); err != nil {
		c.JSON(erros.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if err := control.CreateRelevanceService(c, relevance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("Response", "Control Relevance created successfully")
	c.Status(http.StatusOK)

}

// @Summary Create Control Strength
// @Description Create new Control Strength
// @Tags Control
// @Accept json
// @Produce json
// @Param request body interfaces.InputControlls true "Data for create new Risk"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Strength "Risk Create"
// @Router /api/create-strength [post]
func CreateStrength(c *gin.Context) {
	var strength interfaces.InputControlls
	strength.ControlType = "Strength"
	if err := c.ShouldBindJSON(&strength); err != nil {
		c.Set("Error", "Invalid asset ID")
		c.Status(http.StatusBadRequest)
		return
	}

	if err := control.CreateRelevanceService(c, strength); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Set("Response", "Control Strength created successfully")
	c.Status(http.StatusOK)

}

// @Summary Create Control Propused
// @Description Create new Control Propused
// @Tags Control
// @Accept json
// @Produce json
// @Param request body interfaces.InputControlls true "Data for create new Risk"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.Propused "Risk Create"
// @Router /api/create-propused [post]
func CreatePropused(c *gin.Context) {
	var propused interfaces.InputControlls
	propused.ControlType = "Propused"
	if err := c.ShouldBindJSON(&propused); err != nil {
		c.Set("Error", "Parameters are invalid, need a JSON")
		c.Status(http.StatusBadRequest)
		return
	}

	if err := control.CreateRelevanceService(c, propused); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Set("Response", "Control Propused created successfully")
	c.Status(http.StatusOK)

}

// @Summary Create Control Library
// @Description Create new Control Library
// @Tags Control
// @Accept json
// @Produce json
// @Param request body interfaces.InputControlLibrary true "Data for create new Risk"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.ControlLibrary "Risk Create"
// @Router /api/create-library [post]
func CreateLibrary(c *gin.Context) {
	var library interfaces.InputControlLibrary
	if err := c.ShouldBindJSON(&library); err != nil {
		c.Set("Error", "Parameters are invalid, need a JSON")
		c.Status(http.StatusBadRequest)
		return
	}

	if err := control.CreateLibraryService(c, library); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Set("Response", "Control library created successfully")
	c.Status(http.StatusOK)

}

// @Summary Create Control Implementation
// @Description Create new Control Implementation
// @Tags Control
// @Accept json
// @Produce json
// @Param request body interfaces.InputControlImplementation true "Data for create new Risk"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.ControlImplementation "Risk Create"
// @Router /api/create-implementation [post]
func CreateImplementation(c *gin.Context) {
	var implementation interfaces.InputControlImplementation

	if err := c.ShouldBindJSON(&implementation); err != nil {
		c.Set("Error", "Parameters are invalid, need a JSON")
		c.Status(http.StatusBadRequest)
		return
	}

	if err := control.CreateImplementationService(c, implementation); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Set("Response", "Control implementation created successfully")
	c.Status(http.StatusOK)

}
