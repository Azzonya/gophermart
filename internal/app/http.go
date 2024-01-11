package app

import (
	"context"
	"fmt"
	"github.com/Azzonya/gophermart/internal/handler"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type Rest struct {
	server       *http.Server
	userHandlers *handler.UserHandlers

	ErrorChan chan error
}

func NewRest(userHandlers *handler.UserHandlers) *Rest {
	return &Rest{
		userHandlers: userHandlers,

		ErrorChan: make(chan error, 1),
	}
}

func (o *Rest) Start(lAddr string) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	o.SetRouters(r)

	o.server = &http.Server{
		Addr:              lAddr,
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       time.Minute,
		MaxHeaderBytes:    300 * 1024,
	}

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Println("Panic")
			}
		}()

		err := o.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			o.ErrorChan <- err
		}
	}()
}

func (o *Rest) Stop(ctx context.Context) error {
	defer close(o.ErrorChan)

	err := o.server.Shutdown(ctx)
	if err != nil {
		slog.Error("http-server shutdown error", "error", err)
		return err
	}

	return nil
}

func (o *Rest) SetRouters(r *gin.Engine) {
	apiGroup := r.Group("/api")
	userGroup := apiGroup.Group("/user")

	userGroup.POST("/register", o.userHandlers.RegisterUser)
	userGroup.POST("/login", o.userHandlers.LoginUser)
	userGroup.POST("/orders", o.userHandlers.UploadOrder)
	userGroup.GET("/orders", o.userHandlers.GetOrders)
	userGroup.GET("/balance", o.userHandlers.GetBalance)
	userGroup.POST("/balance/withdraw", o.userHandlers.WithdrawBalance)
	userGroup.POST("/withdrawals", o.userHandlers.GetWithdrawals)
}
