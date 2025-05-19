package http

import (
	config "github.com/danik-tro/weather-subscriber/pkg"
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	handlers "github.com/danik-tro/weather-subscriber/pkg/presenter/http/handlers"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(config config.Config, subscribeUC usecase.SubscribeWeatherUseCase, getWeatherUC usecase.GetWeatherUseCase, confirmUC usecase.ConfirmSubscriptionUseCase, unsubscribeUC usecase.UnsubscribeUseCase, checkTokensUC usecase.CheckTokensUseCase) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	url := ginSwagger.URL(config.SwaggerURL)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	api := router.Group("/api")
	{
		api.POST("/subscribe", handlers.SubscribeHandler(subscribeUC))
		api.GET("/weather", handlers.GetWeatherHandler(getWeatherUC))
		api.POST("/confirm/:token", handlers.ConfirmHandler(confirmUC))
		api.POST("/unsubscribe/:token", handlers.UnsubscribeHandler(unsubscribeUC))
	}

	router.GET("/confirm/:token", handlers.CheckConfirmationTokenHandler(checkTokensUC))
	router.GET("/unsubscribe/:token", handlers.CheckUnsubscribeTokenHandler(checkTokensUC))

	return router
}
