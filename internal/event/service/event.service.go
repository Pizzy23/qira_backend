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

	for _, asset := range input.AffectedAsset {
		var eventAsset db.ThreatEventAssets
		has, err := engine.(*xorm.Engine).Where("threat_i_d = ? AND affected_asset = ?", id, asset).Get(&eventAsset)
		if err != nil {
			return err
		}

		if has {
			eventAsset.ThreatEvent = input.ThreatEvent
			if _, err := engine.(*xorm.Engine).ID(eventAsset.ID).Update(&eventAsset); err != nil {
				return err
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
	return nil
}
