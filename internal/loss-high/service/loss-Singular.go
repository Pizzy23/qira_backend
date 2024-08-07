package losshigh

import (
	"errors"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateSingularLossService(c *gin.Context, LossHigh interfaces.InputLossHigh, id int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingLoss db.LossHigh
	found, err := engine.(*xorm.Engine).Where("threat_event_id = ? AND loss_type = ?", id, "Singular").Get(&existingLoss)
	if err != nil {
		return err
	}

	if found {
		existingLoss.ThreatEvent = LossHigh.ThreatEvent
		existingLoss.MinimumLoss = LossHigh.MinimumLoss
		existingLoss.MaximumLoss = LossHigh.MaximumLoss
		existingLoss.MostLikelyLoss = LossHigh.MostLikelyLoss

		if _, err := engine.(*xorm.Engine).ID(existingLoss.ID).Update(&existingLoss); err != nil {
			return err
		}
	} else {
		newLoss := db.LossHigh{
			ThreatEventID:  id,
			ThreatEvent:    LossHigh.ThreatEvent,
			LossType:       "Singular",
			MinimumLoss:    LossHigh.MinimumLoss,
			MaximumLoss:    LossHigh.MaximumLoss,
			MostLikelyLoss: LossHigh.MostLikelyLoss,
		}

		if err := db.Create(engine.(*xorm.Engine), &newLoss); err != nil {
			return err
		}
	}

	var existingTotal db.LossHighTotal
	totalFound, err := engine.(*xorm.Engine).Where("threat_event_id = ? AND name = 'Total' AND type_of_loss = 'Singular'", id).Get(&existingTotal)
	if err != nil {
		return err
	}

	if totalFound {
		existingTotal.MinimumLoss = LossHigh.MinimumLoss
		existingTotal.MaximumLoss = LossHigh.MaximumLoss
		existingTotal.MostLikelyLoss = LossHigh.MostLikelyLoss

		if _, err := engine.(*xorm.Engine).ID(existingTotal.ID).Update(&existingTotal); err != nil {
			return err
		}
	} else {
		newTotal := db.LossHighTotal{
			ThreatEventID:  id,
			ThreatEvent:    LossHigh.ThreatEvent,
			Name:           "Total",
			TypeOfLoss:     "Singular",
			MinimumLoss:    LossHigh.MinimumLoss,
			MaximumLoss:    LossHigh.MaximumLoss,
			MostLikelyLoss: LossHigh.MostLikelyLoss,
		}

		if err := db.Create(engine.(*xorm.Engine), &newTotal); err != nil {
			return err
		}
	}

	return nil
}

func GetSingularLosses(c *gin.Context) ([]AggregatedLossResponse, error) {
	var lossHighs []db.LossHigh
	var lossHighTotals []db.LossHighTotal
	engine, exists := c.Get("db")
	if !exists {
		return nil, errors.New("database connection not found")
	}

	dbEngine, ok := engine.(*xorm.Engine)
	if !ok {
		return nil, errors.New("failed to cast database connection to *xorm.Engine")
	}

	if err := db.GetAllWithCondition(dbEngine, &lossHighs, "loss_type = ?", "Singular"); err != nil {
		return nil, err
	}

	if err := db.GetAllWithCondition(dbEngine, &lossHighTotals, "name = 'Total' AND type_of_loss = 'Singular'"); err != nil {
		return nil, err
	}

	assets, err := getAssetsLossHigh(dbEngine, lossHighs[0])
	if err != nil {
		return nil, err
	}

	aggregatedData := make(map[int64]*AggregatedLossResponse)

	for _, loss := range lossHighs {
		if _, exists := aggregatedData[loss.ThreatEventID]; !exists {
			aggregatedData[loss.ThreatEventID] = &AggregatedLossResponse{
				ThreatEventID: loss.ThreatEventID,
				ThreatEvent:   loss.ThreatEvent,
				Assets:        assets,
				Losses:        []AggregatedLossDetail{},
			}
		}
		detail := AggregatedLossDetail{
			LossType:       loss.LossType,
			MinimumLoss:    loss.MinimumLoss,
			MaximumLoss:    loss.MaximumLoss,
			MostLikelyLoss: loss.MostLikelyLoss,
		}
		aggregatedData[loss.ThreatEventID].Losses = append(aggregatedData[loss.ThreatEventID].Losses, detail)
	}

	for _, agg := range aggregatedData {
		total := AggregatedLossDetail{
			LossType:       "Total",
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
		found, err := dbEngine.Where("threat_event_id = ? AND type_of_loss = 'Singular' AND name = 'Total'", agg.ThreatEventID).Get(&existingTotal)
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
				TypeOfLoss:     "Singular",
				MinimumLoss:    total.MinimumLoss,
				MaximumLoss:    total.MaximumLoss,
				MostLikelyLoss: total.MostLikelyLoss,
			}
			if _, err := dbEngine.Insert(&newTotal); err != nil {
				return nil, err
			}
		}
	}

	var result []AggregatedLossResponse
	for _, v := range aggregatedData {
		result = append(result, *v)
	}

	filteredResult, err := filterOutOfScopeAggregatedLosses(result, dbEngine)
	if err != nil {
		return nil, err
	}

	return filteredResult, nil
}
