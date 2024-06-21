package control

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateRelevanceService(c *gin.Context, data interfaces.InputControlls) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}
	table, err := setTypes(data)
	if err != nil {
		if err := db.Create(engine.(*xorm.Engine), table); err != nil {
			return err
		}
		return nil
	}
	return nil
}
func CreateStrengthService(c *gin.Context, data interfaces.InputControlls) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}
	table, err := setTypes(data)
	if err != nil {
		if err := db.Create(engine.(*xorm.Engine), table); err != nil {
			return err
		}
		return nil
	}
	return nil
}
func CreatePropusedService(c *gin.Context, data interfaces.InputControlls) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}
	table, err := setTypes(data)
	if err != nil {
		if err := db.Create(engine.(*xorm.Engine), table); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func CreateLibraryService(c *gin.Context, data db.ControlLibrary) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), data); err != nil {
		return err
	}
	return nil

}

func CreateImplementationService(c *gin.Context, data db.ControlImplementation) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), data); err != nil {
		return err
	}
	return nil

}

func GetControl(c *gin.Context, controlType string) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	table, err := setTypesGet(controlType)
	if err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), table); err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", table)
	c.Status(http.StatusOK)
}

func GetById(c *gin.Context, controlType string, id int64) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	table, err := setTypesGet(controlType)
	if err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), table, id)
	if err != nil {
		c.Set("Error", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", table)
	c.Status(http.StatusOK)
}
