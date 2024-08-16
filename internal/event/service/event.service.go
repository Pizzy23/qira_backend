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
	var threatEvents []db.ThreatEventCatalog
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.InScope(engine.(*xorm.Engine).NewSession(), &threatEvents); err != nil {
		c.Set("Response", "Error fetching threat events")
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, threatEvent := range threatEvents {
		var existingEvent db.ThreatEventAssets
		found, err := db.GetByEventID(engine.(*xorm.Engine), &existingEvent, threatEvent.ID)
		if err != nil {
			c.Set("Response", "Error checking event existence")
			c.Status(http.StatusInternalServerError)
			return
		}

		if !found {
			newEvent := db.ThreatEventAssets{
				ThreatID:      threatEvent.ID,
				ThreatEvent:   threatEvent.ThreatEvent,
				AffectedAsset: "Generic asset",
			}

			if err := db.Create(engine.(*xorm.Engine), &newEvent); err != nil {
				c.Set("Response", "Error inserting new threat event with generic asset")
				c.Status(http.StatusInternalServerError)
				return
			}
		}
	}

	if err := db.GetAll(engine.(*xorm.Engine), &res); err != nil {
		c.Set("Response", "Error fetching updated threat event assets")
		c.Status(http.StatusInternalServerError)
		return
	}

	eventMap := make(map[int64]*interfaces.OutPutThreatEventAssets)
	for _, item := range res {
		if _, exists := eventMap[item.ThreatID]; !exists {
			eventMap[item.ThreatID] = &interfaces.OutPutThreatEventAssets{
				ThreatID:      item.ThreatID,
				ThreatEvent:   item.ThreatEvent,
				AffectedAsset: []string{},
			}
		}
		eventMap[item.ThreatID].AffectedAsset = append(eventMap[item.ThreatID].AffectedAsset, item.AffectedAsset)
	}

	var output []interfaces.OutPutThreatEventAssets
	for _, value := range eventMap {
		output = append(output, *value)
	}

	c.Set("Response", output)
	c.Status(http.StatusOK)
}

func CreateEventService(c *gin.Context, input interfaces.InputThreatEventAssets, id int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingAssets []db.ThreatEventAssets
	if err := engine.(*xorm.Engine).Where("threat_id = ?", id).Find(&existingAssets); err != nil {
		return err
	}

	existingMap := make(map[string]db.ThreatEventAssets)
	for _, ea := range existingAssets {
		existingMap[ea.AffectedAsset] = ea
	}

	newAssetsMap := make(map[string]bool)
	for _, asset := range input.AffectedAsset {
		newAssetsMap[asset] = true

		if eventAsset, exists := existingMap[asset]; exists {
			if eventAsset.ThreatEvent != input.ThreatEvent {
				eventAsset.ThreatEvent = input.ThreatEvent
				if _, err := engine.(*xorm.Engine).ID(eventAsset.ID).Update(&eventAsset); err != nil {
					return err
				}
			}
		} else {
			newEventAsset := db.ThreatEventAssets{
				ThreatID:      id,
				ThreatEvent:   input.ThreatEvent,
				AffectedAsset: asset,
			}
			if _, err := engine.(*xorm.Engine).Insert(&newEventAsset); err != nil {
				return err
			}
		}
	}

	for key, eventAsset := range existingMap {
		if !newAssetsMap[key] {
			if _, err := engine.(*xorm.Engine).ID(eventAsset.ID).Delete(&db.ThreatEventAssets{}); err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteEvent(c *gin.Context, id int) error {
	var asset db.ThreatEventAssets
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	has, err := engine.(*xorm.Engine).ID(id).Get(&asset)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("Threat Event Assets not found")
	}

	_, err = engine.(*xorm.Engine).ID(id).Delete(&asset)
	if err != nil {
		return err
	}

	return nil
}
