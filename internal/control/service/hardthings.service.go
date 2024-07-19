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
				Porcent:      fmt.Sprintf("%.2f%%", porcentMap[control.ID]),
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

	// Mapping control IDs to their relevances and implementations
	relevanceMap := make(map[int64][]db.Relevance)
	for _, relevance := range relevances {
		relevanceMap[relevance.ControlID] = append(relevanceMap[relevance.ControlID], relevance)
	}

	implMap := make(map[int64]db.Implements)
	for _, impl := range implementations {
		implMap[impl.ControlID] = impl
	}

	controlStrengthMap := make(map[string]float64)
	totalRelevanceMap := make(map[string]float64)

	for controlID, impl := range implMap {
		relevances, relevanceExists := relevanceMap[controlID]
		if !relevanceExists {
			continue
		}

		for _, relevance := range relevances {
			relevanceAvgStr, err := mock.FindAverageByScore(int(relevance.Porcent))
			if err != nil {
				continue
			}
			relevanceValue, err := strconv.ParseFloat(strings.TrimSuffix(relevanceAvgStr, "%"), 64)
			if err != nil {
				continue
			}

			totalRelevanceMap[relevance.TypeOfAttack] += relevanceValue
			controlStrengthMap[relevance.TypeOfAttack] += CalculateValue(relevanceValue/100.0, float64(impl.Proposed)/100.0)
		}
	}

	var finalResults []db.Propused

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
		if _, err := engine.Insert(&result); err != nil {
			session.Rollback()
			return err
		}
	}

	return session.Commit()
}

func saveResultsControl(engine *xorm.Engine, results []db.Control) error {
	for _, result := range results {
		if err := db.Create(engine, &result); err != nil {
			return err
		}
	}
	return nil
}

func CalculateAggregatedControlStrength(engine *xorm.Engine) ([]db.AggregatedStrength, error) {
	var controlStrengths []db.Control
	var proposedStrengths []db.Propused
	var threatEvents []db.ThreatEventAssets

	if err := db.GetAll(engine, &controlStrengths); err != nil {
		return nil, err
	}
	if err := db.GetAll(engine, &proposedStrengths); err != nil {
		return nil, err
	}
	if err := db.GetAll(engine, &threatEvents); err != nil {
		return nil, err
	}

	aggregatedMap := make(map[int]db.AggregatedStrength)

	// Aggregate current control strengths
	for _, cs := range controlStrengths {
		threatEventID, err := strconv.Atoi(cs.TypeOfAttack)
		if err != nil {
			continue
		}
		if _, exists := aggregatedMap[threatEventID]; !exists {
			aggregatedMap[threatEventID] = db.AggregatedStrength{
				ThreatID: threatEventID,
				Current:  "0%",
				Proposed: "0%",
			}
		}
		acs := aggregatedMap[threatEventID]
		acs.Current = addPercentages(acs.Current, cs.Aggregate)
		aggregatedMap[threatEventID] = acs
	}

	// Aggregate proposed control strengths
	for _, ps := range proposedStrengths {
		threatEventID, err := strconv.Atoi(ps.TypeOfAttack)
		if err != nil {
			continue
		}
		if _, exists := aggregatedMap[threatEventID]; !exists {
			aggregatedMap[threatEventID] = db.AggregatedStrength{
				ThreatID: threatEventID,
				Current:  "0%",
				Proposed: "0%",
			}
		}
		acs := aggregatedMap[threatEventID]
		acs.Proposed = addPercentages(acs.Proposed, ps.Aggregate)
		aggregatedMap[threatEventID] = acs
	}

	// Assign threat event names
	for _, te := range threatEvents {
		if acs, exists := aggregatedMap[te.ID]; exists {
			acs.ThreatEvent = te.ThreatEvent
			aggregatedMap[te.ID] = acs
		}
	}

	// Collect results
	var finalResults []db.AggregatedStrength
	for _, acs := range aggregatedMap {
		finalResults = append(finalResults, acs)
	}

	return finalResults, nil
}

func parsePercentage(percentageStr string) float64 {
	percentageStr = strings.TrimSuffix(percentageStr, "%")
	value, _ := strconv.ParseFloat(percentageStr, 64)
	return value
}

func addPercentages(p1, p2 string) string {
	v1 := parsePercentage(p1)
	v2 := parsePercentage(p2)
	total := v1 + v2
	return strconv.FormatFloat(total, 'f', 2, 64) + "%"
}
