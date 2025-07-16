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
		fx.ParamTags(``, ``, `group:"handler"`),
	),
)

// RegisterRoutes 注册路由
// 注意：
//
//	这里声明 ioc.App 是为了保证 fx 在 ioc.RegisterRoutes 之前完成 engine.Use(middlewares...) 。
//	ioc.RegisterRoutes 在 engine.Use(middlewares...) 之前调用
//	会导致这里注册的路由“错过”这里注册的 middleware ，即导致 middleware 失效。
func RegisterRoutes(_ *App, engine *gin.Engine, registries []ginpkg.RouteRegistry) {
	for _, registry := range registries {
		registry.RegisterRoutes(engine)
	}
}
