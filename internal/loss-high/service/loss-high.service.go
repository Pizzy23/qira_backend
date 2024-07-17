package losshigh

import (
	"errors"
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateLossHighService(c *gin.Context, LossHigh db.LossHigh) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &LossHigh); err != nil {
		return err
	}
	return nil

}

func GetAggregatedLosses(c *gin.Context) ([]AggregatedLossResponse, error) {
	var lossHighs []db.LossHigh
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
