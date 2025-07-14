package ioc

import (
	"strings"
	"time"

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
		InitCorsBuilder,
		fx.As(new(middleware.Builder)),
		fx.ResultTags(`group:"middleware-builder"`),
	),
	fx.Annotate(
		InitJwtBuilder,
		fx.As(new(middleware.Builder)),
		fx.ResultTags(`group:"middleware-builder"`),
	),
)

func InitCorsBuilder() *middleware.CorsBuilder {
	type config struct {
		MaxAge      int      `mapstructure:"max_age"`
		DomainNames []string `mapstructure:"domain_names"`
	}
	cfg := &config{}
	if err := viper.UnmarshalKey("cors", cfg); err != nil {
		panic(err)
	}

	builder := middleware.NewCorsBuilder().
		MaxAge(time.Duration(cfg.MaxAge) * time.Second).
		AllowOriginFunc(func(origin string) bool {
			for _, domainName := range cfg.DomainNames {
				if strings.Contains(origin, domainName) {
					return true
				}
			}
			return false
		})
	return builder
}

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
