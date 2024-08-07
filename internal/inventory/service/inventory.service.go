package inventory

import (
	"encoding/json"
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateAssetService(c *gin.Context, asset interfaces.InputAssetsInventory) error {
	newAsset := db.AssetInventory{
		Name:                    asset.Name,
		Description:             asset.Description,
		Location:                asset.Location,
		Responsible:             asset.Responsible,
		BusinessValue:           asset.BusinessValue,
		ReplacementCost:         asset.ReplacementCost,
		Criticality:             asset.Criticality,
		Users:                   asset.Users,
		RoleInTargetEnvironment: asset.RoleInTargetEnvironment,
	}
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &newAsset); err != nil {
		return err
	}
	return nil

}

func PullAllAsset(c *gin.Context) {
	var assets []db.AssetInventory
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &assets); err != nil {
		c.Set("Response", "Error retrieving assets: "+err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if assets == nil {
		var empty []string
		c.Set("Response", empty)
		c.Status(http.StatusOK)
		return
	}
	c.Set("Response", assets)
	c.Status(http.StatusOK)
}

func PullAssetId(c *gin.Context, id int) {
	var asset db.AssetInventory
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &asset, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving asset")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "Asset not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", asset)
	c.Status(http.StatusOK)
}

func UpdateAssetService(c *gin.Context, id int64, asset interfaces.InputAssetsInventory) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	assetToUpdate := db.AssetInventory{
		Name:                    asset.Name,
		Description:             asset.Description,
		Location:                asset.Location,
		Responsible:             asset.Responsible,
		BusinessValue:           asset.BusinessValue,
		ReplacementCost:         asset.ReplacementCost,
		Criticality:             asset.Criticality,
		Users:                   asset.Users,
		RoleInTargetEnvironment: asset.RoleInTargetEnvironment,
	}

	if err := db.UpdateByID(engine.(*xorm.Engine), &assetToUpdate, id); err != nil {
		return err
	}
	return nil
}

func DeleteAsset(c *gin.Context, id int64) error {
	var asset db.AssetInventory

	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	has, err := engine.(*xorm.Engine).ID(id).Get(&asset)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("asset not found")
	}

	_, err = engine.(*xorm.Engine).ID(id).Delete(&asset)
	if err != nil {
		return err
	}

	events := make([]db.ThreatEventAssets, 0)
	err = engine.(*xorm.Engine).Find(&events)
	if err != nil {
		return err
	}

	for _, event := range events {
		affectedAssets := make([]string, 0)
		err = json.Unmarshal([]byte(event.AffectedAsset), &affectedAssets)
		if err != nil {
			return err
		}

		updatedAssets := make([]string, 0)
		for _, assetName := range affectedAssets {
			if assetName != asset.Name {
				updatedAssets = append(updatedAssets, assetName)
			}
		}

		if len(updatedAssets) > 0 {
			updatedAssetsJSON, err := json.Marshal(updatedAssets)
			if err != nil {
				return err
			}
			event.AffectedAsset = string(updatedAssetsJSON)
			_, err = engine.(*xorm.Engine).ID(event.ID).Update(&event)
			if err != nil {
				return err
			}
		} else {
			_, err = engine.(*xorm.Engine).ID(event.ID).Delete(&event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
