package losshigh

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateLossHighService(c *gin.Context, LossHigh db.LossHigh) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &LossHigh); err != nil {
		return err
	}
	return nil

}

func PullAllLossHigh(c *gin.Context) {
	var lossHighs []interfaces.InputLossHigh
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &lossHighs); err != nil {
		c.Set("Error", "Error")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", lossHighs)
	c.Status(http.StatusOK)
}

func PullLossHighId(c *gin.Context, id int) {
	var lossHigh interfaces.InputLossHigh
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &lossHigh, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving LossHigh")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "LossHigh not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", lossHigh)
	c.Status(http.StatusOK)
}
