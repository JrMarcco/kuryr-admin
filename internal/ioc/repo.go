package ioc

import (
	"time"

	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/repository/cache"
	iredis "github.com/JrMarcco/kuryr-admin/internal/repository/cache/redis"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var RepoFxOpt = fx.Options(
	// dao
	fx.Provide(
		// user dao
		fx.Annotate(
			dao.NewUserDAO,
			fx.As(new(dao.UserDAO)),
		),
		// biz dao
		fx.Annotate(
			dao.NewBizDAO,
			fx.As(new(dao.BizDAO)),
		),
	),
	// cache
	fx.Provide(
		// session cache
		fx.Annotate(
			InitRedisSessionCache,
			fx.As(new(cache.SessionCache)),
		),
	),
	// repo
	fx.Provide(
		// session repo
		fx.Annotate(
			repository.NewDefaultSessionRepo,
			fx.As(new(repository.SessionRepo)),
		),
		// user repo
		fx.Annotate(
			repository.NewUserRepo,
			fx.As(new(repository.UserRepo)),
		),
		// biz repo
		fx.Annotate(
			repository.NewBizRepo,
			fx.As(new(repository.BizRepo)),
		),
	),
)

func InitRedisSessionCache(rc redis.Cmdable) *iredis.RSessionCache {
	var expiration int
	if err := viper.UnmarshalKey("session.expiration", &expiration); err != nil {
		panic(err)
	}
	return iredis.NewRSessionCache(rc, time.Duration(expiration)*time.Second)
}
