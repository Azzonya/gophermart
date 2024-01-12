package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandlers struct {
}

func New() *UserHandlers {
	return &UserHandlers{}
}

func (u *UserHandlers) RegisterUser(c *gin.Context) {
	// Реализация регистрации пользователя
	var err error

	req := &RegisterRequest{}

	err = c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

}

func (u *UserHandlers) LoginUser(c *gin.Context) {
	// Реализация аутентификации пользователя
}

func (u *UserHandlers) UploadOrder(c *gin.Context) {
	// Реализация загрузки номера заказа
}

func (u *UserHandlers) GetOrders(c *gin.Context) {
	// Реализация получения списка заказов пользователя
}

func (u *UserHandlers) GetBalance(c *gin.Context) {
	// Реализация получения баланса баллов лояльности пользователя
}

func (u *UserHandlers) WithdrawBalance(c *gin.Context) {
	// Реализация запроса на списание баллов
}

func (u *UserHandlers) GetWithdrawals(c *gin.Context) {
	// Реализация получения информации о выводе средств
}
