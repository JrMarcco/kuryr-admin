package main

import (
	"github.com/JrMarcco/kuryr-admin/internal/ioc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	initViper()

	fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		// 初始化 zap.Logger
		ioc.LoggerFxOpt,
		// 初始化 redis.Client
		ioc.RedisFxOpt,
		// 初始化 gorm.DB
		ioc.DBFxOpt,
		// 初始化 etcd
		ioc.EtcdFxOpt,
		// 初始化 grpc registry
		ioc.RegistryFxOpt,
		// 初始化 grpc client manager
		ioc.GrpcClientFxOpt,
		// 初始化 jwt manager
		ioc.JwtManagerOpt,
		// 初始化 middleware builder
		ioc.MiddlewareBuilderOpt,
		// 初始化 repo
		ioc.RepoFxOpt,
		// 初始化 service
		ioc.ServiceFxOpt,
		// 初始化 handler
		ioc.HandlerFxOpt,

		// 注册 gin 路由，需要在 app 启动前完成
		// ioc.HandlerFxInvoke,

		// 初始化 ioc.App
		ioc.AppFxOpt,
	).Run()
}

// initViper 初始化 viper
func initViper() {
	configFile := pflag.String("config", "etc/config.yaml", "配置文件路径")
	pflag.Parse()

	viper.SetConfigFile(*configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
