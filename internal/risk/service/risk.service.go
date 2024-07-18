package risk

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullAllRisk(c *gin.Context) {
	risk, err := CreateRiskService(c)
	if err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", risk)
	c.Status(http.StatusOK)
}

func PullRiskId(c *gin.Context, id int64) {
	var Risk db.RiskCalculation
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetRiskById(engine.(*xorm.Engine), id)
	if err != nil {
		c.Set("Error", "Error retrieving Risk")
		c.Status(http.StatusInternalServerError)
		return
	}
	if found == nil {
		c.Set("Error", "Risk not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", Risk)
	c.Status(http.StatusOK)
}
