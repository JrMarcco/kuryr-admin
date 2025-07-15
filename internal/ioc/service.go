package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/JrMarcco/kuryr-admin/internal/service/session"
	"go.uber.org/fx"
)

var ServiceFxOpt = fx.Options(
	fx.Provide(
		// session service
		fx.Annotate(
			session.NewRedisSessionService,
			fx.As(new(session.Service)),
		),
		// user service
		fx.Annotate(
			service.NewJwtUserService,
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
