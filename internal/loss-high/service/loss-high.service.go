package losshigh

import (
	"errors"
	"fmt"
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

func GetAggregatedLosses(c *gin.Context) ([]AggregatedLoss, error) {
	var lossHighs []db.LossHigh
	engine, exists := c.Get("db")
	if !exists {
		return nil, errors.New("database connection not found")
	}

	if err := db.GetAll(engine.(*xorm.Engine), &lossHighs); err != nil {
		return nil, err
	}

	aggregatedData := make(map[string]*AggregatedLoss)

	// Aggregate data
	for _, loss := range lossHighs {
		key := fmt.Sprintf("%s-%s", loss.ThreatEvent, loss.LossType)
		if _, exists := aggregatedData[key]; !exists {
			aggregatedData[key] = &AggregatedLoss{
				ThreatEvent:    loss.ThreatEvent,
				Assets:         loss.Assets,
				LossType:       loss.LossType,
				MinimumLoss:    0,
				MaximumLoss:    0,
				MostLikelyLoss: 0,
			}
		}
		aggregatedData[key].MinimumLoss += loss.MinimumLoss
		aggregatedData[key].MaximumLoss += loss.MaximumLoss
		aggregatedData[key].MostLikelyLoss += loss.MostLikelyLoss
	}

	// Convert map to slice
	var result []AggregatedLoss
	for _, v := range aggregatedData {
		result = append(result, *v)
	}

	// Calculate totals
	total := AggregatedLoss{
		ThreatEvent:    "Total",
		MinimumLoss:    0,
		MaximumLoss:    0,
		MostLikelyLoss: 0,
	}
	for _, agg := range result {
		total.MinimumLoss += agg.MinimumLoss
		total.MaximumLoss += agg.MaximumLoss
		total.MostLikelyLoss += agg.MostLikelyLoss
	}
	result = append(result, total)

	return result, nil
}

func PullLossHighId(c *gin.Context, id int) {
	var lossHigh db.LossHigh
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &lossHigh, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving LossHigh")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "LossHigh not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", lossHigh)
	c.Status(http.StatusOK)
}
