package implementation

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	"qira/internal/mock"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullAllImplements(c *gin.Context) {
	var implementations []db.Implements
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &implementations); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", implementations)
	c.Status(http.StatusOK)
}

func PullImplementsId(c *gin.Context, id int) {
	var implementation db.Implements
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &implementation, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving Implementation")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "Implementation not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", implementation)
	c.Status(http.StatusOK)
}

func CreateImplementsService(c *gin.Context, Implements interfaces.ImplementsInput) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}
	scoreP, err := mock.FindAverageByScore(Implements.Proposed)
	if err != nil {
		return err
	}
	scoreC, err := mock.FindAverageByScore(Implements.Current)
	if err != nil {
		return err
	}
	dataDB := db.Implements{
		ControlID:       Implements.ControlID,
		Current:         Implements.Current,
		Proposed:        Implements.Proposed,
		PercentCurrent:  scoreC,
		PercentProposed: scoreP,
		Cost:            Implements.Cost,
	}
	if err := db.Create(engine.(*xorm.Engine), dataDB); err != nil {
		return err
	}
	return nil

}
