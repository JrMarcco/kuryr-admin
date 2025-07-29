package ioc

import (
	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret/base64"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	configv1 "github.com/JrMarcco/kuryr-api/api/config/v1"
	providerv1 "github.com/JrMarcco/kuryr-api/api/provider/v1"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var ServiceFxOpt = fx.Options(
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
			InitBizService,
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

func grpcServerNameBizConfig() string {
	type config struct {
		Name string `mapstructure:"name"`
	}

	cfg := config{}
	if err := viper.UnmarshalKey("grpc_servers.biz_config", &cfg); err != nil {
		panic(err)
	}
	return cfg.Name
}

func InitBizService(
	db *gorm.DB, bizRepo repository.BizRepo, userRepo repository.UserRepo, generator secret.Generator,
	grpcClients *client.Manager[configv1.BizConfigServiceClient],
) *service.DefaultBizService {
	return service.NewDefaultBizService(
		grpcServerNameBizConfig(), grpcClients, db, bizRepo, userRepo, generator,
	)
}

func InitBizConfigService(
	grpcClients *client.Manager[configv1.BizConfigServiceClient], bizRepo repository.BizRepo,
) *service.DefaultBizConfigService {
	return service.NewDefaultBizConfigService(
		grpcServerNameBizConfig(), grpcClients, bizRepo,
	)
}

func InitProviderService(
	grpcClients *client.Manager[providerv1.ProviderServiceClient],
) *service.DefaultProviderService {
	return service.NewDefaultProviderService(
		grpcServerNameBizConfig(), grpcClients,
	)
}
