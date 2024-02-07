package app

import (
	"context"
	"fmt"
	"github.com/Azzonya/gophermart/internal/client/accrual/http"
	"github.com/Azzonya/gophermart/internal/config"
	"github.com/Azzonya/gophermart/internal/domain/auth"
	bonusService "github.com/Azzonya/gophermart/internal/domain/bonus"
	bonusTransactionsService "github.com/Azzonya/gophermart/internal/domain/bonustransactions"
	orderService "github.com/Azzonya/gophermart/internal/domain/order"
	userService "github.com/Azzonya/gophermart/internal/domain/user"
	"github.com/Azzonya/gophermart/internal/handler"
	bonusTransactionsUsecase "github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	orderUsecase "github.com/Azzonya/gophermart/internal/usecase/order"
	userUsecase "github.com/Azzonya/gophermart/internal/usecase/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

type App struct {
	pgpool *pgxpool.Pool

	// clients
	accrualClient *http.Client

	// bonus
	bonusService *bonusService.Service

	// user
	userService *userService.Service

	// order
	orderService *orderService.Service

	// bonustransactions
	bonusTransactionsService *bonusTransactionsService.Service

	// handlers
	userHandlers *handler.UserHandlers

	// http
	rest *Rest

	exitCode int
}

func (a *App) Init() {
	var err error

	// pgpool
	{
		a.pgpool, err = pgxpool.New(context.Background(), config.Conf.PgDsn)
		errCheck(err, "pgxpool.New")
	}

	// accrual client
	{
		a.accrualClient = http.New(config.Conf.AccrualSystemAddress)
	}

	// http-gw server
	{
		authorizer := auth.New(config.Conf.JwtSecret)

		//bonus transaction
		bonusTransactionsRepoV := bonusTransactionsService.NewRepo(a.pgpool)
		a.bonusTransactionsService = bonusTransactionsService.New(bonusTransactionsRepoV)
		bonusTransactionsUsecaseV := bonusTransactionsUsecase.New(a.bonusTransactionsService)

		//user
		userRepoV := userService.NewRepo(a.pgpool)
		a.userService = userService.New(userRepoV, a.bonusTransactionsService)
		userUsecaseV := userUsecase.New(a.userService)

		//order
		orderRepoV := orderService.NewRepo(a.pgpool)
		a.orderService = orderService.New(orderRepoV, a.bonusTransactionsService)
		orderUsecaseV := orderUsecase.New(a.orderService)

		//bonus
		a.bonusService = bonusService.New(a.accrualClient, a.bonusTransactionsService, a.orderService)

		//handers
		a.userHandlers = handler.New(authorizer, userUsecaseV, orderUsecaseV, bonusTransactionsUsecaseV)

		// server
		a.rest = NewRest(a.userHandlers, config.Conf.JwtSecret)
	}

}

func (a *App) PreStartHook() {
	slog.Info("PreStartHook")
}

func (a *App) Start() {
	slog.Info("Starting")

	// services
	{
		a.bonusService.Start(10 * time.Second)
	}

	// http server
	{
		a.rest.Start(config.Conf.RunAddress)
		slog.Info("http-server started " + config.Conf.RunAddress)
	}
}

func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
	a.bonusService.Wait()
}

func (a *App) Stop() {
	slog.Info("Shutting down...")

	// http server
	{
		ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer ctxCancel()

		if err := a.rest.Stop(ctx); err != nil {
			a.exitCode = 1
		}
	}
}

func (a *App) WaitJobs() {
	slog.Info("waiting jobs")
}

func (a *App) Exit() {
	slog.Info("Exit")

	os.Exit(a.exitCode)
}

func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
