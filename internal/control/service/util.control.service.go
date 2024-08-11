package control

import (
	"qira/db"
	"strings"

	"xorm.io/xorm"
)

type ControlProposed struct {
	ControlID       int64
	AggregateTable  string
	Aggregate       float64
	ControlGapTable string
	ControlGap      float64
}

type ControlStrength struct {
	ControlID    int64
	TypeOfAttack string
	Strength     float64
	Porcent      float64
}

func CalculateValue(relevance float64, current float64) float64 {
	return (relevance * relevance) * current
}

func saveResultsPropused(engine *xorm.Engine, results []db.Propused) error {
	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	for _, result := range results {
		existing := db.Propused{}
		has, err := engine.Where("control_id = ? AND type_of_attack = ?", result.ControlID, result.TypeOfAttack).Get(&existing)
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
func saveResultsStrength(engine *xorm.Engine, results []db.Control) error {
	session := engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	for _, result := range results {
		existing := db.Control{}
		has, err := engine.Where("control_id = ? AND type_of_attack = ?", result.ControlID, result.TypeOfAttack).Get(&existing)
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
		found, err := engine.Where("control_id = ? AND type_of_attack = ?", result.ControlID, result.TypeOfAttack).Get(&existing)
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
	if err := db.GetAllWithCondition(engine, &threatEvents, "in_scope = ?", true); err != nil {
		return nil, err
	}
	aggregatedMap := make(map[string]db.AggregatedStrength)

	for _, cs := range controlStrengths {
		if cs.ControlID == -1 {
			aggregatedMap[cs.TypeOfAttack] = db.AggregatedStrength{
				ThreatID:    0,
				ThreatEvent: cs.TypeOfAttack,
				Current:     cs.Aggregate,
				Proposed:    "",
			}
		}
	}

	for _, ps := range proposedStrengths {
		if ps.ControlID == -1 {
			acs, exists := aggregatedMap[ps.TypeOfAttack]
			if exists {
				acs.Proposed = ps.Aggregate
				aggregatedMap[ps.TypeOfAttack] = acs
			}
		}
	}

	for _, te := range threatEvents {
		acs, exists := aggregatedMap[te.ThreatEvent]
		if exists {
			acs.ThreatID = te.ID
			aggregatedMap[te.ThreatEvent] = acs
		}
	}

	var finalResults []db.AggregatedStrength
	for _, acs := range aggregatedMap {
		finalResults = append(finalResults, acs)
	}

	return finalResults, nil
}

func validateStrengthDataExist(engine *xorm.Engine, finalResults []db.Control) ([]db.Control, []db.Control, error) {
	var existingProposedData []db.Control
	err := engine.Find(&existingProposedData)
	if err != nil {
		return nil, nil, err
	}

	var dataToUpdate []db.Control
	var dataToAdd []db.Control
	found := false

	for _, result := range finalResults {
		found = false
		for _, existing := range existingProposedData {
			if result.ControlID == existing.ControlID && result.TypeOfAttack == existing.TypeOfAttack {
				found = true
				if result.Porcent != existing.Porcent || result.Aggregate != existing.Aggregate || result.ControlGap != existing.ControlGap || result.Information != existing.Information {
					dataToUpdate = append(dataToUpdate, result)
				}
				break
			}
		}
		if !found {
			dataToAdd = append(dataToAdd, result)
		}
	}

	return dataToUpdate, dataToAdd, nil
}
func validateProposedDataExist(engine *xorm.Engine, finalResults []db.Propused) ([]db.Propused, []db.Propused, error) {
	var existingProposedData []db.Propused
	err := engine.Find(&existingProposedData)
	if err != nil {
		return nil, nil, err
	}

	var dataToUpdate []db.Propused
	var dataToAdd []db.Propused
	found := false

	for _, result := range finalResults {
		found = false
		for _, existing := range existingProposedData {
			if result.ControlID == existing.ControlID && result.TypeOfAttack == existing.TypeOfAttack {
				found = true
				if result.Porcent != existing.Porcent || result.Aggregate != existing.Aggregate || result.ControlGap != existing.ControlGap || result.Information != existing.Information {
					dataToUpdate = append(dataToUpdate, result)
				}
				break
			}
		}
		if !found {
			dataToAdd = append(dataToAdd, result)
		}
	}

	return dataToUpdate, dataToAdd, nil
}

func uniqueTypeOfAttacks(typeOfAttacks []string) []string {
	attackMap := make(map[string]bool)
	uniqueAttacks := []string{}

	for _, attack := range typeOfAttacks {
		if _, exists := attackMap[attack]; !exists {
			attackMap[attack] = true
			uniqueAttacks = append(uniqueAttacks, attack)
		}
	}

	return uniqueAttacks
}

func validateEventsInScope(engine *xorm.Engine, typeOfAttacks []string) (map[string]bool, error) {
	uniqueAttacks := uniqueTypeOfAttacks(typeOfAttacks)

	var events []db.ThreatEventCatalog
	err := engine.In("threat_event", uniqueAttacks).Find(&events)
	if err != nil {
		return nil, err
	}

	inScopeMap := make(map[string]bool)
	for _, event := range events {
		inScopeMap[strings.ToLower(event.ThreatEvent)] = event.InScope
	}

	return inScopeMap, nil
}
