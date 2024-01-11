package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

var Conf = struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	PgDsn                string `env:"PG_DSN"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}{}

func init() {
	flag.StringVar(&Conf.RunAddress, "a", "localhost:8080", "address and port where server start")
	flag.StringVar(&Conf.PgDsn, "d", "", "database connection line")
	flag.StringVar(&Conf.AccrualSystemAddress, "r", "", "accrual system address")

	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
