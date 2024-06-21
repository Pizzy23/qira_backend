package risk

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateRiskService(c *gin.Context, Risk interfaces.InputRiskCalculator) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &Risk); err != nil {
		return err
	}
	return nil

}

func PullAllRisk(c *gin.Context) {
	var Risks []interfaces.InputRiskCalculator
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &Risks); err != nil {
		c.Set("Error", "Error")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", Risks)
	c.Status(http.StatusOK)
}

func PullRiskId(c *gin.Context, id int) {
	var Risk interfaces.InputRiskCalculator
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &Risk, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving Risk")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Risk not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", Risk)
	c.Status(http.StatusOK)
}
