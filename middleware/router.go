package middleware

import (
	"qira/db"
	_ "qira/docs"
	catalogue "qira/internal/catalogue/handler"
	control "qira/internal/control/handler"
	event "qira/internal/event/handler"
	frequency "qira/internal/frequency/handler"

	inventory "qira/internal/inventory/handler"
	losshigh "qira/internal/loss-high/handler"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Qira
// @version 1.0
// @description API
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(CORSConfig())
	r.Use(ResponseHandler())

	r.Use(func(c *gin.Context) {
		c.Set("db", db.Repo)
		c.Next()
	})

	//Use response, but not Token
	r.GET("/token", generateTokenHandler)

	auth := r.Group("/api")
	auth.Use(authMiddleware)
	//Response and token service

	auth.GET("/asset/:id", inventory.PullAssetId)
	auth.GET("/all-frequency", frequency.PullAllFrequency)
	auth.GET("/assets", inventory.PullAllAsset)
	auth.GET("/frequency/:id", frequency.PullFrequencyById)
	auth.GET("/all-catalogue", catalogue.PullAllEvent)
	auth.GET("/catalogue/:id", catalogue.PullEventId)
	auth.GET("/all-event", event.PullAllForEvent)
	auth.GET("/losshigh/:id", losshigh.PullLosstId)
	auth.GET("/losshigh", losshigh.PullAllLoss)

	auth.GET("/all-relevance", control.PullLibraryControl)
	auth.GET("/all-implementation", control.PullImplementationControl)
	auth.GET("/all-propused", control.PullPropusedControl)
	auth.GET("/all-library", control.PullLibraryControl)
	auth.GET("/all-strength", control.PullStrengthControl)

	auth.GET("/relevance/:id", control.PullRelevanceId)
	auth.GET("/propused/:id", control.PullImplementationId)
	auth.GET("/implementation/:id", control.PullPropusedId)
	auth.GET("/library/:id", control.PullLibraryId)
	auth.GET("/strength/:id", control.PullStrengthId)

	auth.PUT("/frequency", frequency.EditFrequency)

	auth.POST("/catalogue", catalogue.CreateEvent)
	auth.POST("/create-asset", inventory.CreateAsset)
	auth.POST("/event", event.CreateEvent)
	auth.POST("/losshigh", losshigh.CreateLossHigh)
	auth.POST("/create-strength", control.CreateStrength)
	auth.POST("/create-library", control.CreateLibrary)
	auth.POST("/create-propusedry", control.CreatePropused)
	auth.POST("/create-implementation", control.CreateImplementation)
	auth.POST("/create-relevance", control.CreateStrength)
	return r
}
