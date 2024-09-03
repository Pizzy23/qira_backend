package simulation

import (
	"errors"
	"qira/db"
	"reflect"
	"strings"

	"xorm.io/xorm"
)

func retrieveFrequencyAndLossEntries(dbEngine *xorm.Engine, threatEvent string, lossType string) ([]db.Frequency, []db.LossHighTotal, error) {
	var frequencies []db.Frequency
	var losses []db.LossHighTotal

	threatEvent = strings.ReplaceAll(threatEvent, "\xa0", " ")

	err := dbEngine.Where("threat_event = ?", threatEvent).Find(&frequencies)
	if err != nil {
		return nil, nil, err
	}

	if lossType == "Singular" {
		err = dbEngine.Where("threat_event = ?", threatEvent).Find(&losses)
		if err != nil {
			return nil, nil, err
		}
	} else if lossType == "LossHigh" {
		err := dbEngine.Where("threat_event = ? AND type_of_loss = ?", threatEvent, "LossHigh").Find(&losses)
		if err != nil {
			return nil, nil, err
		}
	} else if lossType == "Granular" {
		err := dbEngine.Where("threat_event = ? AND type_of_loss = ?", threatEvent, "Granular").Find(&losses)
		if err != nil {
			return nil, nil, err
		}
	}

	return frequencies, losses, nil
}

func validateSimulationData(input interface{}) error {
	val := reflect.ValueOf(input)

	var hasZeroFreq, hasZeroLoss, hasZeroProposed bool

	minFreq := val.FieldByName("FrequencyMin")
	pertFreq := val.FieldByName("FrequencyEstimate")
	maxFreq := val.FieldByName("FrequencyMax")

	minLoss := val.FieldByName("LossMin")
	pertLoss := val.FieldByName("LossEstimate")
	maxLoss := val.FieldByName("LossMax")

	proposedMin := val.FieldByName("ProposedMin")
	proposedPert := val.FieldByName("ProposedPert")
	proposedMax := val.FieldByName("ProposedMax")

	if minFreq.IsValid() && pertFreq.IsValid() && maxFreq.IsValid() {
		if minFreq.Float() == 0 && pertFreq.Float() == 0 && maxFreq.Float() == 0 {
			hasZeroFreq = true
		} else if minFreq.Float() == 0 && maxFreq.Float() == 0 {
			return errors.New("min and Max Frequency values cannot both be zero")
		} else if pertFreq.Float() == 0 {
			return errors.New("pert Frequency value cannot be zero")
		}
	}

	if minLoss.IsValid() && pertLoss.IsValid() && maxLoss.IsValid() {
		if minLoss.Float() == 0 && pertLoss.Float() == 0 && maxLoss.Float() == 0 {
			hasZeroLoss = true
		} else if minLoss.Float() == 0 && maxLoss.Float() == 0 {
			return errors.New("min and Max Loss values cannot both be zero")
		} else if pertLoss.Float() == 0 {
			return errors.New("pert Loss value cannot be zero")
		}
	}

	if proposedMin.IsValid() && proposedPert.IsValid() && proposedMax.IsValid() {
		if proposedMin.Float() == 0 && proposedPert.Float() == 0 && proposedMax.Float() == 0 {
			hasZeroProposed = true
		} else if proposedMin.Float() == 0 && proposedMax.Float() == 0 {
			return errors.New("min and Max Loss values cannot both be zero")
		} else if proposedPert.Float() == 0 {
			return errors.New("pert Proposed value cannot be zero")
		}
	}

	if hasZeroFreq && hasZeroLoss && hasZeroProposed {
		return errors.New("all frequency, loss, and proposed values cannot be zero")
	}

	return nil
}

func calculationRisk(risks []db.RiskCalculation, events []db.ThreatEventCatalog, lossType string) (OutputProcess, error) {
	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	if len(events) == 0 {
		return OutputProcess{}, errors.New("not have events")
	}

	if len(risks) != 0 {
		for _, event := range events {
			for _, risk := range risks {
				if risk.RiskType == "Frequency" && risk.Categorie == lossType && risk.ThreatEvent == event.ThreatEvent {
					totalMinFreq += risk.Min
					totalPertFreq += risk.Estimate
					totalMaxFreq += risk.Max
				} else if risk.Categorie == lossType && risk.ThreatEvent == event.ThreatEvent {
					totalMinLoss += risk.Min
					totalPertLoss += risk.Estimate
					totalMaxLoss += risk.Max
				}
			}
		}
	}

	finalResponse := OutputProcess{
		FrequencyMax:      totalMaxFreq,
		FrequencyMin:      totalMinFreq,
		FrequencyEstimate: totalPertFreq,
		LossMax:           totalMaxLoss,
		LossMin:           totalMinLoss,
		LossEstimate:      totalPertLoss,
	}
	return finalResponse, nil
}

func calculationLossAndFreq(frequencies []db.Frequency, losses []db.LossHighTotal) (OutputProcess, error) {
	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	if len(frequencies) == 0 {
		return OutputProcess{}, errors.New("not have frequencies")
	}

	if len(losses) == 0 {
		return OutputProcess{}, errors.New("not have losses")
	}

	for _, freq := range frequencies {
		totalMinFreq += freq.MinFrequency
		totalPertFreq += freq.MostLikelyFrequency
		totalMaxFreq += freq.MaxFrequency
	}

	for _, loss := range losses {
		totalMinLoss += loss.MinimumLoss
		totalPertLoss += loss.MostLikelyLoss
		totalMaxLoss += loss.MaximumLoss
	}
	finalResponse := OutputProcess{
		FrequencyMax:      totalMaxFreq,
		FrequencyMin:      totalMinFreq,
		FrequencyEstimate: totalPertFreq,
		LossMax:           totalMaxLoss,
		LossMin:           totalMinLoss,
		LossEstimate:      totalPertLoss,
	}
	return finalResponse, nil
}
