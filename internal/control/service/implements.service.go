package control

import (
	"errors"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	"qira/internal/mock"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func CreateOrUpdateImplementService(c *gin.Context, data interfaces.ImplementsInputNoID, id int64) error {
	averageC, err := mock.FindAverageByScoreImp(data.Current)
	if err != nil {
		return errors.New("score not found for Percent Current")
	}
	averageP, err := mock.FindAverageByScoreImp(data.Proposed)
	if err != nil {
		return errors.New("score not found for Percent Proposed")
	}

	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingImplement db.Implements
	has, err := engine.(*xorm.Engine).Where("control_id = ?", id).Get(&existingImplement)
	if err != nil {
		return err
	}

	if !has {
		var control db.Control
		hasControl, err := engine.(*xorm.Engine).Where("id = ?", id).Get(&control)
		if err != nil {
			return err
		}
		if !hasControl {
			return errors.New("ControlID not found")
		}

		newImplement := db.Implements{
			ControlID:       id,
			Current:         data.Current,
			Proposed:        data.Proposed,
			Cost:            data.Cost,
			PercentCurrent:  averageC,
			PercentProposed: averageP,
		}
		if err := db.Create(engine.(*xorm.Engine), &newImplement); err != nil {
			return err
		}
		return nil
	}

	if existingImplement.Current == data.Current &&
		existingImplement.Proposed == data.Proposed &&
		existingImplement.Cost == data.Cost &&
		existingImplement.PercentCurrent == averageC &&
		existingImplement.PercentProposed == averageP {
		return errors.New("no update made, all fields are the same")
	}

	updatedImplement := db.Implements{
		ControlID:       id,
		Current:         data.Current,
		Proposed:        data.Proposed,
		Cost:            data.Cost,
		PercentCurrent:  averageC,
		PercentProposed: averageP,
	}
	if err := db.UpdateByControlId(engine.(*xorm.Engine), &updatedImplement, id); err != nil {
		return err
	}

	return nil
}

func PullAllImplements(c *gin.Context) {
	var controls []db.ControlLibrary
	var implementations []db.Implements

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &controls); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &implementations); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	implementationMap := make(map[int64]bool)
	for _, impl := range implementations {
		implementationMap[impl.ControlID] = true
	}

	var newImplementations []db.Implements
	for _, control := range controls {
		if !implementationMap[control.ID] {
			newImpl := db.Implements{
				ControlID:       control.ID,
				Current:         0,
				Proposed:        0,
				PercentCurrent:  "3%",
				PercentProposed: "3%",
				Cost:            0,
			}
			newImplementations = append(newImplementations, newImpl)
		}
	}

	if len(newImplementations) > 0 {
		if _, err := engine.(*xorm.Engine).Insert(&newImplementations); err != nil {
			c.Set("Response", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	if err := db.GetAll(engine.(*xorm.Engine), &implementations); err != nil {
		c.Set("Response", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("Response", implementations)
	c.Status(http.StatusOK)
}

func PullControlId(c *gin.Context, id int) {
	var control db.Implements
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &control, int64(id))
	if err != nil {
		c.Set("Response", "Error retrieving control")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Response", "control not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", control)
	c.Status(http.StatusOK)
}
