package control

import (
	"fmt"
	"net/http"
	"qira/db"
	calculations "qira/internal/math"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullAllControlStrength(c *gin.Context) {
	var relevances []db.RelevanceDinamic
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

	implMap := make(map[int]db.Implements)
	for _, impl := range implementations {
		implMap[impl.ControlID] = impl
	}

	var wg sync.WaitGroup
	resultsChan := make(chan db.ControlDinamic)

	for controlID, impl := range implMap {
		wg.Add(1)
		go func(controlID int, impl db.Implements) {
			defer wg.Done()
			var totalRelevance int
			var totalAggregated float64
			for _, relevance := range relevances {
				if relevance.ControlID == controlID {
					val := reflect.ValueOf(relevance)
					for i := 0; i < val.NumField(); i++ {
						field := val.Type().Field(i)
						if strings.HasSuffix(field.Name, "Attack") {
							relevanceValue, err := strconv.Atoi(val.Field(i).String())
							if err != nil {
								continue
							}
							totalRelevance += relevanceValue
							totalAggregated += calculations.CalculateValue(relevanceValue, impl.Current)
						}
					}
				}
			}
			if totalRelevance > 0 {
				aggregated := totalAggregated / float64(totalRelevance)
				controlGap := 100.0 - aggregated
				result := db.ControlDinamic{
					ControlID:       controlID,
					AggregateTable:  "AuthenticationAttack",
					Aggregate:       fmt.Sprintf("%f", aggregated),
					ControlGapTable: "AuthenticationAttack",
					ControlGap:      fmt.Sprintf("%f", controlGap),
				}
				resultsChan <- result
			}
		}(controlID, impl)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var finalResults []db.ControlDinamic
	for result := range resultsChan {
		finalResults = append(finalResults, result)
	}

	if err := saveResultsControlDinamic(engine.(*xorm.Engine), finalResults); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", finalResults)
	c.Status(http.StatusOK)
}

func PullAllControlProposed(c *gin.Context) {
	var relevances []db.RelevanceDinamic
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

	implMap := make(map[int]db.Implements)
	for _, impl := range implementations {
		implMap[impl.ControlID] = impl
	}

	var wg sync.WaitGroup
	resultsChan := make(chan db.PropusedDinamic)

	for controlID, impl := range implMap {
		wg.Add(1)
		go func(controlID int, impl db.Implements) {
			defer wg.Done()
			var totalRelevance int
			var totalAggregated float64
			for _, relevance := range relevances {
				if relevance.ControlID == controlID {
					val := reflect.ValueOf(relevance)
					for i := 0; i < val.NumField(); i++ {
						field := val.Type().Field(i)
						if strings.HasSuffix(field.Name, "Attack") {
							relevanceValue, err := strconv.Atoi(val.Field(i).String())
							if err != nil {
								continue
							}
							totalRelevance += relevanceValue
							totalAggregated += calculations.CalculateValue(relevanceValue, impl.Proposed)
						}
					}
				}
			}
			if totalRelevance > 0 {
				aggregated := totalAggregated / float64(totalRelevance)
				controlGap := 100.0 - aggregated
				result := db.PropusedDinamic{
					ControlID:       controlID,
					AggregateTable:  "AuthenticationAttack",
					Aggregate:       fmt.Sprintf("%f", aggregated),
					ControlGapTable: "AuthenticationAttack",
					ControlGap:      fmt.Sprintf("%f", controlGap),
				}
				resultsChan <- result
			}
		}(controlID, impl)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var finalResults []db.PropusedDinamic
	for result := range resultsChan {
		finalResults = append(finalResults, result)
	}

	if err := saveResultsPropusedDinamic(engine.(*xorm.Engine), finalResults); err != nil {
		c.Set("Response", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", finalResults)
	c.Status(http.StatusOK)
}

func saveResultsControlDinamic(engine *xorm.Engine, results []db.ControlDinamic) error {
	for _, result := range results {
		if err := db.Create(engine, &result); err != nil {
			return err
		}
	}
	return nil
}
func saveResultsPropusedDinamic(engine *xorm.Engine, results []db.PropusedDinamic) error {
	for _, result := range results {
		if err := db.Create(engine, &result); err != nil {
			return err
		}
	}
	return nil
}

func CalculateAggregatedControlStrength(engine *xorm.Engine) ([]db.AggregatedStrength, error) {
	var controlStrengths []db.ControlDinamic
	var proposedStrengths []db.PropusedDinamic
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
		threatEventID, err := strconv.Atoi(cs.AggregateTable)
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
		threatEventID, err := strconv.Atoi(ps.AggregateTable)
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
