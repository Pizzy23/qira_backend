package db

import (
	"xorm.io/xorm"
)

func Migrate(engine *xorm.Engine) error {
	if err := engine.Sync2(new(AssetsInventory)); err != nil {
		return err
	}
	if err := engine.Sync2(new(ThreatEventCatalogue)); err != nil {
		return err
	}
	if err := engine.Sync2(new(Frequency)); err != nil {
		return err
	}
	if err := engine.Sync2(new(ThreatEventAssets)); err != nil {
		return err
	}
	if err := engine.Sync2(new(LossHigh)); err != nil {
		return err
	}
	if err := engine.Sync2(new(RiskCalculator)); err != nil {
		return err
	}
	if err := engine.Sync2(new(Relevance)); err != nil {
		return err
	}
	if err := engine.Sync2(new(ControlLibrary)); err != nil {
		return err
	}
	if err := engine.Sync2(new(ControlImplementation)); err != nil {
		return err
	}
	if err := engine.Sync2(new(AggregatedControlStrength)); err != nil {
		return err
	}
	return nil
}
