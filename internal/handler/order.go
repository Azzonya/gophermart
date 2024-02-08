package handler

import (
	"context"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (u *UserHandlers) UploadOrder(c *gin.Context) {
	// Реализация загрузки номера заказа
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	ctx := c.Request.Context()

	orderNumber := string(body)

	if !u.orderUsecase.IsLuhnValid(orderNumber) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Неверный формат номера заказа (алгоритм Луна)"})
		return
	}

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

	c.JSON(http.StatusAccepted, nil)
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
