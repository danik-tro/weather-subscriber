package http

import (
	"errors"
	"fmt"
	"net/http"

	"strings"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	entity "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SubscribeRequest struct {
	Email     string `form:"email" json:"email" binding:"required,email"`
	City      string `form:"city"  json:"city"  binding:"required"`
	Frequency string `form:"frequency" json:"frequency" binding:"required,oneof=hourly daily"`
}

func SubscribeHandler(uc usecase.SubscribeWeatherUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SubscribeRequest
		if err := c.ShouldBind(&req); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				out := make([]string, len(ve))
				for i, fe := range ve {
					out[i] = fmt.Sprintf("Field '%s' failed on the '%s' rule", fe.Field(), fe.Tag())
				}
				c.JSON(http.StatusBadRequest, gin.H{"errors": out})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		freq := entity.Frequency(strings.ToUpper(req.Frequency))

		err := uc.Subscribe(c.Request.Context(), req.Email, req.City, freq)

		if err != nil {
			switch err {
			case domain.ErrSubscriptionAlreadyExists:
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			case domain.ErrCityNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "subscription created, confirmation email sent"})
	}
}
