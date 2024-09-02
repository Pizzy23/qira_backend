package risk

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullAllRisk(c *gin.Context, typeLoss string) {
	var calcRisk []db.RiskCalculation
	var events []db.ThreatEventCatalog
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	_, err := CreateRiskService(c, typeLoss)
	if err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.InScope(engine.(*xorm.Engine).NewSession(), &events); err != nil {
		c.Set("Response", "Error fetching threat events")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &calcRisk); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	var filteredRisk []db.RiskCalculation
	for _, event := range events {
		for _, risk := range calcRisk {
			if risk.Categorie == typeLoss && risk.ThreatEvent == event.ThreatEvent {
				filteredRisk = append(filteredRisk, risk)
			}
		}
	}
	c.Set("Response", filteredRisk)
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
