package app

import (
	"context"
	"fmt"
	"github.com/Azzonya/gophermart/internal/auth"
	"github.com/Azzonya/gophermart/internal/client/accrual/http-gw"
	"github.com/Azzonya/gophermart/internal/config"
	bonusService "github.com/Azzonya/gophermart/internal/domain/bonus/service"
	bonusTransactionsRepo "github.com/Azzonya/gophermart/internal/domain/bonus_transactions/repo/db"
	bonusTransactionsService "github.com/Azzonya/gophermart/internal/domain/bonus_transactions/service"
	orderRepo "github.com/Azzonya/gophermart/internal/domain/order/repo/db"
	OrderService "github.com/Azzonya/gophermart/internal/domain/order/service"
	userRepo "github.com/Azzonya/gophermart/internal/domain/user/repo/db"
	UserService "github.com/Azzonya/gophermart/internal/domain/user/service"
	"github.com/Azzonya/gophermart/internal/handler"
	"github.com/Azzonya/gophermart/internal/usecase/bonus_transactions"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"github.com/Azzonya/gophermart/internal/usecase/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

type App struct {
	pgpool *pgxpool.Pool

	// clients
	accrualClient *http_gw.Client

	// bonus
	bonusService *bonusService.Service

	// user
	userService *UserService.Service

	// order
	orderService *OrderService.Service

	// bonus_transactions
	bonusTransactionsService *bonusTransactionsService.Service

	// handlers
	userHandlers *handler.UserHandlers

	// http-gw
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
		a.accrualClient = http_gw.New(config.Conf.AccrualSystemAddress)
	}

	// http-gw-gw server
	{
		authorizer := auth.New(config.Conf.JwtSecret)

		//bonus transaction
		bonusTransactionsRepoV := bonusTransactionsRepo.New(a.pgpool)
		a.bonusTransactionsService = bonusTransactionsService.New(bonusTransactionsRepoV)
		bonusTransactionsUsecase := bonus_transactions.New(a.bonusTransactionsService)

		//user
		userRepoV := userRepo.New(a.pgpool)
		a.userService = UserService.New(userRepoV, a.bonusTransactionsService)
		userUsecase := user.New(a.userService)

		//order
		orderRepoV := orderRepo.New(a.pgpool)
		a.orderService = OrderService.New(orderRepoV, a.bonusTransactionsService)
		orderUsecase := order.New(a.orderService)

		//bonus
		a.bonusService = bonusService.New(a.accrualClient, a.bonusTransactionsService, a.orderService)

		//handers
		a.userHandlers = handler.New(authorizer, userUsecase, orderUsecase, bonusTransactionsUsecase)

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

	// http-gw server
	{
		a.rest.Start(config.Conf.RunAddress)
		slog.Info("http-gw-server started " + config.Conf.RunAddress)
	}
}

func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
	a.bonusService.Wg.Wait()
}

func (a *App) Stop() {
	slog.Info("Shutting down...")

	// http-gw server
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
