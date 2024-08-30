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
		if !event.InScope {
			continue
		}

		var loss db.LossHighTotal
		var freq db.Frequency

		for _, l := range lossTotals {
			if l.ThreatEventID == event.ID && l.TypeOfLoss == typeLoss {
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

		// Verificar e criar/atualizar registros para Loss com base no typeLoss enviado
		estimateLoss := calculations.CalcRisks(loss.MinimumLoss, loss.MostLikelyLoss, loss.MaximumLoss)
		risk := db.RiskCalculation{
			ThreatEventID: event.ID,
			ThreatEvent:   event.ThreatEvent,
			Categorie:     typeLoss,
			RiskType:      "Loss",
			Min:           loss.MinimumLoss,
			Max:           loss.MaximumLoss,
			Mode:          loss.MostLikelyLoss,
			Estimate:      estimateLoss,
		}
		if err := checkAndUpdateRiskCalculation(dbEngine, risk); err != nil {
			return nil, err
		}

		// Verificar e criar/atualizar registros para Frequency com a categoria enviada
		estimateFrequency := calculations.CalcRisks(freq.MinFrequency, freq.MostLikelyFrequency, freq.MaxFrequency)
		risk = db.RiskCalculation{
			ThreatEventID: event.ID,
			ThreatEvent:   event.ThreatEvent,
			Categorie:     typeLoss,
			RiskType:      "Frequency",
			Min:           freq.MinFrequency,
			Max:           freq.MaxFrequency,
			Mode:          freq.MostLikelyFrequency,
			Estimate:      estimateFrequency,
		}
		if err := checkAndUpdateRiskCalculation(dbEngine, risk); err != nil {
			return nil, err
		}

		// Verificar e criar/atualizar registros para Risk com a categoria enviada
		minRisk := freq.MinFrequency * loss.MinimumLoss
		maxRisk := freq.MaxFrequency * loss.MaximumLoss
		modeRisk := freq.MostLikelyFrequency * loss.MostLikelyLoss
		estimateRisk := calculations.CalcRisks(minRisk, modeRisk, maxRisk)
		risk = db.RiskCalculation{
			ThreatEventID: event.ID,
			ThreatEvent:   event.ThreatEvent,
			Categorie:     typeLoss,
			RiskType:      "Risk",
			Min:           minRisk,
			Max:           maxRisk,
			Mode:          modeRisk,
			Estimate:      estimateRisk,
		}
		if err := checkAndUpdateRiskCalculation(dbEngine, risk); err != nil {
			return nil, err
		}
	}

	// Filtrar os resultados para retornar apenas a categoria solicitada
	var filteredRiskCalculations []db.RiskCalculation
	for _, risk := range riskCalculations {
		if risk.Categorie == typeLoss {
			filteredRiskCalculations = append(filteredRiskCalculations, risk)
		}
	}

	return filteredRiskCalculations, nil
}

func checkAndUpdateRiskCalculation(engine *xorm.Engine, risk db.RiskCalculation) error {
	var existingRisk db.RiskCalculation
	exists, err := engine.Where("threat_event_id = ? AND risk_type = ? AND categorie = ?", risk.ThreatEventID, risk.RiskType, risk.Categorie).Get(&existingRisk)
	if err != nil {
		return err
	}

	if exists {
		if existingRisk.Min != risk.Min || existingRisk.Max != risk.Max || existingRisk.Mode != risk.Mode || existingRisk.Estimate != risk.Estimate {
			existingRisk.Min = risk.Min
			existingRisk.Max = risk.Max
			existingRisk.Mode = risk.Mode
			existingRisk.Estimate = risk.Estimate
			if _, err := engine.ID(existingRisk.ID).Update(&existingRisk); err != nil {
				return err
			}
		}
	} else {
		if _, err := engine.Insert(&risk); err != nil {
			return err
		}
	}

	return nil
}

func getAll(engine *xorm.Engine, lossType string) ([]db.ThreatEventCatalog, []db.LossHighTotal, []db.Frequency, error) {
	var threatEvents []db.ThreatEventCatalog
	var loss []db.LossHighTotal
	var frequency []db.Frequency

	if err := engine.Where("in_scope = ?", true).Find(&threatEvents); err != nil {
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
