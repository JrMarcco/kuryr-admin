package ioc

import (
	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret/base64"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	configv1 "github.com/JrMarcco/kuryr-api/api/config/v1"
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
		grpcServerNameBizConfig(), db, bizRepo, userRepo, generator, grpcClients,
	)
}

func InitBizConfigService(
	grpcClients *client.Manager[configv1.BizConfigServiceClient],
) *service.DefaultBizConfigService {
	return service.NewDefaultBizConfigService(
		grpcServerNameBizConfig(), grpcClients,
	)
}
