package handler

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
	"github.com/Azzonya/gophermart/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandlers struct {
	userUsecase *user.Usecase
}

func New() *UserHandlers {
	return &UserHandlers{}
}

func (u *UserHandlers) RegisterUser(c *gin.Context) {
	var err error
	ctx := context.Background()

	req := &model.User{}

	err = c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	exist, err := u.userUsecase.IsLoginTaken(ctx, req.Login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to check login",
			"error":   err.Error(),
		})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Login is already taken",
		})
		return
	}

	if err = u.userUsecase.Register(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	c.Header("Set-Cookie", req.Login)
	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) LoginUser(c *gin.Context) {
	// Реализация аутентификации пользователя
	var err error
	req := &model.GetPars{}

	err = c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	exist, err := u.userUsecase.CheckAuth(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check auth",
			"error":   err.Error(),
		})
		return
	}
	if !exist {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	c.Header("Set-Cookie", req.Login)
	c.JSON(http.StatusOK, nil)
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
