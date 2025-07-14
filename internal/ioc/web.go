package ioc

import (
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var HandlerFxOpt = fx.Provide(
	fx.Annotate(
		web.NewUserHandler,
		fx.As(new(ginpkg.Registry)),
		fx.ResultTags(`group:"handler"`),
	),
)

var HandlerFxInvoke = fx.Invoke(
	fx.Annotate(
		RegisterRoutes,
		fx.ParamTags(``, `group:"handler"`),
	),
)

func RegisterRoutes(engine *gin.Engine, registries []ginpkg.Registry) {
	for _, registry := range registries {
		registry.RegisterRoutes(engine)
	}
}
