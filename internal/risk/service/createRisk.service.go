package risk

import (
	"errors"
	"net/http"
	"qira/db"
	calculations "qira/internal/math"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateRiskService(c *gin.Context, typeLoss string) ([]db.RiskCalculation, error) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return nil, errors.New("database connection not found")
	}

	// Assegure que `engine` é do tipo `*xorm.Engine`
	dbEngine, ok := engine.(*xorm.Engine)
	if !ok {
		return nil, errors.New("failed to cast database connection to *xorm.Engine")
	}

	threatEvents, lossTotals, frequencies, err := getAll(dbEngine, typeLoss)
	if err != nil {
		return nil, err
	}

	var riskCalculations []db.RiskCalculation

	for _, event := range threatEvents {
		var loss db.LossHighTotal
		var freq db.Frequency

		for _, l := range lossTotals {
			if l.ThreatEventID == event.ID {
				loss = l
				break
			}
		}

		for _, f := range frequencies {
			if f.ThreatEventID == event.ID {
				freq = f
				break
			}
		}

		if loss.ThreatEventID == 0 || freq.ThreatEventID == 0 {
			continue
		}

		estimateLoss := calculations.CalcRisks(loss.MinimumLoss, loss.MostLikelyLoss, loss.MaximumLoss)
		riskCalculations = append(riskCalculations, db.RiskCalculation{
			ThreatEventID: event.ID,
			ThreatEvent:   event.ThreatEvent,
			RiskType:      typeLoss,
			Min:           loss.MinimumLoss,
			Max:           loss.MaximumLoss,
			Mode:          loss.MostLikelyLoss,
			Estimate:      estimateLoss,
		})

		estimateFrequency := calculations.CalcRisks(freq.MinFrequency, freq.MostLikelyFrequency, freq.MaxFrequency)
		riskCalculations = append(riskCalculations, db.RiskCalculation{
			ThreatEventID: event.ID,
			ThreatEvent:   event.ThreatEvent,
			RiskType:      "Frequency",
			Min:           freq.MinFrequency,
			Max:           freq.MaxFrequency,
			Mode:          freq.MostLikelyFrequency,
			Estimate:      estimateFrequency,
		})

		minRisk := freq.MinFrequency * loss.MinimumLoss
		maxRisk := freq.MaxFrequency * loss.MaximumLoss
		modeRisk := freq.MostLikelyFrequency * loss.MostLikelyLoss
		estimateRisk := calculations.CalcRisks(minRisk, modeRisk, maxRisk)
		riskCalculations = append(riskCalculations, db.RiskCalculation{
			ThreatEventID: event.ID,
			ThreatEvent:   event.ThreatEvent,
			RiskType:      "Risk",
			Min:           minRisk,
			Max:           maxRisk,
			Mode:          modeRisk,
			Estimate:      estimateRisk,
		})
	}

	for _, risk := range riskCalculations {
		if _, err := dbEngine.Insert(&risk); err != nil {
			return nil, err
		}
	}

	return riskCalculations, nil
}

func getAll(engine *xorm.Engine, lossType string) ([]db.ThreatEventCatalog, []db.LossHighTotal, []db.Frequency, error) {
	var threatEvents []db.ThreatEventCatalog
	var loss []db.LossHighTotal
	var frequency []db.Frequency

	if err := engine.Find(&threatEvents); err != nil {
		return nil, nil, nil, err
	}
	if err := engine.Where("type_of_loss = ?", lossType).Find(&loss); err != nil {
		return nil, nil, nil, err
	}
	if err := engine.Find(&frequency); err != nil {
		return nil, nil, nil, err
	}
	if len(threatEvents) <= 0 {
		return nil, nil, nil, errors.New("not have threat events")
	}
	if len(loss) <= 0 {
		return nil, nil, nil, errors.New("not have loss")
	}
	if len(frequency) <= 0 {
		return nil, nil, nil, errors.New("not have frequency")
	}
	return threatEvents, loss, frequency, nil
}
