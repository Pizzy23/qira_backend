package control

import (
	"fmt"
	"net/http"
	"qira/db"
	"qira/internal/mock"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type ControlProposed struct {
	ControlID       int64
	AggregateTable  string
	Aggregate       float64
	ControlGapTable string
	ControlGap      float64
}

func PullAllControlStrength(c *gin.Context) {
	var controls []db.ControlLibrary
	var relevances []db.Relevance
	var implementations []db.Implements
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := db.GetAll(engine.(*xorm.Engine), &relevances); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := db.GetAll(engine.(*xorm.Engine), &implementations); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Mapping control IDs to their relevances and implementations
	relevanceMap := make(map[int64][]db.Relevance)
	for _, relevance := range relevances {
		relevanceMap[relevance.ControlID] = append(relevanceMap[relevance.ControlID], relevance)
	}

	implMap := make(map[int64]db.Implements)
	for _, impl := range implementations {
		implMap[impl.ControlID] = impl
	}

	type ControlStrength struct {
		ControlID    int64
		TypeOfAttack string
		Strength     float64
		Porcent      float64
	}

	controlStrengths := []ControlStrength{}
	totalRelevanceMap := make(map[string]float64)

	for _, control := range controls {
		impl, implExists := implMap[control.ID]
		if !implExists {
			continue
		}

		relevances, relevanceExists := relevanceMap[control.ID]
		if !relevanceExists {
			continue
		}

		for _, relevance := range relevances {
			typeOfAttack := relevance.TypeOfAttack

			relevanceAvgStr, err := mock.FindAverageByScore(int(relevance.Porcent))
			if err != nil {
				continue
			}
			relevanceValue, err := strconv.ParseFloat(strings.TrimSuffix(relevanceAvgStr, "%"), 64)
			if err != nil {
				continue
			}

			totalRelevanceMap[typeOfAttack] += relevanceValue

			currentValue, err := strconv.ParseFloat(strings.TrimSuffix(impl.PercentCurrent, "%"), 64)
			if err != nil {
				continue
			}

			porcent := CalculateValue(relevanceValue/100.0, currentValue/100.0)

			controlStrengths = append(controlStrengths, ControlStrength{
				ControlID:    control.ID,
				TypeOfAttack: typeOfAttack,
				Strength:     porcent,
				Porcent:      porcent,
			})
		}
	}

	controlStrengthMap := make(map[string]float64)
	porcentMap := make(map[int64]float64)
	for _, result := range controlStrengths {
		controlStrengthMap[result.TypeOfAttack] += result.Strength
		porcentMap[result.ControlID] = result.Porcent
	}

	var finalResults []db.Control
	for _, control := range controls {
		relevances, relevanceExists := relevanceMap[control.ID]
		if !relevanceExists {
			continue
		}

		for _, relevance := range relevances {
			finalResults = append(finalResults, db.Control{
				ControlID:    control.ID,
				TypeOfAttack: relevance.TypeOfAttack,
				Porcent:      fmt.Sprintf("%.2f%%", porcentMap[control.ID]*100), // Convertendo de decimal para percentual
			})
		}
	}

	for typeOfAttack, totalStrength := range controlStrengthMap {
		totalRelevance := totalRelevanceMap[typeOfAttack]
		aggregated := (totalStrength / totalRelevance) * 100.0
		controlGap := 100.0 - aggregated

		finalResults = append(finalResults, db.Control{
			ControlID:    -1,
			TypeOfAttack: typeOfAttack,
			Aggregate:    fmt.Sprintf("%.2f%%", aggregated),
		})

		finalResults = append(finalResults, db.Control{
			ControlID:    -2,
			TypeOfAttack: typeOfAttack,
			ControlGap:   fmt.Sprintf("%.2f%%", controlGap),
		})
	}

	if err := saveResultsControl(engine.(*xorm.Engine), finalResults); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", finalResults)
	c.Status(http.StatusOK)
}

func CalculateValue(relevance float64, current float64) float64 {
	return relevance * relevance * current
}

func PullAllControlProposed(c *gin.Context) {
	var relevances []db.Relevance
	var implementations []db.Implements

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &relevances); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &implementations); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	relevanceMap := make(map[int64][]db.Relevance)
	for _, relevance := range relevances {
		relevanceMap[relevance.ControlID] = append(relevanceMap[relevance.ControlID], relevance)
	}

	implMap := make(map[int64]db.Implements)
	for _, impl := range implementations {
		implMap[impl.ControlID] = impl
	}

	type ControlStrength struct {
		ControlID    int64
		TypeOfAttack string
		Strength     float64
		Porcent      float64
	}

	controlStrengths := []ControlStrength{}
	totalRelevanceMap := make(map[string]float64)

	for controlID, impl := range implMap {
		relevances, relevanceExists := relevanceMap[controlID]
		if !relevanceExists {
			continue
		}

		for _, relevance := range relevances {
			typeOfAttack := relevance.TypeOfAttack

			relevanceAvgStr, err := mock.FindAverageByScore(int(relevance.Porcent))
			if err != nil {
				continue
			}
			relevanceValue, err := strconv.ParseFloat(strings.TrimSuffix(relevanceAvgStr, "%"), 64)
			if err != nil {
				continue
			}

			totalRelevanceMap[typeOfAttack] += relevanceValue

			currentValue, err := strconv.ParseFloat(strings.TrimSuffix(impl.PercentCurrent, "%"), 64)
			if err != nil {
				continue
			}

			porcent := CalculateValue(relevanceValue/100.0, currentValue/100.0)

			controlStrengths = append(controlStrengths, ControlStrength{
				ControlID:    controlID,
				TypeOfAttack: typeOfAttack,
				Strength:     porcent,
				Porcent:      float64(impl.Proposed),
			})
		}
	}

	controlStrengthMap := make(map[string]float64)
	porcentMap := make(map[int64]float64)
	for _, result := range controlStrengths {
		controlStrengthMap[result.TypeOfAttack] += result.Strength
		porcentMap[result.ControlID] = result.Porcent
	}

	var finalResults []db.Propused
	for _, control := range relevances {
		finalResults = append(finalResults, db.Propused{
			ControlID:    control.ControlID,
			TypeOfAttack: control.TypeOfAttack,
			Porcent:      fmt.Sprintf("%.2f%%", porcentMap[control.ControlID]),
		})
	}

	for typeOfAttack, totalStrength := range controlStrengthMap {
		totalRelevance := totalRelevanceMap[typeOfAttack]
		aggregated := (totalStrength / totalRelevance) * 100.0
		controlGap := 100.0 - aggregated

		finalResults = append(finalResults, db.Propused{
			ControlID:    -1,
			TypeOfAttack: typeOfAttack,
			Aggregate:    fmt.Sprintf("%.2f%%", aggregated),
		})

		finalResults = append(finalResults, db.Propused{
			ControlID:    -2,
			TypeOfAttack: typeOfAttack,
			ControlGap:   fmt.Sprintf("%.2f%%", controlGap),
		})
	}

	if err := saveResultsPropused(engine.(*xorm.Engine), finalResults); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", finalResults)
	c.Status(http.StatusOK)
}

