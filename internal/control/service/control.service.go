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
	controlInput := db.ControlLibrary{
		ControlType:      control.ControlType,
		ControlReference: control.ControlReference,
		Information:      control.Information,
		InScope:          control.InScope,
	}

	if err := db.Create(engine.(*xorm.Engine), &controlInput); err != nil {
		return err
	}

	if err := createTables(engine.(*xorm.Engine), &controlInput); err != nil {
		return err
	}

	return nil
}

func UpdateControlService(c *gin.Context, controlID int64, control interfaces.InputControlLibrary) error {
	engineInterface, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	engine, ok := engineInterface.(*xorm.Engine)
	if !ok {
		return errors.New("invalid database connection")
	}

	var existingControl db.ControlLibrary
	has, err := engine.ID(controlID).Get(&existingControl)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("control not found")
	}

	// Atualizar os campos existentes com os novos valores
	existingControl.ControlType = control.ControlType
	existingControl.ControlReference = control.ControlReference
	existingControl.Information = control.Information
	existingControl.InScope = control.InScope

	if _, err := engine.ID(controlID).Update(&existingControl); err != nil {
		return err
	}

	return nil
}

func createTables(engine *xorm.Engine, control *db.ControlLibrary) error {
	var events []db.RiskController
	if err := db.GetAll(engine, &events); err != nil {
		return err
	}

	var relevances []db.Relevance
	var propuseds []db.Propused
	var controls []db.Control

	for _, e := range events {
		relevances = append(relevances, db.Relevance{
			ControlID:    control.ID,
			TypeOfAttack: e.Name,
			Porcent:      0,
		})
		propuseds = append(propuseds, db.Propused{
			ControlID:    control.ID,
			TypeOfAttack: e.Name,
			Porcent:      "0",
			Aggregate:    "0",
			ControlGap:   "0",
		})
		controls = append(controls, db.Control{
			ControlID:    control.ID,
			TypeOfAttack: e.Name,
			Porcent:      "0",
			Aggregate:    "0",
			ControlGap:   "0",
		})
	}

	if len(relevances) > 0 {
		if _, err := engine.Insert(&relevances); err != nil {
			return err
		}
	}
	if len(propuseds) > 0 {
		if _, err := engine.Insert(&propuseds); err != nil {
			return err
		}
	}
	if len(controls) > 0 {
		if _, err := engine.Insert(&controls); err != nil {
			return err
		}
	}

	return nil
}

func CreateImplementService(c *gin.Context, data interfaces.ImplementsInputNoID, id int64) error {
	averageC, err := mock.FindAverageByScore(data.Current)
	if err != nil {
		return errors.New("score not found for Percent Current")
	}
	averageP, err := mock.FindAverageByScore(data.Proposed)
	if err != nil {
		return errors.New("score not found for Percent Proposed")
	}
	implement := db.Implements{
		ControlID:       id,
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
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", controls)
	c.Status(http.StatusOK)
}
func PullAllImplements(c *gin.Context) {
	var controls []db.Implements
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", controls)
	c.Status(http.StatusOK)
}

func PullControlId(c *gin.Context, id int) {
	var control db.Implements
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &control, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving control")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "control not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", control)
	c.Status(http.StatusOK)
}
