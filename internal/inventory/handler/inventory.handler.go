package inventory

import (
	"net/http"
	"qira/db"
	inventory "qira/internal/inventory/service"
	erros "qira/middleware/interfaces/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create Asset
// @Description Create new Asset
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body db.AssetsInventory true "Data for create new Asset"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} db.AssetsInventory "Asset Create"
// @Router /api/create-asset [post]
func CreateAsset(c *gin.Context) {
	var asset db.AssetsInventory

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(erros.StatusNotAcceptable, gin.H{"error": "Parameters are invalid, need a JSON"})
		return
	}

	if err := inventory.CreateAssetService(c, asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("Response", "Asset created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Assets
// @Description Retrieve all assets
// @Tags Inventory
// @Accept json
// @Produce json
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} []interfaces.AssetsInventory "List of All Assets"
// @Router /api/assets [get]
func PullAllAsset(c *gin.Context) {
	inventory.PullAllAsset(c)
}

// @Summary Retrieve Asset by ID
// @Description Retrieve an asset by its ID
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param Authorization header string true "Auth Token" default(Bearer <token>)
// @Success 200 {object} interfaces.AssetsInventory "Asset Details"
// @Router /api/asset/{id} [get]
func PullAssetId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}
	inventory.PullAssetId(c, id)

}
