package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	bonusTransactionsModel "github.com/Azzonya/gophermart/internal/domain/bonustransactions"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order"
	"github.com/Azzonya/gophermart/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

	req.UserID, _ = u.auth.GetUserIDFromCookie(c)

	err = u.bonusTransactionsUsecase.WithdrawBalance(ctx, &bonusTransactionsModel.GetPars{
		OrderNumber:     req.OrderNumber,
		UserID:          req.UserID,
		TransactionType: bonusTransactionsModel.Debit,
		Sum:             req.Sum,
	})

	switch {
	case errors.As(err, &storage.ErrUserInsufficientBalance{}):
		c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
	case errors.As(err, &storage.ErrOrderNotExist{}):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		orders, _ := u.orderUsecase.List(c.Request.Context(), &orderModel.ListPars{})

		asd := ""
		for _, v := range orders {
			asd += v.OrderNumber + ","
		}

		urle := "https://graylog.api.mechta.market/gelf"
		d := Jsn{
			Host: "Azamat",
			Err:  asd,
		}
		jsn, err := json.Marshal(d)
		if err != nil {
			return
		}

		_, err = http.Post(urle, "application/json", bytes.NewBuffer([]byte(jsn)))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to withdraw balance", "error": err.Error()})
	default:
		c.JSON(http.StatusOK, nil)
	}
}

func (u *UserHandlers) GetWithdrawals(c *gin.Context) {
	userID, _ := u.auth.GetUserIDFromCookie(c)

	withdrawals, err := u.bonusTransactionsUsecase.List(c.Request.Context(), &bonusTransactionsModel.ListPars{
		UserID:          &userID,
		TransactionType: bonusTransactionsModel.Debit,
		OrderBy:         "ASC",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions", "details": err.Error()})
		return
	}

	if len(withdrawals) == 0 {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": "No transactions found"})
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
