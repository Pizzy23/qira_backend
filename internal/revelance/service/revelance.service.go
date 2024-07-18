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
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &revelances); err != nil {
		c.Set("Response", err)
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
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &revelance, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving revelance")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "Revelance not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", revelance)
	c.Status(http.StatusOK)
}

func CreateRelevanceService(c *gin.Context, Relevance db.Relevance) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.UpdateByControlIdAndRisk(engine.(*xorm.Engine), &Relevance, Relevance.ControlID, Relevance.TypeOfAttack); err != nil {
		return err
	}
	return nil

}
