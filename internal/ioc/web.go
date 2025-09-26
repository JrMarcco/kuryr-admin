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

		// template handler
		fx.Annotate(
			web.NewTemplateHandler,
			fx.As(new(pkggin.RouteRegistry)),
			fx.ResultTags(`group:"handler"`),
		),
	),
)

func InitRedisJwtHandler(rc redis.Cmdable) ijwt.Handler {
	var expiration int
	if err := viper.UnmarshalKey("session.expiration", &expiration); err != nil {
		panic(err)
	}
	return ijwt.NewRedisHandler(rc, time.Duration(expiration)*time.Second)
}
