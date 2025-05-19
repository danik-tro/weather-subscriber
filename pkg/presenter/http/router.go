package http

import (
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	handlers "github.com/danik-tro/weather-subscriber/pkg/presenter/http/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(subscribeUC usecase.SubscribeWeatherUseCase, getWeatherUC usecase.GetWeatherUseCase, confirmUC usecase.ConfirmSubscriptionUseCase, unsubscribeUC usecase.UnsubscribeUseCase, checkTokensUC usecase.CheckTokensUseCase) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

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
