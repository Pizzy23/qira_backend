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
	var allInOne interfaces.InputThreatEventAndAsset
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	eventChan := make(chan interfaces.InputThreatEventAssets)
	assetChan := make(chan interfaces.InputAssetsInventory)
	errChan := make(chan error, 2)

	go func() {
		var event interfaces.InputThreatEventAssets
		if err := db.Read(engine.(*xorm.Engine), &event, nil); err != nil {
			errChan <- err
		} else {
			eventChan <- event
			errChan <- nil
		}
	}()

	go func() {
		var asset interfaces.InputAssetsInventory
		if err := db.Read(engine.(*xorm.Engine), &asset, nil); err != nil {
			errChan <- err
		} else {
			assetChan <- asset
			errChan <- nil
		}
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			c.Set("Error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	allInOne.Event = <-eventChan
	allInOne.Asset = <-assetChan

	c.Set("Response", allInOne)
	c.Status(http.StatusOK)
}

func CreateEventService(c *gin.Context, input db.ThreatEventAssets) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	event := interfaces.ThreatEventAssets{
		ThreatID:      input.ThreatID,
		ThreatEvent:   input.ThreatEvent,
		AffectedAsset: input.AffectedAsset,
	}

	errChan := make(chan error, 1)

	go func() {
		if err := db.Create(engine.(*xorm.Engine), &event); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}
