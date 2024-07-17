package catalogue

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qira/db"
	frequency "qira/internal/frequency/service"
	"qira/internal/interfaces"
	"qira/util"
	"strings"

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

	if err := createTables(engine.(*xorm.Engine), event.ThreatEvent); err != nil {
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
	if err := db.CreateColumn(engine, "relevance_dinamic", name, "VARCHAR(255)"); err != nil {
		return err
	}
	if err := db.CreateColumn(engine, "control_dinamic", name, "VARCHAR(255)"); err != nil {
		return err
	}
	if err := db.CreateColumn(engine, "propused_dinamic", name, "VARCHAR(255)"); err != nil {
		return err
	}
	RiskController := db.RiskController{
		Name: name,
	}
	if err := db.Create(engine, RiskController); err != nil {
		return err
	}
	modifyTables(name)
	return nil
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func modifyTables(columnName string) error {
	filename := "db/migration.dinamic.go"

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var modifiedLines []string
	inStruct := false

	capitalizedColumnName := capitalizeFirstLetter(columnName)

	for _, line := range lines {
		if strings.Contains(line, "type") && strings.Contains(line, "struct") {
			inStruct = true
		}

		if inStruct && strings.TrimSpace(line) == "}" {
			newLine := fmt.Sprintf("    %s %s `json:\"%s\" xorm:\"%s notnull\"`", capitalizedColumnName, "string", columnName, "VARCHAR(255)")
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
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &events); err != nil {
		c.Set("Response", err)
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
