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
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	session := engine.(*xorm.Engine).NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	var (
		threatEvents []db.ThreatEventCatalog
		controls     []db.ControlLibrary
		relevances   []db.Relevance
	)

	// Buscar eventos de ameaça
	if err := db.InScope(session, &threatEvents); err != nil {
		c.Set("Response", "Error fetching threat events")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Buscar controles
	if err := db.InScope(session, &controls); err != nil {
		c.Set("Response", "Error fetching controls")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Buscar relevâncias
	if err := session.Find(&relevances); err != nil {
		c.Set("Response", "Error fetching relevances")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Filtro de relevâncias baseado nos eventos de ameaça
	var filteredRelevances []db.Relevance
	for _, event := range threatEvents {
		for _, relevance := range relevances {
			if relevance.TypeOfAttack == event.ThreatEvent {
				filteredRelevances = append(filteredRelevances, relevance)
			}
		}
	}

	// Validar relevâncias filtradas com os controles
	for _, event := range threatEvents {
		for _, control := range controls {
			key := fmt.Sprintf("%d-%s", control.ID, event.ThreatEvent)
			found := false
			for _, relevance := range filteredRelevances {
				if fmt.Sprintf("%d-%s", relevance.ControlID, relevance.TypeOfAttack) == key {
					found = true
					break
				}
			}
			// Se não encontrar uma relevância existente, cria uma nova
			if !found {
				newRelevance := db.Relevance{
					ControlID:    control.ID,
					Information:  control.Information,
					TypeOfAttack: event.ThreatEvent,
					Porcent:      0,
				}
				if _, err := session.Insert(&newRelevance); err != nil {
					session.Rollback()
					c.Set("Response", err.Error())
					c.Status(http.StatusInternalServerError)
					return
				}
				filteredRelevances = append(filteredRelevances, newRelevance)
			}
		}
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", filteredRelevances)
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
