package frequency

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateFrequencyService(c *gin.Context, data db.Frequency) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &data); err != nil {
		return err
	}
	return nil
}

func EditFrequencyService(c *gin.Context, freq interfaces.InputFrequency) error {
	var frequencyTable *db.Frequency
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Update(engine.(*xorm.Engine), frequencyTable, &freq); err != nil {
		return err
	}
	return nil
}

func PullAllEventService(c *gin.Context) {
	var frequency []interfaces.InputFrequency
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &frequency); err != nil {
		c.Set("Error", "Error")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", frequency)
	c.Status(http.StatusOK)
}

func PullEventIdService(c *gin.Context, id int) {
	var frequency interfaces.InputFrequency
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &frequency, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving Frequency")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Frequency not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", frequency)
	c.Status(http.StatusOK)
}
