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
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
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
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := implementation.CreateImplementsService(c, ImplementationInput); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Implements created successfully")
	c.Status(http.StatusOK)

}
