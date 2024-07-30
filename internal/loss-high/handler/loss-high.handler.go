package losshigh

import (
	"net/http"
	"qira/internal/interfaces"
	losshigh "qira/internal/loss-high/service"
	"qira/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create LossHigh
// @Description Create new LossHigh
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Param id path int true "Threat Event ID"
// @Param request body interfaces.InputLossHigh true "Data for create new LossHigh"
// @Success 200 {object} db.LossHigh "LossHigh Create"
// @Router /api/update-losshigh/{id} [put]
func CreateLossHigh(c *gin.Context) {
	var LossHigh interfaces.InputLossHigh

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&LossHigh); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := losshigh.CreateLossHighService(c, LossHigh, id); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "LossHigh created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All LossHigh
// @Description Retrieve and aggregate all LossHigh records
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Success 200 {object} []db.LossHigh "List of All LossHigh with Aggregated Data"
// @Router /api/losshigh [get]
func PullAllLossHigh(c *gin.Context) {
	aggregatedLosses, err := losshigh.GetAggregatedLosses(c)
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", aggregatedLosses)
	c.Status(http.StatusOK)
}

// @Summary Create LossHigh Singular
// @Description Create new LossHigh
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Param id path int true "Threat Event ID"
// @Param request body interfaces.InputLossHigh true "Data for create new LossHigh"
// @Success 200 {object} db.LossHigh "LossHigh Create"
// @Router /api/update-losshigh-singular/{id} [put]
func CreateLossHighSingular(c *gin.Context) {
	var LossHigh interfaces.InputLossHigh

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&LossHigh); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := losshigh.CreateSingularLossService(c, LossHigh, id); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "LossHigh created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All LossHigh Singular
// @Description Retrieve and aggregate all LossHigh records
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Success 200 {object} []db.LossHigh "List of All LossHigh with Aggregated Data"
// @Router /api/losshigh-singular [get]
func PullAllLossHighSingular(c *gin.Context) {
	aggregatedLosses, err := losshigh.GetSingularLosses(c)
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", aggregatedLosses)
	c.Status(http.StatusOK)
}

// @Summary Create LossHigh Granuled
// @Description Create new LossHigh
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Param id path int true "Threat Event ID"
// @Param request body interfaces.InputLossHighGranulade true "Data for create new LossHigh"
// @Success 200 {object} db.LossHigh "LossHigh Create"
// @Router /api/update-losshigh-granuled/{id} [put]
func CreateLossHighGranuled(c *gin.Context) {
	var LossHigh interfaces.InputLossHighGranulade

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&LossHigh); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := losshigh.CreateLossHighGranularService(c, LossHigh, id); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "LossHigh created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All LossHigh Granuled
// @Description Retrieve and aggregate all LossHigh records
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Success 200 {object} []db.LossHighGranular "List of All LossHigh with Aggregated Data"
// @Router /api/losshigh-granuled [get]
func PullAllLossHighGranuled(c *gin.Context) {
	aggregatedLosses, err := losshigh.GetGranularLosses(c)
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", aggregatedLosses)
	c.Status(http.StatusOK)
}

// @Summary Retrieve LossHigh by ID
// @Description Retrieve an LossHigh by its ID
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Param id path int true "LossHigh ID"
// @Success 200 {object} db.LossHigh "LossHigh Details"
// @Router /api/losshigh/{id} [get]
func PullLosstId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	losshigh.PullLossHighId(c, id)
}

// @Summary Create LossHigh Specific
// @Description Create new LossHigh Specific
// @Tags 5 - Loss-High
// @Accept json
// @Produce json
// @Param Loss header string true "Tipo de loss" Enums("Singular","LossHigh","Granular")
// @Success 200 {object} db.LossHigh "LossHigh Create"
// @Router /api/losshigh-specific [post]
func CreateLossHighSpecific(c *gin.Context) {
	typeLoss := c.GetHeader("Loss")
	vali := util.EnumLoss(typeLoss)
	if !vali {
		c.Set("Response", "TypeLoss its not valid use: `Singular,LossHigh,Granular`")
		c.Status(http.StatusInternalServerError)
	}
	losshigh.CreateLossSpecific(c, typeLoss)
}
