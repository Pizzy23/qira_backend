package catalogue

import (
	"errors"
	"net/http"
	"qira/db"
	frequency "qira/internal/frequency/service"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateEventService(c *gin.Context, event interfaces.InputThreatEventCatalogue) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if _, err := engine.(*xorm.Engine).Insert(&event); err != nil {
		return err
	}

	errChan := make(chan error)

	go func(eventID int, eventName string) {
		frequencyInput := interfaces.InputFrequency{
			ThreatEventID:       eventID,
			ThreatEvent:         eventName,
			MinFrequency:        0,
			MaxFrequency:        0,
			MostCommonFrequency: 0,
			SupportInformation:  "",
		}

		if err := frequency.CreateFrequencyService(c, frequencyInput); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}(event.ID, event.ThreatEvent)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func PullAllEventService(c *gin.Context) {
	var events []db.ThreatEventCatalogue
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &events); err != nil {
		c.Set("Error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", events)
	c.Status(http.StatusOK)
}

func PullEventIdService(c *gin.Context, id int) {
	var event db.ThreatEventCatalogue
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &event, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving event")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Event not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", event)
	c.Status(http.StatusOK)
}
