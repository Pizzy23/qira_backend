package implementation

import (
	"net/http"
	implementation "qira/internal/implemention/service"
	"qira/internal/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve All Implementss
// @Description Retrieve all Implementss
// @Tags Implements
// @Accept json
// @Produce json
// @Success 200 {object} []db.Implements "List of All Implementss"
// @Router /api/implements [get]
func PullAllImplements(c *gin.Context) {
	implementation.PullAllImplements(c)
}

// @Summary Retrieve Implements by ID
// @Description Retrieve an Implements by its ID
// @Tags Implements
// @Accept json
// @Produce json
// @Param id path int true "Implements ID"
// @Success 200 {object} db.Implements "Implements Details"
// @Router /api/implements/{id} [get]
func PullImplementsId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Implements ID"})
		return
	}
	implementation.PullImplementsId(c, id)
}

// @Summary Create Implements
// @Description Create new Implements
// @Tags Revelance
// @Accept json
// @Produce json
// @Param request body interfaces.ImplementsInput true "Data for create new Implements"
// @Success 200 {object} db.Implements "Implements Create"
// @Router /api/create-implements [post]
func CreateImplements(c *gin.Context) {
	var ImplementationInput interfaces.ImplementsInput

	if err := c.ShouldBindJSON(&ImplementationInput); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if err := implementation.CreateImplementsService(c, ImplementationInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("Response", "Implements created successfully")
	c.Status(http.StatusOK)

}
