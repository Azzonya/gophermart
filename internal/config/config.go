package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

var Conf = struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	PgDsn                string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JwtSecret            string `env:"JWT_SECRET"`
}{}

func init() {
	flag.StringVar(&Conf.RunAddress, "a", "localhost:8095", "address and port where server start")
	flag.StringVar(&Conf.PgDsn, "d", "", "database connection line")
	flag.StringVar(&Conf.AccrualSystemAddress, "r", "http://localhost:8080/", "accrual system address")

	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
