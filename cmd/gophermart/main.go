package main

import (
	"github.com/gitslim/gophermart/internal/conf"
	"github.com/gitslim/gophermart/internal/logging"
	"github.com/gitslim/gophermart/internal/logging/sugared"
	"github.com/gitslim/gophermart/internal/service"
	"github.com/gitslim/gophermart/internal/service/balance"
	"github.com/gitslim/gophermart/internal/service/order"
	"github.com/gitslim/gophermart/internal/service/user"
	"github.com/gitslim/gophermart/internal/storage"
	"github.com/gitslim/gophermart/internal/storage/postgres"
	"github.com/gitslim/gophermart/internal/web"
	"github.com/gitslim/gophermart/internal/web/handlers"
	"github.com/gitslim/gophermart/internal/web/middleware"
	"github.com/gitslim/gophermart/internal/web/router"
	"go.uber.org/fx"
)

func main() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		// Конфигурация и логирование
		fx.Provide(
			conf.ParseConfig,
			fx.Annotate(sugared.NewLogger, fx.As(new(logging.Logger))),
		),

		// Хранилище
		fx.Provide(
			postgres.NewConnPool,
			fx.Annotate(postgres.NewPgStorage, fx.As(new(storage.Storage))),
		),

		// Сервисы
		fx.Provide(
			fx.Annotate(user.NewUserService, fx.As(new(service.UserService))),
			fx.Annotate(order.NewOrderService, fx.As(new(service.OrderService))),
			fx.Annotate(balance.NewBalanceService, fx.As(new(service.BalanceService))),
		),

		// Веб-компоненты
		fx.Provide(
			middleware.NewAuthMiddleware,
			handlers.NewHandler,
			router.NewRouter,
		),

		// Запуск сервера
		fx.Invoke(
			postgres.RegisterPoolHooks,
			web.RegisterServerHooks,
		),
	)
}
