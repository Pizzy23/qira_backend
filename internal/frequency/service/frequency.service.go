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

func EditFrequencyService(c *gin.Context, freq interfaces.InputFrequency, ThreatEventID int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingFreq db.Frequency
	found, err := db.GetByID(engine.(*xorm.Engine), &existingFreq, ThreatEventID)
	if err != nil {
		return errors.New("failed to fetch existing frequency data")
	}
	if !found || existingFreq.ThreatEvent != freq.ThreatEvent {
		return errors.New("threat event mismatch or not found")
	}

	// If check passes, proceed with update
	frequencyTable := db.Frequency{
		ThreatEventID:         ThreatEventID,
		ThreatEvent:           freq.ThreatEvent,
		MinFrequency:          freq.MinFrequency,
		MaxFrequency:          freq.MaxFrequency,
		MostLikelyFrequency:   freq.MostCommonFrequency,
		SupportingInformation: freq.SupportInformation,
	}

	if err := db.UpdateByThreat(engine.(*xorm.Engine), frequencyTable, ThreatEventID); err != nil {
		return err
	}
	return nil
}

func PullAllEventService(c *gin.Context) {
	var frequencies []db.Frequency
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

	for _, event := range threatEvents {
		var frequency db.Frequency
		found, err := db.GetByEventIDAndRiskType(engine.(*xorm.Engine), &frequency, event.ID, "")
		if err != nil {
			c.Set("Response", "Error fetching frequency")
			c.Status(http.StatusInternalServerError)
			return
		}

		if !found {
			newFrequency := db.Frequency{
				ThreatEventID:         event.ID,
				ThreatEvent:           event.ThreatEvent,
				MinFrequency:          0.0,
				MaxFrequency:          0.0,
				MostLikelyFrequency:   0.0,
				SupportingInformation: "Default information",
			}

			if err := db.Create(engine.(*xorm.Engine), &newFrequency); err != nil {
				c.Set("Response", "Error inserting new frequency")
				c.Status(http.StatusInternalServerError)
				return
			}

			frequencies = append(frequencies, newFrequency)
		} else {
			frequencies = append(frequencies, frequency)
		}
	}

	c.Set("Response", frequencies)
	c.Status(http.StatusOK)
}

func PullEventIdService(c *gin.Context, id int) {
	var frequency db.Frequency
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &frequency, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving Frequency")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "Frequency not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", frequency)
	c.Status(http.StatusOK)
}
