package ioc

import (
	"context"
	"errors"
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/pkg/gin/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppFxOpt = fx.Provide(
	gin.Default,
	InitMiddlewares,
	InitApp,
)

var AppFxInvoke = fx.Invoke(AppLifecycle)

type App struct {
	svr    *http.Server
	logger *zap.Logger
}

func (app *App) Start() error {
	go func() {
		if err := app.svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Fatal("[kuryr-admin] listen and serve", zap.Error(err))
		}
	}()
	return nil
}

func (app *App) Stop(ctx context.Context) error {
	app.logger.Info("[kuryr-admin] shutdown server ...")
	if err := app.svr.Shutdown(ctx); err != nil {
		app.logger.Error("[kuryr-admin] server shutdown", zap.Error(err))
		return err
	}
	app.logger.Info("[kuryr-admin] server exited")
	return nil
}

func InitApp(engine *gin.Engine, logger *zap.Logger, mbs []middleware.Builder) *App {
	type config struct {
		Addr string `mapstructure:"addr"`
	}

	cfg := config{}
	if err := viper.UnmarshalKey("app", &cfg); err != nil {
	}

	middlewares := make([]gin.HandlerFunc, 0, len(mbs))
	for _, mb := range mbs {
		middlewares = append(middlewares, mb.Build())
	}
	engine.Use(middlewares...)

	svr := &http.Server{
		Addr:    cfg.Addr,
		Handler: engine.Handler(),
	}

	return &App{
		svr:    svr,
		logger: logger,
	}
}

// InitMiddlewares 提供一个用于创建有序中间件切片的函数
func InitMiddlewares(
	corsBuilder *middleware.CorsBuilder,
	jwtBuilder *middleware.JwtBuilder,
) []middleware.Builder {
	// 按顺序排列中间件
	return []middleware.Builder{
		corsBuilder,
		jwtBuilder,
	}
}

func AppLifecycle(app *App, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return app.Start()
		},
		OnStop: func(ctx context.Context) error {
			return app.Stop(ctx)
		},
	})
}
