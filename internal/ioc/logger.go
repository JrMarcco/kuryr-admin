package ioc

import (
	"context"
	"log/slog"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

var (
	LoggerFxOpt    = fx.Provide(InitLogger)
	LoggerFxInvoke = fx.Invoke(LoggerLifecycle)
	SlogFxInvoke   = fx.Invoke(InitSlog)
)

func InitLogger() *zap.Logger {
	type config struct {
		Env string `mapstructure:"env"`
	}

	cfg := config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	var zapLogger *zap.Logger
	var err error
	switch cfg.Env {
	case "prod":
		zapLogger, err = zap.NewProduction()
	default:
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	return zapLogger
}

// LoggerLifecycle 注册 zap.Logger 生命周期
// 在程序退出时 flush buffer 防止日志丢失
func LoggerLifecycle(lc fx.Lifecycle, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = logger.Sync()
			return nil
		},
	})
}

// InitSlog 设置全局 slog 默认使用 zap.Logger 实例。
func InitSlog(logger *zap.Logger) {
	slog.SetDefault(slog.New(zapslog.NewHandler(logger.Core())))
}
