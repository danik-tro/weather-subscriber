package http

import (
	"net/http"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckConfirmationTokenHandler(uc usecase.CheckTokensUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")

		if token == "" {
			c.HTML(http.StatusBadRequest, "404.html", gin.H{"Message": "Token is required"})
			return
		}

		if _, err := uuid.Parse(token); err != nil {
			c.HTML(http.StatusBadRequest, "404.html", gin.H{"Message": "Invalid token format"})
			return
		}

		res, err := uc.CheckConfirmationToken(c.Request.Context(), token)
		if err != nil {
			switch err {
			case domain.ErrSubscriptionNotFound:
				c.HTML(http.StatusNotFound, "404.html", gin.H{"Message": "Token not found"})
			default:
				c.HTML(http.StatusBadRequest, "404.html", gin.H{"Message": err.Error()})
			}
			return
		}

		if !res {
			c.HTML(http.StatusNotFound, "404.html", gin.H{"Message": "Token not found"})
			return
		}

		c.HTML(http.StatusOK, "get_confirmation.html", gin.H{"Token": token})
	}
}

func CheckUnsubscribeTokenHandler(uc usecase.CheckTokensUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")

		if token == "" {
			c.HTML(http.StatusBadRequest, "404.html", gin.H{"Message": "Token is required"})
			return
		}

		if _, err := uuid.Parse(token); err != nil {
			c.HTML(http.StatusBadRequest, "404.html", gin.H{"Message": "Invalid token format"})
			return
		}

		res, err := uc.CheckUnsubscribeToken(c.Request.Context(), token)
		if err != nil {
			switch err {
			case domain.ErrSubscriptionNotFound:
				c.HTML(http.StatusNotFound, "404.html", gin.H{"Message": "Token not found"})
			default:
				c.HTML(http.StatusBadRequest, "404.html", gin.H{"Message": err.Error()})
			}
			return
		}

		if !res {
			c.HTML(http.StatusNotFound, "404.html", gin.H{"Message": "Token not found"})
			return
		}

		c.HTML(http.StatusOK, "get_unsubscribe.html", gin.H{"Token": token})
	}
}
