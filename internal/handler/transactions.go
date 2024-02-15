package handler

import (
	"errors"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/Azzonya/gophermart/internal/errs"
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

	err = u.bonusTransactionsUsecase.WithdrawBalanceU(ctx, &entities.BonusTransaction{
		OrderNumber:     req.OrderNumber,
		UserID:          req.UserID,
		TransactionType: entities.Debit,
		Sum:             req.Sum,
	})

	switch {
	case errors.As(err, &errs.ErrUserInsufficientBalance{}):
		c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to withdraw balance", "error": err.Error()})
	default:
		c.JSON(http.StatusOK, nil)
	}
}

func (u *UserHandlers) GetWithdrawals(c *gin.Context) {
	userID, _ := u.auth.GetUserIDFromCookie(c)

	withdrawals, err := u.bonusTransactionsUsecase.List(c.Request.Context(), &entities.BonusTransactionsListPars{
		UserID:          &userID,
		TransactionType: entities.Debit,
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
