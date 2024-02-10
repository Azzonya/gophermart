package handler

import (
	"errors"
	userModel "github.com/Azzonya/gophermart/internal/domain/user"
	"github.com/Azzonya/gophermart/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Jsn struct {
	Host         string `json:"host"`
	ShortMessage string `json:"short_message"`
	Err          string `json:"_err"`
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

	newUser, err := u.userUsecase.Register(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotUniq{}):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
		}
		return
	}

	sessionCookie, errS := u.auth.CreateJWTCookie(newUser)
	if errS != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT cookie", "details": err.Error()})
		return
	}

	c.Header("Set-Cookie", sessionCookie.String())
	c.JSON(http.StatusOK, nil)
}

func (u *UserHandlers) LoginUser(c *gin.Context) {
	// Реализация аутентификации пользователя
	c.Header("Content-Type", "application/json")

	var err error
	req := &userModel.GetPars{}

	userID, _ := u.auth.GetUserIDFromCookie(c)
	if len(userID) != 0 {
		c.JSON(http.StatusOK, nil)
		return
	}

	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body", "details": err.Error()})
		return
	}

	foundUser, exist, err := u.userUsecase.CheckAuth(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check authentication", "details": err.Error()})
		return
	}
	if !exist {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	sessionCookie, errS := u.auth.CreateJWTCookie(foundUser)
	if errS != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.Header("Set-Cookie", sessionCookie.String())
	c.JSON(http.StatusOK, nil)
}
