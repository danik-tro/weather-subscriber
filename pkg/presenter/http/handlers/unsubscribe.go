package http

import (
	"net/http"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Unsubscribe from weather updates
// @Description Unsubscribe from weather updates using token
// @Tags subscription
// @Accept json
// @Produce json
// @Param token path string true "Unsubscribe token"
// @Success 200 {object} map[string]string "Subscription unsubscribed"
// @Failure 400 {object} map[string]string "Invalid token format or missing token"
// @Failure 404 {object} map[string]string "Subscription not found"
// @Router /unsubscribe/{token} [post]
func UnsubscribeHandler(uc usecase.UnsubscribeUseCase) gin.HandlerFunc {
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

		err := uc.Unsubscribe(c.Request.Context(), token)
		if err != nil {
			switch err {
			case domain.ErrSubscriptionNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription unsubscribed"})
	}
}
