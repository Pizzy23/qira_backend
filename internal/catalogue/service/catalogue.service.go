package catalogue

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qira/db"
	frequency "qira/internal/frequency/service"
	"qira/internal/interfaces"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateEventService(c *gin.Context, event interfaces.InputThreatEventCatalogue) error {

	engine, ok := c.MustGet("db_engine").(*xorm.Engine)
	if !ok {
		return errors.New("database connection not found")
	}

	if err := createTables(engine, event.ThreatEvent); err != nil {
		return err
	}

	if _, err := engine.Insert(&event); err != nil {
		return err
	}

	errChan := make(chan error)

	go func(eventID int64, eventName string) {
		frequencyInput := db.Frequency{
			ThreatEventID: eventID,
			ThreatEvent:   eventName,
			MinFrequency:  0,
			MaxFrequency:  0,
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

func createTables(engine *xorm.Engine, name string) error {
	var relavence db.RelevanceDinamic
	var strength db.ControlDinamic
	var propused db.PropusedDinamic
	if err := db.CreateColumn(engine, relavence, name, "VARCHAR(255)"); err != nil {
		return err
	}
	if err := db.CreateColumn(engine, strength, name, "VARCHAR(255)"); err != nil {
		return err
	}
	if err := db.CreateColumn(engine, propused, name, "VARCHAR(255)"); err != nil {
		return err
	}
	RiskController := db.RiskController{
		Name: name,
	}
	if err := db.Create(engine, RiskController); err != nil {
		return err
	}
	modifyTables(name, "VARCHAR(255)")
	return nil
}

func modifyTables(name string, typeTable string) error {
	filename := "db/migration.dinamic.go"

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var modifiedLines []string
	inStruct := false

	for _, line := range lines {
		if strings.Contains(line, "type") && strings.Contains(line, "struct") {
			inStruct = true
		}

		if inStruct && strings.TrimSpace(line) == "}" {
			newLine := fmt.Sprintf("    %s %s `json:\"%s\" xorm:\"String notnull\"`", name, typeTable, name)
			modifiedLines = append(modifiedLines, newLine)
			inStruct = false
		}
		modifiedLines = append(modifiedLines, line)
	}

	modifiedContent := strings.Join(modifiedLines, "\n")

	err = ioutil.WriteFile(filename, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func PullAllEventService(c *gin.Context) {
	var events []db.ThreatEventCatalog
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
	var event db.ThreatEventCatalog
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
