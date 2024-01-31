package handler

import (
	"context"
	"errors"
	"github.com/Azzonya/gophermart/internal/auth"
	bonusTransactionsModel "github.com/Azzonya/gophermart/internal/domain/bonusTransactions/model"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order/model"
	userModel "github.com/Azzonya/gophermart/internal/domain/user/model"
	"github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"github.com/Azzonya/gophermart/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type UserHandlers struct {
	auth                     *auth.Auth
	userUsecase              *user.Usecase
	orderUsecase             *order.Usecase
	bonusTransactionsUsecase *bonustransactions.Usecase
}

func New(auth *auth.Auth, userUsecase *user.Usecase, orderUsecase *order.Usecase, bonusTransactionsUsecase *bonustransactions.Usecase) *UserHandlers {
	return &UserHandlers{
		auth:                     auth,
		userUsecase:              userUsecase,
		orderUsecase:             orderUsecase,
		bonusTransactionsUsecase: bonusTransactionsUsecase,
	}
}

func (u *UserHandlers) RegisterUser(c *gin.Context) {
	var err error
	req := &userModel.GetPars{}

	ctx := c.Request.Context()

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

	newUser, err := u.userUsecase.Register(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	sessionCookie, errS := u.auth.CreateJWTCookie(newUser)
	if errS != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//c.Header("Set-Cookie", sessionCookie.String())
	c.SetCookie("jwt", sessionCookie.String(), int(sessionCookie.Expires.Second()), "/", "", false, true)

	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) LoginUser(c *gin.Context) {
	// Реализация аутентификации пользователя
	var err error
	req := &userModel.GetPars{}

	userID, err := u.auth.GetUserIDFromCookie(c)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get cookie",
			"error":   err.Error(),
		})
		return
	}
	if len(userID) != 0 {
		c.JSON(http.StatusOK, nil)
		return
	}

	err = c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	foundUser, exist, err := u.userUsecase.CheckAuth(c.Request.Context(), req)
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

	sessionCookie, errS := u.auth.CreateJWTCookie(foundUser)
	if errS != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Set-Cookie", sessionCookie.String())
	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) UploadOrder(c *gin.Context) {
	// Реализация загрузки номера заказа
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	ctx := c.Request.Context()

	orderNumber := string(body)

	//if !u.orderUsecase.IsLuhnValid(orderNumber) {
	//	c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Неверный формат номера заказа (алгоритм Луна)"})
	//	return
	//}

	foundOrder, orderExist, err := u.orderUsecase.Get(ctx, &orderModel.GetPars{
		OrderNumber: orderNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get order",
			"error":   err.Error(),
		})
		return
	}
	userID, err := u.auth.GetUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get cookie",
			"error":   err.Error(),
		})
		return
	}

	if orderExist {
		if userID == foundOrder.UserID {
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
		UserID:      userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create order",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) GetOrders(c *gin.Context) {
	// Реализация получения списка заказов пользователя
	c.Header("Content-Type", "application/json")

	userID, err := u.auth.GetUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get cookie",
			"error":   err.Error(),
		})
		return
	}

	orders, err := u.orderUsecase.ListWithAccrual(c.Request.Context(), &orderModel.ListPars{
		UserID:  &userID,
		OrderBy: "ASC",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get orders",
			"error":   err.Error(),
		})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, orders)
}

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

func (u *UserHandlers) WithdrawBalance(c *gin.Context) {
	// Реализация запроса на списание баллов
	req := &WithdrawalBalanceRequest{}
	ctx := c.Request.Context()

	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	userID, err := u.auth.GetUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get cookie",
			"error":   err.Error(),
		})
		return
	}

	req.UserID = userID

	foundUser, _, err := u.userUsecase.Get(ctx,
		&userModel.GetPars{
			ID: userID,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user balance",
			"error":   err.Error(),
		})
		return
	}

	if foundUser.Balance < req.Sum {
		c.JSON(http.StatusPaymentRequired, nil)
		return
	}

	_, orderExist, err := u.orderUsecase.Get(ctx, &orderModel.GetPars{
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

	err = u.bonusTransactionsUsecase.Create(ctx, &bonusTransactionsModel.GetPars{
		OrderNumber:     req.OrderNumber,
		UserID:          userID,
		TransactionType: bonusTransactionsModel.Debit,
		Sum:             req.Sum,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to withdraw balance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) GetWithdrawals(c *gin.Context) {
	userID, err := u.auth.GetUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get cookie",
			"error":   err.Error(),
		})
		return
	}

	withdrawals, err := u.bonusTransactionsUsecase.List(c.Request.Context(), &bonusTransactionsModel.ListPars{
		UserID:          &userID,
		TransactionType: bonusTransactionsModel.Debit,
		OrderBy:         "ASC",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to lost transactions",
			"error":   err.Error(),
		})
		return
	}

	if len(withdrawals) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "not found",
		})
		return
	}

	result := []*WithdrawalsResult{}
	for _, v := range withdrawals {
		withdrawal := WithdrawalsResult{
			OrderNumber: v.OrderNumber,
			Sum:         v.Sum,
			ProcessedAt: v.ProcessedAt.Format(time.RFC3339),
		}

		result = append(result, &withdrawal)
	}

	c.JSON(http.StatusOK, result)
}
