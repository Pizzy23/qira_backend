package main

import (
	"log"
	"os"
	"qira/db"
	"qira/docs"
	_ "qira/docs" // import generated docs
	"qira/middleware"
)

// @title           Qira
// @version         2.0
// @description     This is a server for app.

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	r := middleware.SetupRouter()

	db.ConnectDatabaseXorm()
	//migrate()

	host := os.Getenv("API_HOST")
	if host == "" {
		host = "localhost:8080"
	}

	docs.SwaggerInfo.Host = host

	r.Run(":8080")
}

func migrate() {
	if err := db.Migrate(db.Repo); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}
