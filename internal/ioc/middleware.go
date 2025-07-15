package ioc

import (
	"net/http"
	"net/url"
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
	InitCorsBuilder,
	fx.Annotate(
		InitJwtBuilder,
		fx.ParamTags(``, `name:"access-token-manager"`),
	),
)

func InitCorsBuilder() *middleware.CorsBuilder {
	type config struct {
		MaxAge    int      `mapstructure:"max_age"`
		Hostnames []string `mapstructure:"hostnames"`
	}
	cfg := &config{}
	if err := viper.UnmarshalKey("cors", cfg); err != nil {
		panic(err)
	}

	builder := middleware.NewCorsBuilder().
		AllowCredentials(true).
		AllowMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}).
		AllowHeaders([]string{"Content-Type", "Content-Length", "Authorization", "Accept", "Origin"}).
		MaxAge(time.Duration(cfg.MaxAge) * time.Second).
		AllowOriginFunc(func(origin string) bool {
			if origin == "" {
				return false
			}
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			reqHostname := u.Hostname()
			for _, hostname := range cfg.Hostnames {
				if reqHostname == hostname {
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
