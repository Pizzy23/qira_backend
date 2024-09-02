package catalogue

import (
	"errors"
	"net/http"
	"qira/db"
	frequency "qira/internal/frequency/service"
	"qira/internal/interfaces"
	"qira/util"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateEventService(c *gin.Context, event interfaces.InputThreatEventCatalogue) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	event = util.SanitizeInputCatalogue(&event)
	eventDB := db.ThreatEventCatalog{
		ThreatGroup: util.CleanString(event.ThreatGroup),
		ThreatEvent: util.CleanString(event.ThreatEvent),
		Description: event.Description,
		InScope:     event.InScope,
	}

	if _, err := engine.(*xorm.Engine).Insert(&eventDB); err != nil {
		return err
	}

	eventID := eventDB.ID

	RiskController := db.RiskController{
		Name: util.CleanString(event.ThreatEvent),
	}

	if err := db.Create(engine.(*xorm.Engine), RiskController); err != nil {
		return err
	}

	frequencyInput := db.Frequency{
		ThreatEventID: eventID,
		ThreatEvent:   util.CleanString(event.ThreatEvent),
		MinFrequency:  0,
		MaxFrequency:  0,
	}
	EventAssets := db.ThreatEventAssets{
		ThreatID:    eventID,
		ThreatEvent: util.CleanString(event.ThreatEvent),
	}

	if err := frequency.CreateFrequencyService(c, frequencyInput); err != nil {
		return err
	}

	if err := db.Create(engine.(*xorm.Engine), &EventAssets); err != nil {
		return err
	}

	return nil
}

func PullAllEventService(c *gin.Context) {
	var events []db.ThreatEventCatalog
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &events); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", events)
	c.Status(http.StatusOK)
}

func PullEventIdService(c *gin.Context, id int) {
	var event db.ThreatEventCatalog
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &event, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving event")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "Event not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", event)
	c.Status(http.StatusOK)
}

func DeleteEventService(c *gin.Context, eventID int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var threatEvent db.ThreatEventCatalog
	found, err := db.GetByID(engine.(*xorm.Engine), &threatEvent, eventID)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("threat event not found")
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.Frequency{}, map[string]interface{}{"threat_event_id": eventID}); err != nil {
		return err
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.LossHigh{}, map[string]interface{}{"threat_event_id": eventID}); err != nil {
		return err
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.LossHighGranular{}, map[string]interface{}{"threat_event_id": eventID}); err != nil {
		return err
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.LossHighTotal{}, map[string]interface{}{"threat_event_id": eventID}); err != nil {
		return err
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.ThreatEventAssets{}, map[string]interface{}{"threat_id": eventID}); err != nil {
		return err
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.RiskCalculation{}, map[string]interface{}{"threat_event_id": eventID}); err != nil {
		return err
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.Relevance{}, map[string]interface{}{"type_of_attack": threatEvent.ThreatEvent}); err != nil {
		return err
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.ThreatEventCatalog{}, map[string]interface{}{"id": eventID}); err != nil {
		return err
	}

	return nil
}

func UpdateEventService(c *gin.Context, id int, updatedEvent interfaces.InputThreatEventCatalogue) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var event db.ThreatEventCatalog
	found, err := db.GetByID(engine.(*xorm.Engine), &event, int64(id))
	if err != nil {
		return err
	}
	if !found {
		return errors.New("event not found")
	}

	event.ThreatGroup = util.CleanString(updatedEvent.ThreatGroup)
	event.ThreatEvent = util.CleanString(updatedEvent.ThreatEvent)
	event.Description = updatedEvent.Description
	event.InScope = updatedEvent.InScope

	if err := db.UpdateByCatalogue(engine.(*xorm.Engine), &event, int64(id)); err != nil {
		return err
	}

	return nil
}
