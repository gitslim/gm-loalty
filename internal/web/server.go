package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/gophermart/internal/conf"
	"github.com/gitslim/gophermart/internal/log"
	"github.com/gitslim/gophermart/internal/middleware"
	"go.uber.org/fx"
)

// NewRouter создает новый экземпляр Gin Engine
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.GzipMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}

// RegisterServerHooks регистрирует lifecycle hooks для запуска и остановки сервера
func RegisterServerHooks(lc fx.Lifecycle, cfg *conf.Config, log *log.Logger, router *gin.Engine) {
	srv := &http.Server{
		Addr:    cfg.RunAddrress,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Infof("Starting HTTP server on %v", srv.Addr)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("HTTP server failed: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})
}
