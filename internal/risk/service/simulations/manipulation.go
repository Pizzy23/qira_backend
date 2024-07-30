package simulation

import (
	"errors"
	"qira/db"

	"xorm.io/xorm"
)

func retrieveFrequencyAndLossEntries(engine *xorm.Engine, threatEvent, lossType string) ([]db.Frequency, []db.LossHighTotal, error) {
	var frequencyEntries []db.Frequency
	var lossEntries []db.LossHighTotal

	err := engine.Where("threat_event = ?", threatEvent).Find(&frequencyEntries)
	if err != nil {
		return nil, nil, errors.New("error retrieving frequency entries")
	}

	err = engine.Where("threat_event = ? AND type_of_loss = ?", threatEvent, lossType).Find(&lossEntries)
	if err != nil {
		return nil, nil, errors.New("error retrieving loss entries")
	}

	if lossZero(lossEntries) {
		return nil, nil, errors.New("all loss values are zero")
	}

	return frequencyEntries, lossEntries, nil
}

func lossZero(loss []db.LossHighTotal) bool {
	for _, l := range loss {
		if l.MinimumLoss != 0 || l.MaximumLoss != 0 || l.MostLikelyLoss != 0 {
			return false
		}
	}
	return true
}
