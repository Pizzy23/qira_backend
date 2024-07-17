package risk

import (
	"errors"
	"fmt"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	losshigh "qira/internal/loss-high/service"
	calculations "qira/internal/math"
	"strings"
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
	risk, threat, err := getThreatAndRisks(engine.(*xorm.Engine))
	if err != nil {
		return nil, err
	}

	if len(risk) == len(threat) {
		return risk, nil
	}

	_, freq, err := getAll(engine.(*xorm.Engine))
	if err != nil {
		return nil, errors.New("database connection not found")
	}

	aggregatedLossControles, err := getAggregatedLossControles(c)
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
			if err := handleRiskCalculation(engine.(*xorm.Engine), &freqCalc); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Process aggregated losses
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, aggregatedLossControl := range aggregatedLossControles {
			if aggregatedLossControl.LossType == "Total" {
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
				if err := handleRiskCalculation(engine.(*xorm.Engine), &lossCalc); err != nil {
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
			if err := handleRiskCalculation(engine.(*xorm.Engine), &riskCalc); err != nil {
				errChan <- err
				return
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

func getAll(engine *xorm.Engine) ([]db.LossHigh, []db.Frequency, error) {
	var loss []db.LossHigh
	var frequency []db.Frequency

	if err := engine.Find(&loss); err != nil {
		return nil, nil, err
	}
	if err := engine.Find(&frequency); err != nil {
		return nil, nil, err
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

func getAggregatedLossControles(c *gin.Context) ([]losshigh.AggregatedLossControl, error) {
	var lossHighs []db.LossHigh
	engine, exists := c.Get("db")
	if !exists {
		return nil, errors.New("database connection not found")
	}

	if err := db.GetAll(engine.(*xorm.Engine), &lossHighs); err != nil {
		return nil, err
	}

	aggregatedData := make(map[string]*losshigh.AggregatedLossControl)
	threatEventTotals := make(map[string]*losshigh.AggregatedLossControl)
	threatEventAssets := make(map[string][]string)

	for _, loss := range lossHighs {
		key := fmt.Sprintf("%s-%s", loss.ThreatEvent, loss.LossType)
		if _, exists := aggregatedData[key]; !exists {
			aggregatedData[key] = &losshigh.AggregatedLossControl{
				ThreatEvent:    loss.ThreatEvent,
				ThreatEventId:  loss.ThreatEventID,
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

		if _, exists := threatEventTotals[loss.ThreatEvent]; !exists {
			threatEventTotals[loss.ThreatEvent] = &losshigh.AggregatedLossControl{
				ThreatEvent:    loss.ThreatEvent,
				ThreatEventId:  loss.ThreatEventID,
				Assets:         "",
				LossType:       "Total",
				MinimumLoss:    0,
				MaximumLoss:    0,
				MostLikelyLoss: 0,
			}
		}
		threatEventTotals[loss.ThreatEvent].MinimumLoss += loss.MinimumLoss
		threatEventTotals[loss.ThreatEvent].MaximumLoss += loss.MaximumLoss
		threatEventTotals[loss.ThreatEvent].MostLikelyLoss += loss.MostLikelyLoss

		if _, exists := threatEventAssets[loss.ThreatEvent]; !exists {
			threatEventAssets[loss.ThreatEvent] = []string{}
		}
		threatEventAssets[loss.ThreatEvent] = append(threatEventAssets[loss.ThreatEvent], loss.Assets)
	}

	var result []losshigh.AggregatedLossControl
	for _, v := range aggregatedData {
		result = append(result, *v)
	}

	for _, total := range threatEventTotals {
		total.Assets = strings.Join(threatEventAssets[total.ThreatEvent], ", ")
		result = append(result, *total)
	}

	return result, nil
}
