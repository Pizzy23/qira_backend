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
	found, err := engine.(*xorm.Engine).Where("threat_event_i_d = ? AND loss_type = ?", id, LossHigh.LossType).Get(&existingLoss)
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
		found, err := engine.(*xorm.Engine).Where("threat_event_id = ? AND type_of_loss = 'LossHigh'", agg.ThreatEventID).Get(&existingTotal)
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
				TypeOfLoss:     "LossHigh",
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

func CreateSingularLossService(c *gin.Context, LossHigh interfaces.InputLossHigh, id int64) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingLoss db.LossHigh
	found, err := engine.(*xorm.Engine).Where("threat_event_i_d = ? AND loss_type = ?", id, "Singular").Get(&existingLoss)
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
			LossType:       "Singular",
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

func GetSingularLosses(c *gin.Context) ([]AggregatedLossResponse, error) {
	var lossHighs []db.LossHigh
	var totalGet []db.LossHighTotal
	engine, exists := c.Get("db")
	if !exists {
		return nil, errors.New("database connection not found")
	}

	if err := db.GetAllWithCondition(engine.(*xorm.Engine), &lossHighs, "loss_type = ?", "Singular"); err != nil {
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
		found, err := engine.(*xorm.Engine).Where("threat_event_id = ? AND type_of_loss = 'Singular'", agg.ThreatEventID).Get(&existingTotal)
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
				TypeOfLoss:     "Singular",
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
	var totalGet []db.LossHighTotal
	engine, exists := c.Get("db")
	if !exists {
		return nil, errors.New("database connection not found")
	}

	if err := db.GetAllWithCondition(engine.(*xorm.Engine), &lossHighs, "loss_type IN ('Direct', 'Indirect')"); err != nil {
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
			LossType:       "Total",
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
		found, err := engine.(*xorm.Engine).Where("threat_event_id = ? AND type_of_loss = 'Granular'", agg.ThreatEventID).Get(&existingTotal)
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
				TypeOfLoss:     "Granular",
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

	var result []AggregatedLossResponseGranulade
	for _, v := range aggregatedData {
		result = append(result, *v)
	}

	return result, nil
}

func CreateLossSpecific(c *gin.Context, typeOfLoss string) {
	var lossesInput []db.LossHigh
	var lossesInputGranular []db.LossHighGranular
	var catalogue []db.ThreatEventCatalog

	engine, exists := c.Get("db")

	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &catalogue); err != nil {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(catalogue) <= 0 {
		c.Set("Response", "Not contein events")
		c.Status(http.StatusInternalServerError)
	}

	for _, event := range catalogue {
		switch typeOfLoss {
		case "Singular":
			lossesInput = append(lossesInput, lossesNotGranu(event, "Singular"))
		case "LossHigh":
			lossesInput = append(lossesInput, lossesNotGranu(event, "Indirect"))
			lossesInput = append(lossesInput, lossesNotGranu(event, "Direct"))
		case "Granular":
			lossesInputGranular = append(lossesInputGranular, lossesWithGranu(event, "Indirect", "Short Term"))
			lossesInputGranular = append(lossesInputGranular, lossesWithGranu(event, "Direct", "Short Term"))
			lossesInputGranular = append(lossesInputGranular, lossesWithGranu(event, "Indirect", "Long Term"))
			lossesInputGranular = append(lossesInputGranular, lossesWithGranu(event, "Direct", "Long Term"))
		}
	}

	if len(lossesInputGranular) != 0 {
		for _, loss := range lossesInputGranular {
			if err := db.Create(engine.(*xorm.Engine), &loss); err != nil {
				continue
			}
		}
		c.Set("Response", lossesInputGranular)
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(lossesInput) != 0 {
		for _, loss := range lossesInput {
			if err := db.Create(engine.(*xorm.Engine), &loss); err != nil {
				continue
			}
		}
		c.Set("Response", lossesInput)
		c.Status(http.StatusInternalServerError)
		return
	}
}

func lossesNotGranu(input db.ThreatEventCatalog, lossType string) db.LossHigh {
	return db.LossHigh{
		ThreatEventID:  input.ID,
		ThreatEvent:    input.ThreatEvent,
		Assets:         "",
		LossType:       lossType,
		MinimumLoss:    0,
		MaximumLoss:    0,
		MostLikelyLoss: 0,
	}
}

func lossesWithGranu(input db.ThreatEventCatalog, lossType string, impact string) db.LossHighGranular {
	return db.LossHighGranular{
		ThreatEventID:  input.ID,
		ThreatEvent:    input.ThreatEvent,
		Assets:         "",
		Impact:         impact,
		LossType:       lossType,
		MinimumLoss:    0,
		MaximumLoss:    0,
		MostLikelyLoss: 0,
	}
}
