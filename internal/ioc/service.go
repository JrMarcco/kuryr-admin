package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"go.uber.org/fx"
)

var ServiceFxOpt = fx.Options(
	fx.Provide(
		// user service
		fx.Annotate(
			service.NewUserService,
			fx.As(new(service.UserService)),
			fx.ParamTags(``, `name:"access-token-manager"`, `name:"refresh-token-manager"`),
		),
		// biz service
		fx.Annotate(
			service.NewBizService,
			fx.As(new(service.BizService)),
		),
	),
)
