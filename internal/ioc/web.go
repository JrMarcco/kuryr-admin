package ioc

import (
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var HandlerFxOpt = fx.Provide(
	// user handler
	fx.Annotate(
		web.NewUserHandler,
		fx.As(new(ginpkg.RouteRegistry)),
		fx.ResultTags(`group:"handler"`),
	),
	// biz handler
	fx.Annotate(
		web.NewBizHandler,
		fx.As(new(ginpkg.RouteRegistry)),
		fx.ResultTags(`group:"handler"`),
	),
)

var HandlerFxInvoke = fx.Invoke(
	fx.Annotate(
		RegisterRoutes,
		fx.ParamTags(``, `group:"handler"`),
	),
)

// RegisterRoutes 注册路由
func RegisterRoutes(engine *gin.Engine, registries []ginpkg.RouteRegistry) {
	for _, registry := range registries {
		registry.RegisterRoutes(engine)
	}
}
