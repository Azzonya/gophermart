package handler

import (
	bonusTransactionsModel "github.com/Azzonya/gophermart/internal/domain/bonustransactions"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order"
	userModel "github.com/Azzonya/gophermart/internal/domain/user"
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

	foundUser, _, err := u.userUsecase.Get(c.Request.Context(), &userModel.GetPars{ID: req.UserID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user balance", "details": err.Error()})
		return
	}

	if foundUser.Balance < req.Sum {
		c.JSON(http.StatusPaymentRequired, gin.H{"message": "Insufficient balance"})
		return
	}

	_, orderExist, err := u.orderUsecase.Get(ctx, &orderModel.GetPars{
		OrderNumber: req.OrderNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check order existence", "details": err.Error()})
		return
	}
	if !orderExist {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Order not found"})
		return
	}

	err = u.bonusTransactionsUsecase.Create(ctx, &bonusTransactionsModel.GetPars{
		OrderNumber:     req.OrderNumber,
		UserID:          req.UserID,
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
