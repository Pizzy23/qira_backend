package control

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	"qira/internal/mock"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateControlService(c *gin.Context, control interfaces.InputControlLibrary) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &control); err != nil {
		return err
	}
	return nil

}
func CreateImplementService(c *gin.Context, data interfaces.ImplementsInput) error {
	averageC, err := mock.FindAverageByScore(data.Current)
	if err != nil {
		return errors.New("score not found for Percent Current")
	}
	averageP, err := mock.FindAverageByScore(data.Proposed)
	if err != nil {
		return errors.New("score not found for Percent Proposed")
	}
	implement := db.Implements{
		ControlID:       data.ControlID,
		Current:         data.Current,
		Proposed:        data.Proposed,
		Cost:            data.Cost,
		PercentCurrent:  averageC,
		PercentProposed: averageP,
	}
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &implement); err != nil {
		return err
	}
	return nil

}

func PullAllControl(c *gin.Context) {
	var controls []db.ControlLibrary
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", controls)
	c.Status(http.StatusOK)
}

func PullControlId(c *gin.Context, id int) {
	var control db.ControlLibrary
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &control, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving control")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "control not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", control)
	c.Status(http.StatusOK)
}
