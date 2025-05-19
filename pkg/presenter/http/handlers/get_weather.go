package http

import (
	"net/http"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	"github.com/gin-gonic/gin"
)

// @Summary Get weather by city
// @Description Get current weather information for a specific city
// @Tags weather
// @Accept json
// @Produce json
// @Param city query string true "City name"
// @Success 200 {object} domain.Weather "Weather information"
// @Failure 400 {object} map[string]string "Invalid request or missing city parameter"
// @Failure 404 {object} map[string]string "City not found"
// @Router /weather [get]
func GetWeatherHandler(uc usecase.GetWeatherUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Query("city")
		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'city' is required"})
			return
		}

		res, err := uc.GetWeather(c.Request.Context(), city)
		if err != nil {
			switch err {
			case domain.ErrCityNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
