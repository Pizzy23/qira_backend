package losshigh

import (
	"errors"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateLossHighGranularService(c *gin.Context, LossHigh interfaces.InputLossHighGranulade, id int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingLoss db.LossHighGranular
	found, err := engine.(*xorm.Engine).Where("threat_event_id = ? AND loss_type = ? AND impact = ? AND loss_edit_number = ?",
		id, LossHigh.LossType, LossHigh.Impact, LossHigh.LossEditNumber).Get(&existingLoss)
	if err != nil {
		return err
	}

	if found {
		existingLoss.ThreatEvent = LossHigh.ThreatEvent
		existingLoss.Impact = LossHigh.Impact
		existingLoss.MinimumLoss = LossHigh.MinimumLoss
		existingLoss.MaximumLoss = LossHigh.MaximumLoss
		existingLoss.MostLikelyLoss = LossHigh.MostLikelyLoss

		if _, err := engine.(*xorm.Engine).ID(existingLoss.ID).Update(&existingLoss); err != nil {
			return err
		}
	} else {
		// Verificar o próximo `LossEditNumber` disponível
		var maxEditNumber int64
		_, err := engine.(*xorm.Engine).Table("loss_high_granular").
			Where("threat_event_id = ? AND loss_type = ? AND impact = ?", id, LossHigh.LossType, LossHigh.Impact).
			Select("MAX(loss_edit_number)").
			Get(&maxEditNumber)
		if err != nil {
			return err
		}

		// Atribuir o próximo número de edição disponível (de 1 a 4)
		nextEditNumber := maxEditNumber + 1
		if nextEditNumber > 4 {
			return errors.New("maximum loss_edit_number exceeded for this combination of threat_event_id, loss_type, and impact")
		}

		newLoss := db.LossHighGranular{
			ThreatEventID:  id,
			ThreatEvent:    LossHigh.ThreatEvent,
			LossEditNumber: nextEditNumber,
			LossType:       LossHigh.LossType,
			Impact:         LossHigh.Impact,
			MinimumLoss:    LossHigh.MinimumLoss,
			MaximumLoss:    LossHigh.MaximumLoss,
			MostLikelyLoss: LossHigh.MostLikelyLoss,
		}

		if err := db.Create(engine.(*xorm.Engine), &newLoss); err != nil {
			return err
		}
	}

	return nil
}

func GetGranularLosses(c *gin.Context) ([]AggregatedLossResponseGranulade, error) {
	var lossHighs []db.LossHighGranular
	var lossHighTotals []db.LossHighTotal
	engine, exists := c.Get("db")
	if !exists {
		return nil, errors.New("database connection not found")
	}

	dbEngine, ok := engine.(*xorm.Engine)
	if !ok {
		return nil, errors.New("failed to cast database connection to *xorm.Engine")
	}

	if err := db.GetAllWithCondition(dbEngine, &lossHighs, "loss_type IN ('Direct', 'Indirect')"); err != nil {
		return nil, err
	}

	if err := db.GetAllWithCondition(dbEngine, &lossHighTotals, "name = 'Total' AND type_of_loss = 'Granular'"); err != nil {
		return nil, err
	}

	aggregatedData := make(map[int64]*AggregatedLossResponseGranulade)

	for _, loss := range lossHighs {
		if _, exists := aggregatedData[loss.ThreatEventID]; !exists {
			assets, err := getAssetsLossGran(dbEngine, loss)
			if err != nil {
				return nil, err
			}
			aggregatedData[loss.ThreatEventID] = &AggregatedLossResponseGranulade{
				ThreatEventID: loss.ThreatEventID,
				ThreatEvent:   loss.ThreatEvent,
				Assets:        assets,
				Losses:        []AggregatedLossDetailGranulade{},
			}
		}

		lossEditNumber := int64(len(aggregatedData[loss.ThreatEventID].Losses) + 1)

		detail := AggregatedLossDetailGranulade{
			LossType:       loss.LossType,
			Impact:         loss.Impact,
			LossEditNumber: lossEditNumber,
			MinimumLoss:    loss.MinimumLoss,
			MaximumLoss:    loss.MaximumLoss,
			MostLikelyLoss: loss.MostLikelyLoss,
		}
		aggregatedData[loss.ThreatEventID].Losses = append(aggregatedData[loss.ThreatEventID].Losses, detail)
	}

	for _, agg := range aggregatedData {
		total := AggregatedLossDetailGranulade{
			LossType:       "Granular",
			Impact:         "Total",
			LossEditNumber: int64(len(agg.Losses) + 1),
			MinimumLoss:    0,
			MaximumLoss:    0,
			MostLikelyLoss: 0,
		}
		for _, detail := range agg.Losses {
			total.MinimumLoss += detail.MinimumLoss
			total.MaximumLoss += detail.MaximumLoss
			total.MostLikelyLoss += detail.MostLikelyLoss
		}
		agg.Losses = append(agg.Losses, total)

		existingTotal := db.LossHighTotal{}
		found, err := dbEngine.Where("threat_event_id = ? AND type_of_loss = 'Granular' AND name = 'Total'", agg.ThreatEventID).Get(&existingTotal)
		if err != nil {
			return nil, err
		}

		if found {
			if existingTotal.MinimumLoss != total.MinimumLoss || existingTotal.MaximumLoss != total.MaximumLoss || existingTotal.MostLikelyLoss != total.MostLikelyLoss {
				existingTotal.MinimumLoss = total.MinimumLoss
				existingTotal.MaximumLoss = total.MaximumLoss
				existingTotal.MostLikelyLoss = total.MostLikelyLoss
				if _, err := dbEngine.ID(existingTotal.ID).Update(&existingTotal); err != nil {
					return nil, err
				}
			}
		} else {
			newTotal := db.LossHighTotal{
				ThreatEventID:  agg.ThreatEventID,
				ThreatEvent:    agg.ThreatEvent,
				Name:           "Total",
				TypeOfLoss:     "Granular",
				MinimumLoss:    total.MinimumLoss,
				MaximumLoss:    total.MaximumLoss,
				MostLikelyLoss: total.MostLikelyLoss,
			}
			if _, err := dbEngine.Insert(&newTotal); err != nil {
				return nil, err
			}
		}
	}

	var result []AggregatedLossResponseGranulade
	for _, v := range aggregatedData {
		result = append(result, *v)
	}

	filteredResult, err := filterOutOfScopeAggregatedLossesGranulade(result, dbEngine)
	if err != nil {
		return nil, err
	}

	return filteredResult, nil
}
