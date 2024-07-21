package risk

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	losshigh "qira/internal/loss-high/service"
	calculations "qira/internal/math"
	"sync"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateRiskService(c *gin.Context) ([]db.RiskCalculation, error) {
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return nil, errors.New("database connection not found")
	}

	xormEngine, ok := engine.(*xorm.Engine)
	if !ok {
		return nil, errors.New("type assertion to *xorm.Engine failed")
	}

	risk, threat, err := getThreatAndRisks(xormEngine)
	if err != nil {
		return nil, err
	}
	if len(threat) <= 0 {
		return nil, errors.New("not have Event catalogue")
	}
	if len(risk) == len(threat)*3 {
		return risk, nil
	}

	_, freq, err := getAll(xormEngine)
	if err != nil {
		return nil, err
	}

	aggregatedLossControles, err := aggregateLosses(xormEngine)
	if err != nil {
		return nil, errors.New("error getting aggregated losses")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// Process frequencies
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, frequency := range freq {
			if hasLoss(frequency.ThreatEventID, aggregatedLossControles) {
				calc := calculations.CalcRisks(frequency.MinFrequency, frequency.MostLikelyFrequency, frequency.MaxFrequency)
				freqCalc := db.RiskCalculation{
					ThreatEventID: frequency.ThreatEventID,
					ThreatEvent:   frequency.ThreatEvent,
					RiskType:      "Frequency",
					Min:           frequency.MinFrequency,
					Max:           frequency.MaxFrequency,
					Mode:          frequency.MostLikelyFrequency,
					Estimate:      calc,
				}
				if err := handleRiskCalculation(xormEngine, &freqCalc); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	// Process aggregated losses
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, aggregatedLossControl := range aggregatedLossControles {
			if hasFrequency(aggregatedLossControl.ThreatEventId, freq) {
				calc := calculations.CalcRisks(aggregatedLossControl.MinimumLoss, aggregatedLossControl.MostLikelyLoss, aggregatedLossControl.MaximumLoss)
				lossCalc := db.RiskCalculation{
					ThreatEventID: aggregatedLossControl.ThreatEventId,
					ThreatEvent:   aggregatedLossControl.ThreatEvent,
					RiskType:      "Loss",
					Min:           aggregatedLossControl.MinimumLoss,
					Max:           aggregatedLossControl.MaximumLoss,
					Mode:          aggregatedLossControl.MostLikelyLoss,
					Estimate:      calc,
				}
				if err := handleRiskCalculation(xormEngine, &lossCalc); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	// Combine frequency and loss data
	combinedRisks := combineFrequencyAndLoss(freq, aggregatedLossControles)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, risk := range combinedRisks {
			if hasLoss(risk.ThreatEventID, aggregatedLossControles) && hasFrequency(risk.ThreatEventID, freq) {
				minRisk := risk.MinFrequency * risk.MinimumLoss
				maxRisk := risk.MaxFrequency * risk.MaximumLoss
				modeRisk := risk.MostLikelyFrequency * risk.MostLikelyLoss
				estimateRisk := calculations.CalcRisks(minRisk, modeRisk, maxRisk)

				riskCalc := db.RiskCalculation{
					ThreatEventID: risk.ThreatEventID,
					ThreatEvent:   risk.ThreatEvent,
					RiskType:      "Risk",
					Min:           minRisk,
					Max:           maxRisk,
					Mode:          modeRisk,
					Estimate:      estimateRisk,
				}
				if err := handleRiskCalculation(xormEngine, &riskCalc); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func handleRiskCalculation(engine *xorm.Engine, riskCalc *db.RiskCalculation) error {
	existingRiskCalc := db.RiskCalculation{}
	has, err := db.GetByEventIDAndRiskType(engine, &existingRiskCalc, riskCalc.ThreatEventID, riskCalc.RiskType)
	if err != nil {
		return err
	}

	if has {
		if existingRiskCalc.Min != riskCalc.Min || existingRiskCalc.Max != riskCalc.Max || existingRiskCalc.Mode != riskCalc.Mode || existingRiskCalc.Estimate != riskCalc.Estimate {
			_, err := engine.ID(existingRiskCalc.ID).Update(riskCalc)
			if err != nil {
				return err
			}
		}
	} else {
		if err := db.Create(engine, riskCalc); err != nil {
			return err
		}
	}
	return nil
}

func getAll(engine *xorm.Engine) ([]db.LossHighTotal, []db.Frequency, error) {
	var loss []db.LossHighTotal
	var frequency []db.Frequency

	if err := engine.Find(&loss); err != nil {
		return nil, nil, err
	}
	if err := engine.Find(&frequency); err != nil {
		return nil, nil, err
	}
	if len(loss) <= 0 {
		return nil, nil, errors.New("not have loss")
	}
	if len(frequency) <= 0 {
		return nil, nil, errors.New("not have frequency")
	}
	return loss, frequency, nil
}

func combineFrequencyAndLoss(freqs []db.Frequency, losses []losshigh.AggregatedLossControl) []interfaces.CombinedRisk {
	combined := []interfaces.CombinedRisk{}
	freqMap := make(map[int64]db.Frequency)
	lossMap := make(map[int64]losshigh.AggregatedLossControl)

	for _, freq := range freqs {
		freqMap[freq.ThreatEventID] = freq
	}

	for _, loss := range losses {
		lossMap[loss.ThreatEventId] = loss
	}

	for id, freq := range freqMap {
		if loss, exists := lossMap[id]; exists {
			combined = append(combined, interfaces.CombinedRisk{
				ThreatEventID:       id,
				ThreatEvent:         freq.ThreatEvent,
				MinFrequency:        freq.MinFrequency,
				MaxFrequency:        freq.MaxFrequency,
				MostLikelyFrequency: freq.MostLikelyFrequency,
				MinimumLoss:         loss.MinimumLoss,
				MaximumLoss:         loss.MaximumLoss,
				MostLikelyLoss:      loss.MostLikelyLoss,
			})
		}
	}

	return combined
}

func getThreatAndRisks(engine *xorm.Engine) ([]db.RiskCalculation, []db.ThreatEventCatalog, error) {
	var risk []db.RiskCalculation
	var threat []db.ThreatEventCatalog

	if err := engine.Find(&risk); err != nil {
		return nil, nil, err
	}
	if err := engine.Find(&threat); err != nil {
		return nil, nil, err
	}
	return risk, threat, nil

}

func aggregateLosses(engine *xorm.Engine) ([]losshigh.AggregatedLossControl, error) {
	var lossHighs []db.LossHigh
	if err := engine.Find(&lossHighs); err != nil {
		return nil, err
	}

	aggregatedData := make(map[int64]*losshigh.AggregatedLossControl)
	for _, loss := range lossHighs {
		if _, exists := aggregatedData[loss.ThreatEventID]; !exists {
			aggregatedData[loss.ThreatEventID] = &losshigh.AggregatedLossControl{
				ThreatEventId:  loss.ThreatEventID,
				ThreatEvent:    loss.ThreatEvent,
				MinimumLoss:    0,
				MaximumLoss:    0,
				MostLikelyLoss: 0,
			}
		}
		aggregatedData[loss.ThreatEventID].MinimumLoss += loss.MinimumLoss
		aggregatedData[loss.ThreatEventID].MaximumLoss += loss.MaximumLoss
		aggregatedData[loss.ThreatEventID].MostLikelyLoss += loss.MostLikelyLoss
	}

	var result []losshigh.AggregatedLossControl
	for _, v := range aggregatedData {
		result = append(result, *v)
	}

	return result, nil
}

func hasLoss(threatEventID int64, losses []losshigh.AggregatedLossControl) bool {
	for _, loss := range losses {
		if loss.ThreatEventId == threatEventID {
			return true
		}
	}
	return false
}

func hasFrequency(threatEventID int64, freqs []db.Frequency) bool {
	for _, freq := range freqs {
		if freq.ThreatEventID == threatEventID {
			return true
		}
	}
	return false
}
func handleLossHighTotal(engine *xorm.Engine, loss *db.LossHighTotal) error {
	existingLoss := db.LossHighTotal{}
	has, err := engine.Where("threat_event_id = ?", loss.ThreatEventID).Get(&existingLoss)
	if err != nil {
		return err
	}

	if has {
		if existingLoss.MinimumLoss != loss.MinimumLoss || existingLoss.MaximumLoss != loss.MaximumLoss || existingLoss.MostLikelyLoss != loss.MostLikelyLoss {
			_, err := engine.ID(existingLoss.ID).Update(loss)
			if err != nil {
				return err
			}
		}
	} else {
		if _, err := engine.Insert(loss); err != nil {
			return err
		}
	}
	return nil
}

func handleFrequency(engine *xorm.Engine, frequency *db.Frequency) error {
	existingFrequency := db.Frequency{}
	has, err := engine.Where("threat_event_id = ?", frequency.ThreatEventID).Get(&existingFrequency)
	if err != nil {
		return err
	}

	if has {
		if existingFrequency.MinFrequency != frequency.MinFrequency || existingFrequency.MaxFrequency != frequency.MaxFrequency || existingFrequency.MostLikelyFrequency != frequency.MostLikelyFrequency {
			_, err := engine.ID(existingFrequency.ID).Update(frequency)
			if err != nil {
				return err
			}
		}
	} else {
		if _, err := engine.Insert(frequency); err != nil {
			return err
		}
	}
	return nil
}
