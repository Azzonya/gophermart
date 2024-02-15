package handler

import (
	"errors"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/Azzonya/gophermart/internal/errs"
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
	userID, _ := u.userUsecase.GetUserIDFromCookieU(c)

	err = u.orderUsecase.Create(c.Request.Context(), &entities.Order{
		OrderNumber: orderNumber,
		Status:      entities.OrderStatusNew,
		UserID:      userID,
	})

	switch {
	case errors.As(err, &errs.ErrOrderNumberLuhnValid{}):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Неверный формат номера заказа (алгоритм Луна)"})
	case errors.As(err, &errs.ErrOrderUploaded{}):
		c.JSON(http.StatusOK, gin.H{"message": "already uploaded"})
	case errors.As(err, &errs.ErrOrderUploadedByAnotherUser{}):
		c.JSON(http.StatusConflict, gin.H{"message": "already uploaded by different user"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order", "error": err.Error()})
	default:
		c.JSON(http.StatusAccepted, nil)
	}
}

func (u *UserHandlers) GetOrders(c *gin.Context) {
	// Реализация получения списка заказов пользователя
	userID, _ := u.userUsecase.GetUserIDFromCookieU(c)

	orders, err := u.orderUsecase.ListWithAccrual(c.Request.Context(), &entities.OrderListPars{
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
