package control

import (
	"fmt"
	"math"
	"net/http"
	"qira/db"
	"qira/internal/mock"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

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
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &relevances); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &implementations); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	typesOfAttack := make([]string, 0)
	for _, relevance := range relevances {
		typesOfAttack = append(typesOfAttack, strings.ToLower(relevance.TypeOfAttack))
	}

	eventsInScope, err := validateEventsInScope(engine.(*xorm.Engine), typesOfAttack)
	if err != nil {
		c.Set("Response", "Failed to validate threat events: "+err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	controlInScopeMap := make(map[int64]bool)
	for _, control := range controls {
		controlInScopeMap[control.ID] = control.InScope
	}

	controlInScopeMap[-1] = true
	controlInScopeMap[-2] = true

	relevanceMap := make(map[int64][]db.Relevance)
	for _, relevance := range relevances {
		if eventsInScope[strings.ToLower(relevance.TypeOfAttack)] {
			relevanceMap[relevance.ControlID] = append(relevanceMap[relevance.ControlID], relevance)
		}
	}

	implMap := make(map[int64]db.Implements)
	for _, impl := range implementations {
		implMap[impl.ControlID] = impl
	}

	type ControlStrength struct {
		ControlID    int64
		TypeOfAttack string
		Information  string
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
			var relevanceValue float64
			typeOfAttack := relevance.TypeOfAttack
			lowerCaseTypeOfAttack := strings.ToLower(typeOfAttack)

			relevanceAvgStr, err := mock.FindAverageByScore(int(relevance.Porcent))
			if err != nil {
				continue
			}
			if relevanceAvgStr != "N/A" {
				relevanceValue, err = strconv.ParseFloat(strings.TrimSuffix(relevanceAvgStr, "%"), 64)
				if err != nil {
					continue
				}
			} else {
				relevanceValue = 0
			}
			totalRelevanceMap[lowerCaseTypeOfAttack] += relevanceValue

			currentValue, err := strconv.ParseFloat(strings.TrimSuffix(impl.PercentCurrent, "%"), 64)
			if err != nil {
				continue
			}

			porcent := CalculateValue(relevanceValue/100.0, currentValue/100.0)

			controlStrengths = append(controlStrengths, ControlStrength{
				ControlID:    control.ID,
				Information:  control.Information,
				TypeOfAttack: strings.Title(lowerCaseTypeOfAttack),
				Strength:     porcent,
				Porcent:      porcent,
			})
		}
	}

	controlStrengthMap := make(map[string]float64)
	porcentMap := make(map[int64]float64)

	allZero := make(map[string]bool)

	for _, result := range controlStrengths {
		controlStrengthMap[result.TypeOfAttack] += result.Strength
		porcentMap[result.ControlID] = result.Porcent
		if result.Porcent == 0 {
			allZero[result.TypeOfAttack] = true
		} else {
			allZero[result.TypeOfAttack] = false
		}
	}

	var finalResults []db.Control
	for _, control := range controls {
		relevances, relevanceExists := relevanceMap[control.ID]
		if !relevanceExists {
			continue
		}

		for _, relevance := range relevances {
			finalResults = append(finalResults, db.Control{
				ID:           relevance.ID,
				ControlID:    control.ID,
				TypeOfAttack: strings.Title(strings.ToLower(relevance.TypeOfAttack)),
				Porcent:      fmt.Sprintf("%.0f%%", math.Floor(porcentMap[control.ID]*100)),
			})
		}
	}

	for typeOfAttack, totalStrength := range controlStrengthMap {
		totalRelevance := totalRelevanceMap[strings.ToLower(typeOfAttack)]
		if allZero[typeOfAttack] {
			finalResults = append(finalResults, db.Control{
				ControlID:    -1,
				TypeOfAttack: strings.Title(typeOfAttack),
				Aggregate:    "0%",
			})
			finalResults = append(finalResults, db.Control{
				ControlID:    -2,
				TypeOfAttack: strings.Title(typeOfAttack),
				ControlGap:   "100%",
			})
		} else {
			notPorcent := (totalStrength * 100.0)
			aggregated := (notPorcent / totalRelevance) * 100.0
			controlGap := 100.0 - int(aggregated)

			finalResults = append(finalResults, db.Control{
				ControlID:    -1,
				TypeOfAttack: strings.Title(typeOfAttack),
				Aggregate:    fmt.Sprintf("%d%%", int(aggregated)),
			})

			finalResults = append(finalResults, db.Control{
				ControlID:    -2,
				TypeOfAttack: strings.Title(typeOfAttack),
				ControlGap:   fmt.Sprintf("%d%%", controlGap),
			})
		}
	}

	dataToUpdate, dataToAdd, err := validateStrengthDataExist(engine.(*xorm.Engine), finalResults)
	if err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(dataToAdd) > 0 {
		if err := saveResultsStrength(engine.(*xorm.Engine), dataToAdd); err != nil {
			c.Set("Response", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	if len(dataToUpdate) > 0 {
		if err := saveResultsStrength(engine.(*xorm.Engine), dataToUpdate); err != nil {
			c.Set("Response", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	if len(dataToAdd) == 0 && len(dataToUpdate) == 0 {
		var proposed []db.Control
		if err := db.GetAll(engine.(*xorm.Engine), &proposed); err != nil {
			c.Set("Response", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		var filteredProposed []db.Control
		for _, prop := range proposed {
			if inScope, exists := controlInScopeMap[prop.ControlID]; exists && inScope {
				if eventInScope, eventExists := eventsInScope[strings.ToLower(prop.TypeOfAttack)]; eventExists && eventInScope {
					filteredProposed = append(filteredProposed, prop)
				}
			}
		}

		c.Set("Response", filteredProposed)
		c.Status(http.StatusOK)
		return
	}
	c.Set("Response", finalResults)
	c.Status(http.StatusOK)
}
