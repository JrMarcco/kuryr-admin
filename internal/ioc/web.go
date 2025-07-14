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
		fx.ParamTags(``, ``, `group:"handler"`),
	),
)

// RegisterRoutes 注册路由
//
// 注意：
//
//	这里声明 *App 是为了让 fx 在 RegisterRoutes 被调用之前，先初始化 ioc.App。
//	RegisterRoutes 在 ioc.App 之前调用会导致这里注册的路由“错过” ioc.APP 初始化时注册的 middleware 导致 middleware 失效。
func RegisterRoutes(_ *App, engine *gin.Engine, registries []ginpkg.Registry) {
	for _, registry := range registries {
		registry.RegisterRoutes(engine)
	}
}
