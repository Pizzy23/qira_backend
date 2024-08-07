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
		ThreatGroup: event.ThreatGroup,
		ThreatEvent: event.ThreatEvent,
		Description: event.Description,
		InScope:     event.InScope,
	}

	if _, err := engine.(*xorm.Engine).Insert(&eventDB); err != nil {
		return err
	}

	eventID := eventDB.ID

	RiskController := db.RiskController{
		Name: event.ThreatEvent,
	}

	if err := db.Create(engine.(*xorm.Engine), RiskController); err != nil {
		return err
	}

	frequencyInput := db.Frequency{
		ThreatEventID: eventID,
		ThreatEvent:   event.ThreatEvent,
		MinFrequency:  0,
		MaxFrequency:  0,
	}
	EventAssets := db.ThreatEventAssets{
		ThreatID:    eventID,
		ThreatEvent: event.ThreatEvent,
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

	// Obter o evento de ameaça para usar seu nome no filtro de relevância
	var threatEvent db.ThreatEventCatalog
	has, err := engine.(*xorm.Engine).ID(eventID).Get(&threatEvent)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("threat event not found")
	}

	if _, err := engine.(*xorm.Engine).Where("threat_event_id = ?", eventID).Delete(&db.Frequency{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).Where("threat_event_id = ?", eventID).Delete(&db.LossHigh{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).Where("threat_id = ?", eventID).Delete(&db.ThreatEventAssets{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).Where("threat_event_id = ?", eventID).Delete(&db.RiskCalculation{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).Where("type_of_attack = ?", threatEvent.ThreatEvent).Delete(&db.Relevance{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).ID(eventID).Delete(&db.ThreatEventCatalog{}); err != nil {
		return err
	}

	return nil
}
