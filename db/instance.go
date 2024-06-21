package db

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"xorm.io/core"
	"xorm.io/xorm"
)

var Repo *xorm.Engine

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DB")
	if databaseURL == "" {
		log.Fatal("DB environment variable not set")
	}

	engine, err := xorm.NewEngine("mysql", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	engine.SetTableMapper(core.SnakeMapper{})
	engine.SetColumnMapper(core.SameMapper{})

	engine.ShowSQL(true)

	err = Migrate(engine)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	Repo = engine
}
