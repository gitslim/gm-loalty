package main

import (
	"github.com/gitslim/gophermart/internal/conf"
	"github.com/gitslim/gophermart/internal/log"
	"github.com/gitslim/gophermart/internal/web"
	"go.uber.org/fx"
)

func main() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			conf.ParseConfig,
			log.NewLogger,
			web.NewRouter),
		fx.Invoke(
			web.RegisterServerHooks,
		),
	)
}
