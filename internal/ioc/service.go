package ioc

import (
	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret/base64"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	businessv1 "github.com/JrMarcco/kuryr-api/api/go/business/v1"
	configv1 "github.com/JrMarcco/kuryr-api/api/go/config/v1"
	providerv1 "github.com/JrMarcco/kuryr-api/api/go/provider/v1"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ServiceFxOpt = fx.Module(
	"service",
	fx.Provide(
		// user service
		fx.Annotate(
			base64.NewGenerator,
			fx.As(new(secret.Generator)),
		),
		fx.Annotate(
			service.NewJwtUserService,
			fx.As(new(service.UserService)),
			fx.ParamTags(``, `name:"access-token-manager"`, `name:"refresh-token-manager"`),
		),

		// biz service
		fx.Annotate(
			InitBizInfoService,
			fx.As(new(service.BizService)),
		),

		// biz config service
		fx.Annotate(
			InitBizConfigService,
			fx.As(new(service.BizConfigService)),
		),

		// provider service
		fx.Annotate(
			InitProviderService,
			fx.As(new(service.ProviderService)),
		),
	),
)

func grpcServerName() string {
	var name string
	if err := viper.UnmarshalKey("grpc.server.name", &name); err != nil {
		panic(err)
	}
	return name
}

func InitBizInfoService(grpcClients *client.Manager[businessv1.BusinessServiceClient], userRepo repository.UserRepo) *service.DefaultBizService {
	return service.NewDefaultBizService(
		grpcServerName(), grpcClients, userRepo,
	)
}

func InitBizConfigService(grpcClients *client.Manager[configv1.BizConfigServiceClient]) *service.DefaultBizConfigService {
	return service.NewDefaultBizConfigService(
		grpcServerName(), grpcClients,
	)
}

func InitProviderService(grpcClients *client.Manager[providerv1.ProviderServiceClient]) *service.DefaultProviderService {
	return service.NewDefaultProviderService(
		grpcServerName(), grpcClients,
	)
}
