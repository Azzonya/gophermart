package handler

import (
	"errors"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order"
	"github.com/Azzonya/gophermart/internal/storage"
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

	orderNumber := string(body)
	userID, _ := u.auth.GetUserIDFromCookie(c)

	err = u.orderUsecase.Create(c.Request.Context(), &orderModel.GetPars{
		OrderNumber: orderNumber,
		Status:      orderModel.OrderStatusNew,
		UserID:      userID,
	})

	switch {
	case errors.Is(err, storage.ErrOrderNumberLuhnValid{}):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Неверный формат номера заказа (алгоритм Луна)"})
	case errors.Is(err, storage.ErrOrderUploaded{}):
		c.JSON(http.StatusOK, gin.H{"message": "already uploaded"})
	case errors.Is(err, storage.ErrOrderUploadedByAnotherUser{}):
		c.JSON(http.StatusConflict, gin.H{"message": "already uploaded by different user"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order", "error": err.Error()})
	default:
		c.JSON(http.StatusAccepted, nil)
	}
}

func (u *UserHandlers) GetOrders(c *gin.Context) {
	// Реализация получения списка заказов пользователя
	userID, _ := u.auth.GetUserIDFromCookie(c)

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
