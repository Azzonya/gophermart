package handler

import (
	bonusTransactionsModel "github.com/Azzonya/gophermart/internal/domain/bonustransactions"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order"
	userModel "github.com/Azzonya/gophermart/internal/domain/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	//
	//if len(withdrawals) == 0 {
	//	c.JSON(http.StatusNoContent, gin.H{
	//		"message": "not found",
	//	})
	//	return
	//}

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
