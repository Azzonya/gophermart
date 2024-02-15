package handler

import (
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserHandlers) GetBalance(c *gin.Context) {
	// Реализация получения баланса баллов лояльности пользователя
	userID, _ := u.userUsecase.GetUserIDFromCookieU(c)

	result, err := u.userUsecase.GetBalanceWithWithdrawn(c.Request.Context(), &entities.UserParameters{
		ID: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get balance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
