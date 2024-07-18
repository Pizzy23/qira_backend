package services

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullEventService(c *gin.Context) {
	var res []db.ThreatEventAssets
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &res); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", res)
	c.Status(http.StatusOK)
}

func CreateEventService(c *gin.Context, input interfaces.InputThreatEventAssets) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	for _, asset := range input.AffectedAsset {
		eventAsset := db.ThreatEventAssets{
			ThreatID:      input.ThreatID,
			ThreatEvent:   input.ThreatEvent,
			AffectedAsset: asset,
		}
		if err := db.Create(engine.(*xorm.Engine), &eventAsset); err != nil {
			return err
		}
	}
	return nil
}
