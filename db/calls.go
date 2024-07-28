package db

import (
	"fmt"

	"xorm.io/xorm"
)

func Create(engine *xorm.Engine, table interface{}) error {
	_, err := engine.Insert(table)
	return err
}

func Read(engine *xorm.Engine, table interface{}, condition interface{}) error {
	_, err := engine.Where(condition).Get(table)
	return err
}

func UpdateByThreatEvent(engine *xorm.Engine, table interface{}, threatEventID int64, riskType string) error {
	_, err := engine.Where("threat_event_i_d = ? AND risk_type = ?", threatEventID, riskType).Update(table)
	return err
}

func UpdateByThreat(engine *xorm.Engine, table interface{}, threatEventID int64) error {
	_, err := engine.Where("threat_event_i_d = ?", threatEventID).Update(table)
	return err
}
func UpdateByControlIdAndRisk(engine *xorm.Engine, table interface{}, controlID int64, attack string) error {
	_, err := engine.Where("control_i_d = ? AND type_of_attack = ?", controlID, attack).Update(table)
	return err
}
func UpdateByControlId(engine *xorm.Engine, table interface{}, controlID int64) error {
	_, err := engine.Where("control_i_d = ?", controlID).Update(table)
	return err
}

func UpdateByID(engine *xorm.Engine, table interface{}, id int64) error {
	_, err := engine.ID(id).Update(table)
	return err
}

func Delete(engine *xorm.Engine, table interface{}, condition interface{}) error {
	_, err := engine.Where(condition).Delete(table)
	return err
}

func GetByID(engine *xorm.Engine, table interface{}, id int64) (bool, error) {
	found, err := engine.ID(id).Get(table)
	return found, err
}

func GetAll(engine *xorm.Engine, tableSlice interface{}) error {
	err := engine.Find(tableSlice)
	return err
}

func GetRiskCalculationsByRiskType(engine *xorm.Engine, threatEventID int64, riskType string) ([]RiskCalculation, error) {
	var riskCalcs []RiskCalculation
	err := engine.Where("threat_event_i_d = ? AND risk_type = ?", threatEventID, riskType).Find(&riskCalcs)
	if err != nil {
		return nil, err
	}
	return riskCalcs, nil
}

func CreateColumn(engine *xorm.Engine, tableName string, columnName string, typeTable string) error {
	query := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s DEFAULT 0", tableName, columnName, typeTable)
	_, err := engine.Exec(query)
	return err
}

func GetRiskById(engine *xorm.Engine, riskType int64) ([]RiskCalculation, error) {
	var riskCalcs []RiskCalculation
	err := engine.Where("threat_event_i_d = ?", riskType).Find(&riskCalcs)
	if err != nil {
		return nil, err
	}
	return riskCalcs, nil
}

func GetByEventIDAndRiskType(engine *xorm.Engine, table interface{}, eventId int64, riskType string) (bool, error) {
	found, err := engine.Where("threat_event_i_d = ? AND risk_type = ?", eventId, riskType).Get(table)
	return found, err
}

func GetAllWithCondition(engine *xorm.Engine, tableSlice interface{}, condition string, args ...interface{}) error {
	return engine.Where(condition, args...).Find(tableSlice)
}
