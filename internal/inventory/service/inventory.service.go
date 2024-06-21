package inventory

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateAssetService(c *gin.Context, asset db.AssetsInventory) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &asset); err != nil {
		return err
	}
	return nil

}

func PullAllAsset(c *gin.Context) {
	var assets []interfaces.AssetsInventory
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &assets); err != nil {
		c.Set("Error", "Error")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", assets)
	c.Status(http.StatusOK)
}

func PullAssetId(c *gin.Context, id int) {
	var asset interfaces.AssetsInventory
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &asset, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving asset")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Asset not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", asset)
	c.Status(http.StatusOK)
}
