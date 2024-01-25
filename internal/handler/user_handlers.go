package handler

import (
	"context"
	"fmt"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order/model"
	userModel "github.com/Azzonya/gophermart/internal/domain/user/model"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"github.com/Azzonya/gophermart/internal/usecase/user"
	"github.com/Azzonya/gophermart/internal/usecase/withdrawal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandlers struct {
	userUsecase  *user.Usecase
	orderUsecase *order.Usecase
	withdrawal   *withdrawal.Usecase
}

func New() *UserHandlers {
	return &UserHandlers{}
}

func (u *UserHandlers) RegisterUser(c *gin.Context) {
	var err error
	ctx := context.Background()

	req := &userModel.GetPars{}

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
		c.JSON(http.StatusInternalServerError, gin.H{
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
	req := &userModel.GetPars{}

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
	orderNumber := c.PostForm("orderNumber")

	if !u.orderUsecase.IsLuhnValid(orderNumber) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Неверный формат номера заказа (алгоритм Луна)"})
		return
	}

	order, orderExist, err := u.orderUsecase.Get(context.Background(), &orderModel.GetPars{
		OrderNumber: orderNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get order",
			"error":   err.Error(),
		})
		return
	}

	login, err := c.Cookie("login")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get login",
			"error":   err.Error(),
		})
		return
	}

	user, _, err := u.userUsecase.Get(context.Background(),
		&userModel.GetPars{
			Login: login,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get login",
			"error":   err.Error(),
		})
		return
	}

	if orderExist {
		if user.ID == order.UserID {
			c.JSON(http.StatusOK, gin.H{
				"message": "already uploaded",
			})
			return
		} else {
			c.JSON(http.StatusConflict, gin.H{
				"message": "already uploaded by different user",
			})
			return
		}
	}

	err = u.orderUsecase.Create(context.Background(), &orderModel.GetPars{
		OrderNumber: orderNumber,
		Status:      orderModel.OrderStatusNew,
		UserID:      user.ID,
	})

	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) GetOrders(c *gin.Context) {
	// Реализация получения списка заказов пользователя
	login, err := c.Cookie("login")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get login",
			"error":   err.Error(),
		})
		return
	}

	user, _, err := u.userUsecase.Get(context.Background(),
		&userModel.GetPars{
			Login: login,
		})

	orders, err := u.orderUsecase.List(context.Background(), &orderModel.ListPars{
		UserID: &user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get login",
			"error":   err.Error(),
		})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	result := []*ListOrdersResult{}

	for _, v := range orders {
		newElement := ListOrdersResult{}
		newElement.Encode(v)

		result = append(result, &newElement)
	}

	c.JSON(http.StatusOK, result)
}

func (u *UserHandlers) GetBalance(c *gin.Context) {
	// Реализация получения баланса баллов лояльности пользователя
	login, err := c.Cookie("login")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get login",
			"error":   err.Error(),
		})
		return
	}

	user, _, err := u.userUsecase.Get(context.Background(),
		&userModel.GetPars{
			Login: login,
		})
	fmt.Println(user)

	c.JSON(http.StatusOK, UserBalanceResult{ // change
		Balance:   0,
		Withdrawn: 0,
	})
}

func (u *UserHandlers) WithdrawBalance(c *gin.Context) {
	// Реализация запроса на списание баллов
	req := WithdrawalBalanceRequest{}

	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	login, err := c.Cookie("login")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get login",
			"error":   err.Error(),
		})
		return
	}

	user, _, err := u.userUsecase.Get(context.Background(),
		&userModel.GetPars{
			Login: login,
		})

	if user.Balance < req.Sum {
		c.JSON(http.StatusPaymentRequired, nil)
		return
	}

	_, orderExist, err := u.orderUsecase.Get(context.Background(), &orderModel.GetPars{
		OrderNumber: req.OrderNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get order",
			"error":   err.Error(),
		})
		return
	}
	if !orderExist {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) GetWithdrawals(c *gin.Context) {
	// Реализация получения информации о выводе средств
}
