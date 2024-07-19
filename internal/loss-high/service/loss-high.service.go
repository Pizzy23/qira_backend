package losshigh

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateLossHighService(c *gin.Context, LossHigh interfaces.InputLossHigh, id int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingLoss db.LossHigh
	found, err := engine.(*xorm.Engine).Where("threat_event_id = ? AND loss_type = ?", id, LossHigh.LossType).Get(&existingLoss)
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
		newLoss := db.LossHigh{
			ThreatEventID:  id,
			ThreatEvent:    LossHigh.ThreatEvent,
			Assets:         strings.Join(LossHigh.Assets, ","),
			LossType:       LossHigh.LossType,
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

func GetAggregatedLosses(c *gin.Context) ([]AggregatedLossResponse, error) {
	var lossHighs []db.LossHigh
	var totalGet []db.LossHighTotal
	engine, exists := c.Get("db")
	if !exists {
		return nil, errors.New("database connection not found")
	}

	if err := db.GetAll(engine.(*xorm.Engine), &lossHighs); err != nil {
		return nil, err
	}

	aggregatedData := make(map[int64]*AggregatedLossResponse)

	for _, loss := range lossHighs {
		if _, exists := aggregatedData[loss.ThreatEventID]; !exists {
			aggregatedData[loss.ThreatEventID] = &AggregatedLossResponse{
				ThreatEventID: loss.ThreatEventID,
				ThreatEvent:   loss.ThreatEvent,
				Assets:        loss.Assets,
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
		found, err := engine.(*xorm.Engine).Where("threat_event_i_d = ?", agg.ThreatEventID).Get(&existingTotal)
		if err != nil {
			return nil, err
		}

		if found {
			if existingTotal.MinimumLoss != total.MinimumLoss || existingTotal.MaximumLoss != total.MaximumLoss || existingTotal.MostLikelyLoss != total.MostLikelyLoss {
				existingTotal.MinimumLoss = total.MinimumLoss
				existingTotal.MaximumLoss = total.MaximumLoss
				existingTotal.MostLikelyLoss = total.MostLikelyLoss
				if _, err := engine.(*xorm.Engine).ID(existingTotal.ID).Update(&existingTotal); err != nil {
					return nil, err
				}
			}
		} else {
			totalGet = append(totalGet, db.LossHighTotal{
				ThreatEventID:  agg.ThreatEventID,
				ThreatEvent:    agg.ThreatEvent,
				MinimumLoss:    total.MinimumLoss,
				MaximumLoss:    total.MaximumLoss,
				MostLikelyLoss: total.MostLikelyLoss,
			})
		}
	}

	for _, total := range totalGet {
		if _, err := engine.(*xorm.Engine).Insert(&total); err != nil {
			return nil, err
		}
	}

	var result []AggregatedLossResponse
	for _, v := range aggregatedData {
		result = append(result, *v)
	}

	return result, nil
}

func PullLossHighId(c *gin.Context, id int) {
	var lossHigh db.LossHigh
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &lossHigh, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving LossHigh")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "LossHigh not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", lossHigh)
	c.Status(http.StatusOK)
}
