package ioc

import (
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/easy-grpc/client/bl"
	"github.com/JrMarcco/easy-grpc/client/sl"
	"github.com/JrMarcco/easy-grpc/registry"
	configv1 "github.com/JrMarcco/kuryr-api/api/config/v1"
	notificationv1 "github.com/JrMarcco/kuryr-api/api/notification/v1"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

var GrpcClientFxOpt = fx.Provide(
	InitNotificationGrpcClients,
	InitBizConfigGrpcClients,
)

type lbConfig struct {
	Name    string `mapstructure:"name"`
	Timeout int    `mapstructure:"timeout"`
}

func loadLoadBalanceConfig() *lbConfig {
	cfg := &lbConfig{}
	if err := viper.UnmarshalKey("load_balance", &cfg); err != nil {
		panic(err)
	}
	return cfg
}

// InitNotificationGrpcClients 初始化 notification grpc client 管理器
func InitNotificationGrpcClients(r registry.Registry) *client.Manager[notificationv1.NotificationServiceClient] {
	cfg := loadLoadBalanceConfig()
	bb := base.NewBalancerBuilder(
		cfg.Name,
		bl.NewRwWeightBalancerBuilder(),
		base.Config{HealthCheck: true},
	)

	// 注册负载均衡
	balancer.Register(bb)

	return client.NewManagerBuilder(
		sl.NewResolverBuilder(r, time.Duration(cfg.Timeout)*time.Millisecond),
		bb,
		func(conn *grpc.ClientConn) notificationv1.NotificationServiceClient {
			return notificationv1.NewNotificationServiceClient(conn)
		},
	).Insecure().Build()
}

// InitBizConfigGrpcClients 初始化 biz config grpc client 管理器
func InitBizConfigGrpcClients(r registry.Registry) *client.Manager[configv1.BizConfigServiceClient] {
	cfg := loadLoadBalanceConfig()
	bb := base.NewBalancerBuilder(
		cfg.Name,
		bl.NewRwWeightBalancerBuilder(),
		base.Config{HealthCheck: true},
	)

	// 注册负载均衡
	balancer.Register(bb)

	return client.NewManagerBuilder(
		sl.NewResolverBuilder(r, time.Duration(cfg.Timeout)*time.Millisecond),
		bb,
		func(conn *grpc.ClientConn) configv1.BizConfigServiceClient {
			return configv1.NewBizConfigServiceClient(conn)
		},
	).Insecure().Build()
}
