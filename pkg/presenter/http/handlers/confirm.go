package http

import (
	"net/http"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ConfirmHandler(uc usecase.ConfirmSubscriptionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'token' is required"})
			return
		}

		if _, err := uuid.Parse(token); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
			return
		}

		err := uc.Confirm(c.Request.Context(), token)
		if err != nil {
			switch err {
			case domain.ErrSubscriptionNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed"})
	}
}
