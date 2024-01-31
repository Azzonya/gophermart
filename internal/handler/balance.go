package handler

import (
	userModel "github.com/Azzonya/gophermart/internal/domain/user/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserHandlers) GetBalance(c *gin.Context) {
	// Реализация получения баланса баллов лояльности пользователя
	c.Header("Content-Type", "application/json")

	userID, err := u.auth.GetUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get cookie",
			"error":   err.Error(),
		})
		return
	}

	result, err := u.userUsecase.GetBalanceWithWithdrawn(c.Request.Context(), &userModel.GetPars{
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
