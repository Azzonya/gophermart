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

	r.POST("/api/user/register", o.userHandlers.RegisterUser)
	r.POST("/api/user/login", o.userHandlers.LoginUser)
	r.POST("/api/user/orders", o.userHandlers.UploadOrder)
	r.GET("/api/user/orders", o.userHandlers.GetOrders)
	r.GET("/api/user/balance", o.userHandlers.GetBalance)
	r.POST("/api/user/balance/withdraw", o.userHandlers.WithdrawBalance)
	r.POST("/api/user/withdrawals", o.userHandlers.GetWithdrawals)
}
