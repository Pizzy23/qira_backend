package revelance

import (
	"errors"
	"fmt"
	"net/http"
	"qira/db"
	"qira/internal/interfaces"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func PullAllRevelance(c *gin.Context) {
	var revelances []db.RelevanceDinamic
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := db.GetAll(engine.(*xorm.Engine), &revelances); err != nil {
		c.Set("Error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", revelances)
	c.Status(http.StatusOK)
}

func PullRevelanceId(c *gin.Context, id int) {
	var revelance db.RelevanceDinamic
	engine, exists := c.Get("db")
	if !exists {
		c.Set("Error", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := db.GetByID(engine.(*xorm.Engine), &revelance, int64(id))
	if err != nil {
		c.Set("Error", "Error retrieving revelance")
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Set("Error", "Revelance not found")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Set("Response", revelance)
	c.Status(http.StatusOK)
}

func CreateRelevanceService(c *gin.Context, Relevance db.RelevanceDinamicInput) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	if err := db.Create(engine.(*xorm.Engine), &Relevance); err != nil {
		return err
	}
	return nil

}

func UpdateRelevanceService(c *gin.Context, id int64, relevanceInput interfaces.RelevanceDinamicInput) error {
	engine, exists := c.Get("db")
	if !exists {
		return errors.New("database connection not found")
	}

	var existingRelevance db.RelevanceDinamic
	if _, err := engine.(*xorm.Engine).ID(id).Get(&existingRelevance); err != nil {
		return err
	}

	v := reflect.ValueOf(&existingRelevance).Elem()
	updateMap := map[string]interface{}{}

	for field, value := range relevanceInput.Attributes {
		fieldName := toCamelCase(field)
		f := v.FieldByName(fieldName)
		if f.IsValid() && f.CanSet() {
			updateMap[strings.ToLower(fieldName)] = value
		} else {
			return fmt.Errorf("field %s does not exist in RelevanceDinamic", field)
		}
	}

	if _, err := engine.(*xorm.Engine).Table(new(db.RelevanceDinamic)).ID(id).Update(updateMap); err != nil {
		return err
	}
	return nil
}

func toCamelCase(input string) string {
	isToUpper := false
	camelCase := ""
	for i, char := range input {
		if i == 0 {
			camelCase += string(char - 32)
		} else {
			if isToUpper {
				camelCase += string(char - 32)
				isToUpper = false
			} else if char == '_' {
				isToUpper = true
			} else {
				camelCase += string(char)
			}
		}
	}
	return camelCase
}
