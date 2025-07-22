package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret/base64"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"go.uber.org/fx"
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
			service.NewDefaultBizService,
			fx.As(new(service.BizService)),
		),

		// biz config service
		fx.Annotate(
			service.NewDefaultBizConfigService,
			fx.As(new(service.BizConfigService)),
		),
	),
)
