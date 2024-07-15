package risk

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	calculations "qira/internal/math"
	"sync"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateRiskService(c *gin.Context, Risk interfaces.RiskCalc) ([]db.RiskCalculation, error) {
	engine, ok := c.MustGet("db_engine").(*xorm.Engine)
	if !ok {
		return nil, errors.New("database connection not found")
	}
	risk, threat, err := getThreatAndRisks(engine)
	if err != nil {
		return nil, err
	}

	if len(risk) == len(threat) {
		return risk, nil
	}

	loss, freq, err := getAll(engine)
	if err != nil {
		return nil, errors.New("database connection not found")
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
				RiskType:      "Loss",
				Min:           frequency.MinFrequency,
				Max:           frequency.MaxFrequency,
				Mode:          frequency.MostLikelyFrequency,
				Estimate:      calc,
			}
			if err := db.Create(engine, &freqCalc); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Process losses
	refinedLoss := processLossHighData(loss)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, refined := range refinedLoss {
			calc := calculations.CalcRisks(refined.MinimumLoss, refined.MostLikelyLoss, refined.MaximumLoss)
			lossCalc := db.RiskCalculation{
				ThreatEventID: refined.ThreatEventID,
				ThreatEvent:   refined.ThreatEvent,
				RiskType:      "Loss",
				Min:           refined.MinimumLoss,
				Max:           refined.MaximumLoss,
				Mode:          refined.MostLikelyLoss,
				Estimate:      calc,
			}
			if err := db.Create(engine, &lossCalc); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Combine frequency and loss data
	combinedRisks := combineFrequencyAndLoss(freq, loss)
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

			if _, err := engine.Insert(&riskCalc); err != nil {
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

func processLossHighData(lossHigh []db.LossHigh) []db.LossHigh {
	var filteredLoss []db.LossHigh
	for _, loss := range lossHigh {
		if loss.LossType == "Total" {
			filteredLoss = append(filteredLoss, loss)
		}
	}
	return filteredLoss
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

func combineFrequencyAndLoss(freqs []db.Frequency, losses []db.LossHigh) []interfaces.CombinedRisk {
	combined := []interfaces.CombinedRisk{}
	freqMap := make(map[int64]db.Frequency)
	lossMap := make(map[int64][]db.LossHigh)

	for _, freq := range freqs {
		freqMap[freq.ThreatEventID] = freq
	}

	for _, loss := range losses {
		lossMap[loss.ThreatEventID] = append(lossMap[loss.ThreatEventID], loss)
	}

	for id, freq := range freqMap {
		if lossList, exists := lossMap[id]; exists {
			for _, loss := range lossList {
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

func PullAllRisk(c *gin.Context) {
	var Risks []db.RiskCalculation
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &Risks); err != nil {
		c.Set("Error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", Risks)
	c.Status(http.StatusOK)
}

func PullRiskId(c *gin.Context, id int) {
	var Risk db.RiskCalculation
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &Risk, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving Risk")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Risk not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", Risk)
	c.Status(http.StatusOK)
}
