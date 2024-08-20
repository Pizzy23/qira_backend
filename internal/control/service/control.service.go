package control

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

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

func PullAllControl(c *gin.Context) {
	var controls []db.ControlLibrary
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", controls)
	c.Status(http.StatusOK)
}

func Stren(c *gin.Context) {
	var controls []db.Control
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", controls)
	c.Status(http.StatusOK)
}

func Prupu(c *gin.Context) {
	var controls []db.Propused
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", controls)
	c.Status(http.StatusOK)
}

func DeleteControl(c *gin.Context, id int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var asset db.ControlLibrary
	has, err := engine.(*xorm.Engine).ID(id).Get(&asset)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("Control not found")
	}

	if _, err := engine.(*xorm.Engine).Where("control_id = ?", id).Delete(&db.Relevance{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).Where("control_id = ?", id).Delete(&db.Control{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).Where("control_id = ?", id).Delete(&db.Propused{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).Where("control_id = ?", id).Delete(&db.Implements{}); err != nil {
		return err
	}

	if _, err := engine.(*xorm.Engine).ID(id).Delete(&db.ControlLibrary{}); err != nil {
		return err
	}

	return nil
}
