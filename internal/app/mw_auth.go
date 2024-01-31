package app

import (
	"errors"
	"github.com/Azzonya/gophermart/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/api/user/register" || path == "/api/user/login" {
			c.Next()
			return
		}

		authorizer := auth.New(jwtSecret)

		userID, err := authorizer.GetUserIDFromCookie(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to get cookie",
				"error":   err.Error(),
			})
			return
		}
		if len(userID) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
