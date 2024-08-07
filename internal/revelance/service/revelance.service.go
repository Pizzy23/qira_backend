package revelance

import (
	"errors"
	"fmt"
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullAllRelevance(c *gin.Context) {
	var relevances []db.Relevance
	var threatEvents []db.ThreatEventCatalog
	var controls []db.ControlLibrary
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := engine.(*xorm.Engine).Where("in_scope = ?", true).Find(&threatEvents); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := engine.(*xorm.Engine).Where("in_scope = ?", true).Find(&controls); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &relevances); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	existingRelevanceMap := make(map[string]bool)
	for _, relevance := range relevances {
		key := fmt.Sprintf("%d-%s", relevance.ControlID, relevance.TypeOfAttack)
		existingRelevanceMap[key] = true
	}

	for _, event := range threatEvents {
		for _, control := range controls {
			key := fmt.Sprintf("%d-%s", control.ID, event.ThreatEvent)
			if !existingRelevanceMap[key] {
				newRelevance := db.Relevance{
					ControlID:    control.ID,
					TypeOfAttack: event.ThreatEvent,
					Porcent:      0,
				}
				if _, err := engine.(*xorm.Engine).Insert(&newRelevance); err != nil {
					c.Set("Response", err.Error())
					c.Status(http.StatusInternalServerError)
					return
				}
			}
		}
	}

	if err := db.GetAll(engine.(*xorm.Engine), &relevances); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", relevances)
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

func CreateRelevanceService(c *gin.Context, relevanceInput db.RelevanceDinamicInput) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var relevanceDb db.Relevance
	found, err := engine.(*xorm.Engine).Where("control_id = ? AND type_of_attack = ?", relevanceInput.ControlID, relevanceInput.TypeOfAttack).Get(&relevanceDb)
	if err != nil {
		return err
	}

	if found {
		relevanceDb.Porcent = relevanceInput.Porcent
		relevanceDb.TypeOfAttack = relevanceInput.TypeOfAttack

		affected, err := engine.(*xorm.Engine).ID(relevanceDb.ID).Update(&relevanceDb)
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("no columns found to be updated")
		}
	} else {
		relevanceDb = db.Relevance{
			ControlID:    relevanceInput.ControlID,
			TypeOfAttack: relevanceInput.TypeOfAttack,
			Porcent:      relevanceInput.Porcent,
		}
		_, err := engine.(*xorm.Engine).Insert(&relevanceDb)
		if err != nil {
			return err
		}
	}

	return nil
}
