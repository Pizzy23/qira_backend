package inventory

import (
	"net/http"
	"qira/internal/interfaces"
	inventory "qira/internal/inventory/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create Asset
// @Description Create new Asset
// @Tags 1 - Inventory
// @Accept json
// @Produce json
// @Param request body interfaces.InputAssetsInventory true "Data for create new Asset"
// @Success 200 {object} db.AssetInventory "Asset Create"
// @Router /api/create-asset [post]
func CreateAsset(c *gin.Context) {
	var asset interfaces.InputAssetsInventory

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := inventory.CreateAssetService(c, asset); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Asset created successfully")
	c.Status(http.StatusOK)

}

// @Summary Retrieve All Assets
// @Description Retrieve all assets
// @Tags 1 - Inventory
// @Accept json
// @Produce json
// @Success 200 {object} []db.AssetInventory "List of All Assets"
// @Router /api/assets [get]
func PullAllAsset(c *gin.Context) {
	inventory.PullAllAsset(c)
}

// @Summary Retrieve Asset by ID
// @Description Retrieve an asset by its ID
// @Tags 1 - Inventory
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} db.AssetInventory "Asset Details"
// @Router /api/asset/{id} [get]
func PullAssetId(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}
	inventory.PullAssetId(c, id)

}

// @Summary Delete Asset
// @Description Update an existing Asset
// @Tags 1 - Inventory
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} db.AssetInventory "Asset Updated"
// @Router /api/asset/{id} [delete]
func DeleteAsset(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := inventory.DeleteAsset(c, id); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Asset delete successfully")
	c.Status(http.StatusOK)
}

// @Summary Update Asset
// @Description Update an existing Asset
// @Tags 1 - Inventory
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param request body interfaces.InputAssetsInventory true "Data to update Asset"
// @Success 200 {object} db.AssetInventory "Asset Updated"
// @Router /api/asset/{id} [put]
func UpdateAsset(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Set("Response", "Invalid ID")
		c.Status(http.StatusInternalServerError)
		return
	}

	var asset interfaces.InputAssetsInventory
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.Set("Response", "Parameters are invalid, need a JSON")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := inventory.UpdateAssetService(c, id, asset); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", "Asset updated successfully")
	c.Status(http.StatusOK)
}
