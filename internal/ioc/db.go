package ioc

import (
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBFxOpt = fx.Provide(InitDB)

func InitDB() *gorm.DB {
	type config struct {
		DSN string `mapstructure:"dsn"`
	}
	cfg := &config{}
	if err := viper.UnmarshalKey("db", cfg); err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