func saveResultsPropused(engine *xorm.Engine, results []db.Propused) error {
	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	for _, result := range results {
		existing := db.Propused{}
		has, err := engine.Where("control_i_d = ? AND type_of_attack = ?", result.ControlID, result.TypeOfAttack).Get(&existing)
		if err != nil {
			session.Rollback()
			return err
		}

		if has {
			result.ID = existing.ID
			if _, err := engine.ID(result.ID).Update(&result); err != nil {
				session.Rollback()
				return err
			}
		} else {
			if _, err := engine.Insert(&result); err != nil {
				session.Rollback()
				return err
			}
		}
	}

	return session.Commit()
}

func saveResultsControl(engine *xorm.Engine, results []db.Control) error {
	for _, result := range results {
		var existing db.Control
		found, err := engine.Where("control_i_d = ? AND type_of_attack = ?", result.ControlID, result.TypeOfAttack).Get(&existing)
		if err != nil {
			return err
		}

		if found {
			if existing.Porcent != result.Porcent || existing.Aggregate != result.Aggregate || existing.ControlGap != result.ControlGap {
				existing.Porcent = result.Porcent
				existing.Aggregate = result.Aggregate
				existing.ControlGap = result.ControlGap
				if _, err := engine.ID(existing.ID).Update(&existing); err != nil {
					return err
				}
			}
		} else {
			if err := db.Create(engine, &result); err != nil {
				return err
			}
		}
	}
	return nil
}

func CalculateAggregatedControlStrength(engine *xorm.Engine) ([]db.AggregatedStrength, error) {
	var controlStrengths []db.Control
	var proposedStrengths []db.Propused
	var threatEvents []db.ThreatEventCatalog

	if err := db.GetAll(engine, &controlStrengths); err != nil {
		return nil, err
	}
	if err := db.GetAll(engine, &proposedStrengths); err != nil {
		return nil, err
	}
	if err := db.GetAll(engine, &threatEvents); err != nil {
		return nil, err
	}

	aggregatedMap := make(map[string]db.AggregatedStrength)

	// Aggregate current control strengths
	for _, cs := range controlStrengths {
		if cs.ControlID == -1 {
			aggregatedMap[cs.TypeOfAttack] = db.AggregatedStrength{
				ThreatID:    0,
				ThreatEvent: cs.TypeOfAttack,
				Current:     cs.Aggregate,
				Proposed:    "", // Deixe proposto vazio por enquanto
			}
		}
	}

	// Aggregate proposed control strengths
	for _, ps := range proposedStrengths {
		if ps.ControlID == -1 {
			acs, exists := aggregatedMap[ps.TypeOfAttack]
			if exists {
				acs.Proposed = ps.Aggregate
				aggregatedMap[ps.TypeOfAttack] = acs
			}
		}
	}

	// Assign threat event names
	for _, te := range threatEvents {
		acs, exists := aggregatedMap[te.ThreatEvent]
		if exists {
			acs.ThreatID = te.ID
			aggregatedMap[te.ThreatEvent] = acs
		}
	}

	// Collect results
	var finalResults []db.AggregatedStrength
	for _, acs := range aggregatedMap {
		finalResults = append(finalResults, acs)
	}

	return finalResults, nil
}
