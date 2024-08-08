package simulation

import (
	"qira/db"

	"xorm.io/xorm"
)

func lossZero(loss []db.LossHighTotal) bool {
	for _, l := range loss {
		if l.MinimumLoss != 0 || l.MaximumLoss != 0 || l.MostLikelyLoss != 0 {
			return false
		}
	}
	return true
}

func retrieveFrequencyAndLossEntries(dbEngine *xorm.Engine, threatEvent string, lossType string) ([]db.Frequency, []db.LossHigh, error) {
	var frequencies []db.Frequency
	var losses []db.LossHigh

	if lossType == "Singular" {
		err := dbEngine.Where("threat_event = ?", threatEvent).Find(&frequencies)
		if err != nil {
			return nil, nil, err
		}

		err = dbEngine.Where("threat_event = ?", threatEvent).Find(&losses)
		if err != nil {
			return nil, nil, err
		}
	} else if lossType == "LossHigh" {
		err := dbEngine.Where("threat_event = ? AND loss_type = ?", threatEvent, "LossHigh").Find(&losses)
		if err != nil {
			return nil, nil, err
		}
	} else if lossType == "Granular" {
		err := dbEngine.Where("threat_event = ? AND loss_type = ?", threatEvent, "Granular").Find(&losses)
		if err != nil {
			return nil, nil, err
		}
	}

	return frequencies, losses, nil
}
