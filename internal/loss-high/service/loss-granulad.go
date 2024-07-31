package losshigh

import (
	"errors"
	"qira/db"
	"qira/internal/interfaces"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateLossHighGranularService(c *gin.Context, LossHigh interfaces.InputLossHighGranulade, id int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingLoss db.LossHighGranular
	found, err := engine.(*xorm.Engine).Where("threat_event_i_d = ? AND loss_type = ? AND impact = ?", id, LossHigh.LossType, LossHigh.Impact).Get(&existingLoss)
	if err != nil {
		return err
	}

	if found {
		existingLoss.ThreatEvent = LossHigh.ThreatEvent
		existingLoss.Assets = strings.Join(LossHigh.Assets, ",")
		existingLoss.MinimumLoss = LossHigh.MinimumLoss
		existingLoss.MaximumLoss = LossHigh.MaximumLoss
		existingLoss.MostLikelyLoss = LossHigh.MostLikelyLoss

		if _, err := engine.(*xorm.Engine).ID(existingLoss.ID).Update(&existingLoss); err != nil {
			return err
		}
	} else {
		newLoss := db.LossHighGranular{
			ThreatEventID:  id,
			ThreatEvent:    LossHigh.ThreatEvent,
			Assets:         strings.Join(LossHigh.Assets, ","),
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
			aggregatedData[loss.ThreatEventID] = &AggregatedLossResponseGranulade{
				ThreatEventID: loss.ThreatEventID,
				ThreatEvent:   loss.ThreatEvent,
				Assets:        loss.Assets,
				Losses:        []AggregatedLossDetailGranulade{},
			}
		}
		detail := AggregatedLossDetailGranulade{
			LossType:       loss.LossType,
			Impact:         loss.Impact,
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
		found, err := dbEngine.Where("threat_event_i_d = ? AND type_of_loss = 'Granular' AND name = 'Total'", agg.ThreatEventID).Get(&existingTotal)
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

	return result, nil
}
