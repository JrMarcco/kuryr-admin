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
	fx.Annotate(
		InitApp,
		// 第一个参数是 zap.Logger 不需要使用标签
		fx.ParamTags(``, ``, `group:"middleware-builder"`),
	),
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

	svr := &http.Server{
		Addr:    cfg.Addr,
		Handler: engine.Handler(),
	}

	// 注册中间件
	middlewares := make([]gin.HandlerFunc, 0, len(mbs))
	for _, mb := range mbs {
		middlewares = append(middlewares, mb.Build())
	}
	engine.Use(middlewares...)

	return &App{
		svr:    svr,
		logger: logger,
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
