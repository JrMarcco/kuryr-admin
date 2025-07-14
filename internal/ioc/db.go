package ioc

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
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

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
