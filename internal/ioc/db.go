package ioc

import (
	"time"

	pkggorm "github.com/JrMarcco/kuryr-admin/internal/pkg/gorm"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/snowflake"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBFxOpt = fx.Module(
	"db",
	fx.Provide(
		InitDB,
		snowflake.NewGenerator,
	),
)

func InitDB(zLogger *zap.Logger) *gorm.DB {
	type config struct {
		LogLevel                  string `mapstructure:"log_level"`
		SlowThreshold             int    `mapstructure:"slow_threshold"`
		IgnoreRecordNotFoundError bool   `mapstructure:"ignore_record_not_found_error"`
		DSN                       string `mapstructure:"dsn"`
	}
	cfg := config{}
	if err := viper.UnmarshalKey("db", &cfg); err != nil {
		panic(err)
	}

	var level logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
	default:
		panic("invalid log level")
	}

	gormLogger := pkggorm.NewZapLogger(
		zLogger,
		pkggorm.WithLogLevel(level),
		pkggorm.WithSlowThreshold(time.Duration(cfg.SlowThreshold)*time.Millisecond),
		pkggorm.WithIgnoreRecordNotFoundError(cfg.IgnoreRecordNotFoundError),
	)

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}
	return db
}
