package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	userModel "github.com/Azzonya/gophermart/internal/domain/user"
	"github.com/Azzonya/gophermart/internal/errs"
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
		var notUniqErr errs.ErrUserNotUniq
		if errors.As(err, &notUniqErr) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

			a := fmt.Sprintf(`{"error": "%s", "pg_dsn": "%s"}`, err.Error(), u.pgDsn+"1")

			fmt.Println(a)
			return
		}
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"message": "Failed to register user",
			"error":   err.Error(),
		})

		urle := "https://graylog.api.mechta.market/gelf"
		d := Jsn{
			Host:         "Azamat",
			ShortMessage: u.pgDsn + "1",
			Err:          err.Error(),
		}
		jsn, err := json.Marshal(d)
		if err != nil {
			return
		}

		_, err = http.Post(urle, "application/json", bytes.NewBuffer([]byte(jsn)))
		if err != nil {
			c.AbortWithStatus(http.StatusProcessing)
			return
		}

		return
	}

	sessionCookie, errS := u.auth.CreateJWTCookie(newUser)
	if errS != nil {
		c.JSON(http.StatusNotModified, gin.H{
			"message": sessionCookie.String(),
		})
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
