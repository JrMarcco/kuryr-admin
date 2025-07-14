package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"go.uber.org/fx"
)

var ServiceFxOpt = fx.Options(
	fx.Provide(
		fx.Annotate(
			service.NewUserService,
			fx.As(new(service.UserService)),
		)),
)
