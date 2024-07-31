package losshigh

import (
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

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
		c.Set("Response", "No events found")
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, event := range catalogue {
		switch typeOfLoss {
		case "Singular":
			loss := lossesNotGranu(event, "Singular")
			exists, err := checkLossExists(engine.(*xorm.Engine), loss.ThreatEventID, loss.LossType)
			if err == nil && !exists {
				lossesInput = append(lossesInput, loss)
			}
		case "LossHigh":
			lossIndirect := lossesNotGranu(event, "Indirect")
			existsIndirect, errIndirect := checkLossExists(engine.(*xorm.Engine), lossIndirect.ThreatEventID, lossIndirect.LossType)
			if errIndirect == nil && !existsIndirect {
				lossesInput = append(lossesInput, lossIndirect)
			}

			lossDirect := lossesNotGranu(event, "Direct")
			existsDirect, errDirect := checkLossExists(engine.(*xorm.Engine), lossDirect.ThreatEventID, lossDirect.LossType)
			if errDirect == nil && !existsDirect {
				lossesInput = append(lossesInput, lossDirect)
			}
		case "Granular":
			losses := []db.LossHighGranular{
				lossesWithGranu(event, "Indirect", "Short Term"),
				lossesWithGranu(event, "Direct", "Short Term"),
				lossesWithGranu(event, "Indirect", "Long Term"),
				lossesWithGranu(event, "Direct", "Long Term"),
			}
			for _, loss := range losses {
				exists, err := checkLossGranularExists(engine.(*xorm.Engine), loss.ThreatEventID, loss.LossType, loss.Impact)
				if err == nil && !exists {
					lossesInputGranular = append(lossesInputGranular, loss)
				}
			}
		}
	}

	if len(lossesInputGranular) != 0 {
		for _, loss := range lossesInputGranular {
			if err := db.Create(engine.(*xorm.Engine), &loss); err != nil {
				continue
			}
		}
		c.Set("Response", lossesInputGranular)
		c.Status(http.StatusOK)
		return
	}

	if len(lossesInput) != 0 {
		for _, loss := range lossesInput {
			if err := db.Create(engine.(*xorm.Engine), &loss); err != nil {
				continue
			}
		}
		c.Set("Response", lossesInput)
		c.Status(http.StatusOK)
		return
	}

	c.Set("Response", "No new losses to add")
	c.Status(http.StatusOK)
}

func checkLossExists(engine *xorm.Engine, threatEventID int64, typeOfLoss string) (bool, error) {
	var loss db.LossHigh
	exists, err := engine.Where("threat_event_i_d = ? AND loss_type = ?", threatEventID, typeOfLoss).Get(&loss)
	return exists, err
}

func checkLossGranularExists(engine *xorm.Engine, threatEventID int64, typeOfLoss, impact string) (bool, error) {
	var loss db.LossHighGranular
	exists, err := engine.Where("threat_event_i_d = ? AND loss_type = ? AND impact = ?", threatEventID, typeOfLoss, impact).Get(&loss)
	return exists, err
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
