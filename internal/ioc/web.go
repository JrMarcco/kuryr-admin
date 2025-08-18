package ioc

import (
	"time"

	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/web"
	ijwt "github.com/JrMarcco/kuryr-admin/internal/web/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var HandlerFxOpt = fx.Module(
	"web",
	fx.Provide(
		// redis jwt handler
		fx.Annotate(
			InitRedisJwtHandler,
			fx.As(new(ijwt.Handler)),
		),

		// user handler
		fx.Annotate(
			web.NewUserHandler,
			fx.As(new(pkggin.RouteRegistry)),
			fx.ResultTags(`group:"handler"`),
		),
		// biz handler
		fx.Annotate(
			web.NewBizHandler,
			fx.As(new(pkggin.RouteRegistry)),
			fx.ResultTags(`group:"handler"`),
		),

		// biz config handler
		fx.Annotate(
			web.NewBizConfigHandler,
			fx.As(new(pkggin.RouteRegistry)),
			fx.ResultTags(`group:"handler"`),
		),

		// provider handler
		fx.Annotate(
			web.NewProviderHandler,
			fx.As(new(pkggin.RouteRegistry)),
			fx.ResultTags(`group:"handler"`),
		),
	),
)

// var HandlerFxInvoke = fx.Invoke(
// 	fx.Annotate(
// 		RegisterRoutes,
// 		fx.ParamTags(``, ``, `group:"handler"`),
// 	),
// )

// // RegisterRoutes 注册路由
// // 注意：
// //
// //	这里声明 ioc.App 是为了保证 fx 在 ioc.RegisterRoutes 之前完成 engine.Use(middlewares...) 。
// //	ioc.RegisterRoutes 在 engine.Use(middlewares...) 之前调用
// //	会导致这里注册的路由“错过”这里注册的 middleware ，即导致 middleware 失效。
// func RegisterRoutes(_ *App, engine *gin.Engine, registries []pkggin.RouteRegistry) {
// 	for _, registry := range registries {
// 		registry.RegisterRoutes(engine)
// 	}
// }

func InitRedisJwtHandler(rc redis.Cmdable) ijwt.Handler {
	var expiration int
	if err := viper.UnmarshalKey("session.expiration", &expiration); err != nil {
		panic(err)
	}
	return ijwt.NewRedisHandler(rc, time.Duration(expiration)*time.Second)
}
