package environment

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateRiskAssessmentService(c *gin.Context, risk interfaces.InputRiskAssessment) error {
	newRisk := db.RiskAssessment{
		RiskAssessmentName:       risk.RiskAssessmentName,
		Practitioner:             risk.Practitioner,
		AssessmentStartDate:      risk.AssessmentStartDate,
		AssessmentEndDate:        risk.AssessmentEndDate,
		Sponsor:                  risk.Sponsor,
		TargetEnvironment:        risk.TargetEnvironment,
		Profit:                   risk.Profit,
		AnnualRevenue:            risk.AnnualRevenue,
		StockPrice:               risk.StockPrice,
		IndustrySector:           risk.IndustrySector,
		GeographicRegion:         risk.GeographicRegion,
		NumberOfOperationalSites: risk.NumberOfOperationalSites,
		NumberOfEmployees:        risk.NumberOfEmployees,
	}
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &newRisk); err != nil {
		return err
	}
	return nil
}

func PullAllRiskAssessments(c *gin.Context) {
	var risks []db.RiskAssessment
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &risks); err != nil {
		c.Set("Response", "Error retrieving risk assessments: "+err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if risks == nil {
		var empty []string
		c.Set("Response", empty)
		c.Status(http.StatusOK)
		return
	}
	c.Set("Response", risks)
	c.Status(http.StatusOK)
}

func PullRiskAssessmentById(c *gin.Context, id int) {
	var risk db.RiskAssessment
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &risk, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving risk assessment")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "Risk assessment not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", risk)
	c.Status(http.StatusOK)
}

func UpdateRiskAssessmentService(c *gin.Context, id int64, risk interfaces.InputRiskAssessment) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	riskToUpdate := db.RiskAssessment{
		RiskAssessmentName:       risk.RiskAssessmentName,
		Practitioner:             risk.Practitioner,
		AssessmentStartDate:      risk.AssessmentStartDate,
		AssessmentEndDate:        risk.AssessmentEndDate,
		Sponsor:                  risk.Sponsor,
		TargetEnvironment:        risk.TargetEnvironment,
		Profit:                   risk.Profit,
		AnnualRevenue:            risk.AnnualRevenue,
		StockPrice:               risk.StockPrice,
		IndustrySector:           risk.IndustrySector,
		GeographicRegion:         risk.GeographicRegion,
		NumberOfOperationalSites: risk.NumberOfOperationalSites,
		NumberOfEmployees:        risk.NumberOfEmployees,
	}

	if err := db.UpdateByID(engine.(*xorm.Engine), &riskToUpdate, id); err != nil {
		return err
	}
	return nil
}

func DeleteRiskAssessment(c *gin.Context, id int64) error {
	var risk db.RiskAssessment

	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &risk, id)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("risk assessment not found")
	}

	if err := db.Delete(engine.(*xorm.Engine), &db.RiskAssessment{}, map[string]interface{}{"id": id}); err != nil {
		return err
	}

	return nil
}
