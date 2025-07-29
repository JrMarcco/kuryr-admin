package ioc

import (
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/easy-grpc/client/br"
	"github.com/JrMarcco/easy-grpc/client/rr"
	"github.com/JrMarcco/easy-grpc/registry"
	configv1 "github.com/JrMarcco/kuryr-api/api/config/v1"
	notificationv1 "github.com/JrMarcco/kuryr-api/api/notification/v1"
	providerv1 "github.com/JrMarcco/kuryr-api/api/provider/v1"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/keepalive"
)

var GrpcClientFxOpt = fx.Provide(
	InitBizConfigGrpcClients,
	InitProviderGrpcClients,
	InitNotificationGrpcClients,
)

type keepaliveConfig struct {
	Time                int  `mapstructure:"time"`
	Timeout             int  `mapstructure:"timeout"`
	PermitWithoutStream bool `mapstructure:"permit_without_stream"`
}

type lbConfig struct {
	Name      string           `mapstructure:"name"`
	Timeout   int              `mapstructure:"timeout"`
	KeepAlive *keepaliveConfig `mapstructure:"keep_alive"`
}

func loadLoadBalanceConfig() *lbConfig {
	cfg := &lbConfig{}
	if err := viper.UnmarshalKey("load_balance", &cfg); err != nil {
		panic(err)
	}
	return cfg
}

// InitBizConfigGrpcClients 初始化 biz config grpc client 管理器
func InitBizConfigGrpcClients(r registry.Registry) *client.Manager[configv1.BizConfigServiceClient] {
	cfg := loadLoadBalanceConfig()
	bb := base.NewBalancerBuilder(
		cfg.Name,
		br.NewRwWeightBalancerBuilder(),
		base.Config{HealthCheck: true},
	)

	// 注册负载均衡
	balancer.Register(bb)

	return client.NewManagerBuilder(
		rr.NewResolverBuilder(r, time.Duration(cfg.Timeout)*time.Millisecond),
		bb,
		func(conn *grpc.ClientConn) configv1.BizConfigServiceClient {
			return configv1.NewBizConfigServiceClient(conn)
		},
	).KeepAlive(keepalive.ClientParameters{
		Time:                time.Duration(cfg.KeepAlive.Time) * time.Millisecond,
		Timeout:             time.Duration(cfg.KeepAlive.Timeout) * time.Millisecond,
		PermitWithoutStream: cfg.KeepAlive.PermitWithoutStream,
	}).Insecure().Build()
}

// InitProviderGrpcClients 初始化 provider grpc client 管理器
func InitProviderGrpcClients(r registry.Registry) *client.Manager[providerv1.ProviderServiceClient] {
	cfg := loadLoadBalanceConfig()
	bb := base.NewBalancerBuilder(
		cfg.Name,
		br.NewRwWeightBalancerBuilder(),
		base.Config{HealthCheck: true},
	)

	// 注册负载均衡
	balancer.Register(bb)

	return client.NewManagerBuilder(
		rr.NewResolverBuilder(r, time.Duration(cfg.Timeout)*time.Millisecond),
		bb,
		func(conn *grpc.ClientConn) providerv1.ProviderServiceClient {
			return providerv1.NewProviderServiceClient(conn)
		},
	).KeepAlive(keepalive.ClientParameters{
		Time:                time.Duration(cfg.KeepAlive.Time) * time.Millisecond,
		Timeout:             time.Duration(cfg.KeepAlive.Timeout) * time.Millisecond,
		PermitWithoutStream: cfg.KeepAlive.PermitWithoutStream,
	}).Insecure().Build()
}

// InitNotificationGrpcClients 初始化 notification grpc client 管理器
func InitNotificationGrpcClients(r registry.Registry) *client.Manager[notificationv1.NotificationServiceClient] {
	cfg := loadLoadBalanceConfig()
	bb := base.NewBalancerBuilder(
		cfg.Name,
		br.NewRwWeightBalancerBuilder(),
		base.Config{HealthCheck: true},
	)

	// 注册负载均衡
	balancer.Register(bb)

	return client.NewManagerBuilder(
		rr.NewResolverBuilder(r, time.Duration(cfg.Timeout)*time.Millisecond),
		bb,
		func(conn *grpc.ClientConn) notificationv1.NotificationServiceClient {
			return notificationv1.NewNotificationServiceClient(conn)
		},
	).KeepAlive(keepalive.ClientParameters{
		Time:                time.Duration(cfg.KeepAlive.Time) * time.Millisecond,
		Timeout:             time.Duration(cfg.KeepAlive.Timeout) * time.Millisecond,
		PermitWithoutStream: cfg.KeepAlive.PermitWithoutStream,
	}).Insecure().Build()
}
