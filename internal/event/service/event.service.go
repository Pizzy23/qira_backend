package services

import (
	"errors"
	"net/http"
	"qira/db"

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

func CreateEventService(c *gin.Context, input db.ThreatEventAssets) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &input); err != nil {
		return err
	}
	return nil
}
