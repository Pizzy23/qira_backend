package revelance

import (
	"errors"
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullAllRevelance(c *gin.Context) {
	var revelances []db.Relevance
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &revelances); err != nil {
		c.Set("Error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", revelances)
	c.Status(http.StatusOK)
}

func PullRevelanceId(c *gin.Context, id int) {
	var revelance db.Relevance
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &revelance, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving revelance")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Revelance not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", revelance)
	c.Status(http.StatusOK)
}

func CreateRelevanceService(c *gin.Context, Relevance db.RelevanceDinamicInput) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &Relevance); err != nil {
		return err
	}
	return nil

}
