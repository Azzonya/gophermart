package app

import (
	"context"
	"fmt"
	"github.com/Azzonya/gophermart/internal/client/accrual"
	"github.com/Azzonya/gophermart/internal/config"
	OrderService "github.com/Azzonya/gophermart/internal/domain/order/service"
	UserService "github.com/Azzonya/gophermart/internal/domain/user/service"
	WithdrawalService "github.com/Azzonya/gophermart/internal/domain/withdrawal/service"
	"github.com/Azzonya/gophermart/internal/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

type App struct {
	pgpool *pgxpool.Pool

	// clients
	accrual *accrual.Client

	// user
	userService *UserService.Service

	// order
	orderService *OrderService.Service

	// withdrawal
	withdrawalService *WithdrawalService.Service

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

	// http-gw server
	{
		//userHandlers
		a.userHandlers = handler.New()

		// server
		a.rest = NewRest(a.userHandlers)
	}
}

func (a *App) PreStartHook() {
	slog.Info("PreStartHook")
}

func (a *App) Start() {
	slog.Info("Starting")

	// services
	{

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
