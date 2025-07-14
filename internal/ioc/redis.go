package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var RedisFxOpt = fx.Provide(
	fx.Annotate(
		InitRedis,
		fx.As(new(redis.Cmdable)),
	),
)

func InitRedis() *redis.Client {
	type config struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
	}
	cfg := &config{}
	if err := viper.UnmarshalKey("redis", cfg); err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})
}
