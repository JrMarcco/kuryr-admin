package ioc

import (
	"strings"

	easyjwt "github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/easy-kit/set"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/gin/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var MiddlewareBuilderOpt = fx.Provide(
	fx.Annotate(
		InitJwtBuilder,
		fx.As(new(middleware.Builder)),
		fx.ResultTags(`group:"middleware-builder"`),
	),
)

func InitJwtBuilder(rc redis.Cmdable, jwtManager easyjwt.Manager[domain.AuthUser]) *middleware.JwtBuilder {
	var ignores []string
	if err := viper.UnmarshalKey("ignores", &ignores); err != nil {
		panic(err)
	}

	ts, err := set.NewTreeSet[string](strings.Compare)
	if err != nil {
		panic(err)
	}
	for _, ignore := range ignores {
		ts.Add(ignore)
	}
	return middleware.NewJwtBuilder(rc, jwtManager, ts)
}

func InitAccessLogBuilder() *middleware.AccessLogBuilder {
	return &middleware.AccessLogBuilder{}
}
