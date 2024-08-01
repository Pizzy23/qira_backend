package middleware

import (
	"qira/db"
	_ "qira/docs"
	catalogue "qira/internal/catalogue/handler"
	control "qira/internal/control/handler"
	event "qira/internal/event/handler"
	frequency "qira/internal/frequency/handler"
	revelance "qira/internal/revelance/handler"
	risk "qira/internal/risk/handler"

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
	r.Use(func(c *gin.Context) {
		c.Set("db", db.Repo)
		c.Next()
	})

	r.GET("/simulation", risk.RiskMount)
	r.GET("/simulation-aggregated", risk.RiskMountAggregated)
	r.GET("/simulation-appetite", risk.RiskMountAppetite)
	r.GET("/simulation-report", risk.RiskMountReport)

	r.Use(CORSConfig())
	r.Use(ResponseHandler())

	//Use response, but not Token

	auth := r.Group("/api")
	//Response and token service

	//invetory
	auth.GET("/asset/:id", inventory.PullAssetId)
	auth.POST("/create-asset", inventory.CreateAsset)
	auth.PUT("/asset/:id", inventory.UpdateAsset)
	auth.GET("/assets", inventory.PullAllAsset)
	auth.DELETE("/asset/:id", inventory.DeleteAsset)

	//frequency
	auth.GET("/frequency/:id", frequency.PullFrequencyById)
	auth.GET("/all-frequency", frequency.PullAllFrequency)
	auth.PUT("/frequency/:id", frequency.EditFrequency)

	//events
	auth.GET("/all-event", event.PullAllForEvent)
	auth.PUT("/event/:id", event.CreateEvent)
	auth.DELETE("/event/:id", event.DeleteControlId)

	//Catalogue
	auth.GET("/all-catalogue", catalogue.PullAllEvent)
	auth.GET("/catalogue/:id", catalogue.PullEventId)
	auth.POST("/catalogue", catalogue.CreateEvent)
	auth.DELETE("/catalogue/:id", catalogue.DeleteEventId)

	//losshigh
	auth.GET("/losshigh/:id", losshigh.PullLosstId)
	auth.PUT("/update-losshigh/:id", losshigh.CreateLossHigh)
	auth.GET("/losshigh-singular", losshigh.PullAllLossHighSingular)
	auth.PUT("/update-losshigh-singular/:id", losshigh.CreateLossHighSingular)
	auth.GET("/losshigh-granuled", losshigh.PullAllLossHighGranuled)
	auth.PUT("/update-losshigh-granuled/:id", losshigh.CreateLossHighGranuled)
	auth.GET("/losshigh", losshigh.PullAllLossHigh)
	auth.POST("/losshigh-specific", losshigh.CreateLossHighSpecific)

	//Control
	auth.PUT("/control/:id", control.UpdateControl)
	auth.GET("/all-control", control.PullAllControl)
	auth.GET("/control/:id", control.PullControlId)
	auth.DELETE("/control/:id", control.DeleteControlId)
	auth.POST("/control", control.CreateControl)

	//Implementation.
	auth.GET("/all-implementation", control.PullAllControlImplementation)
	auth.GET("/implementation/:id", control.PullControlImplementationId)
	auth.PUT("/implementation/:id", control.EditControlImplementation)

	//Risk
	auth.GET("/risk", risk.PullAllRisk)
	auth.GET("/risk/:id", risk.PullRiskId)

	// Revelance
	auth.GET("/revelance", revelance.PullAllRevelance)
	auth.GET("/revelance/:id", revelance.PullRevelanceId)
	auth.PUT("/update-revelance", revelance.UpdateRelevance)

	auth.GET("/aggregated-control-strength", control.PullAggregatedControlStrength)
	auth.GET("/all-proposed", control.PullAllControlProposed)
	auth.GET("/all-strength", control.PullAllControlStrength)
	auth.GET("/all-prupu", control.PullPrupu)
	auth.GET("/all-stren", control.PullStren)

	auth.PUT("/upload-appetite", risk.UploadAppetite)
	return r
}
