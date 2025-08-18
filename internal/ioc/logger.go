package ioc

import (
	"context"
	"log/slog"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

var LoggerFxOpt = fx.Module("logger", fx.Provide(InitLogger))

func InitLogger(lc fx.Lifecycle) *zap.Logger {
	type config struct {
		Env string `mapstructure:"env"`
	}

	cfg := config{}
	if err := viper.UnmarshalKey("profile", &cfg); err != nil {
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

	// 初始化 slog
	slog.SetDefault(slog.New(zapslog.NewHandler(zapLogger.Core())))

	// 注册生命周期 hook
	lc.Append(fx.Hook{
		// 程序停止时 flush buffer 防止日志丢失
		OnStop: func(ctx context.Context) error {
			_ = zapLogger.Sync()
			return nil
		},
	})

	return zapLogger
}
