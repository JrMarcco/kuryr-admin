package ioc

import (
	"time"

	easyjwt "github.com/JrMarcco/easy-kit/jwt"
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var JwtManagerOpt = fx.Provide(
	fx.Annotate(
		InitAccessTokenManager,
		fx.ResultTags(`name:"access-token-manager"`),
	),
	fx.Annotate(
		InitRefreshTokenManager,
		fx.ResultTags(`name:"refresh-token-manager"`),
	),
)

type jwtTokenConfig struct {
	Issuer     string `mapstructure:"issuer"`
	Expiration int    `mapstructure:"expiration"`
}

type jwtConfig struct {
	Private string         `mapstructure:"private"`
	Public  string         `mapstructure:"public"`
	Access  jwtTokenConfig `mapstructure:"access"`
	Refresh jwtTokenConfig `mapstructure:"refresh"`
}

// loadJwtConfig 加载配置
func loadJwtConfig() *jwtConfig {
	cfg := &jwtConfig{}
	if err := viper.UnmarshalKey("jwt", cfg); err != nil {
		panic(err)
	}
	return cfg
}

// InitAccessTokenManager 创建用于 Access Token 的 Manager
func InitAccessTokenManager() easyjwt.Manager[ginpkg.AuthUser] {
	jwtCfg := loadJwtConfig()

	claimsCfg := easyjwt.NewClaimsConfig(
		time.Duration(jwtCfg.Access.Expiration)*time.Second,
		easyjwt.WithIssuer(jwtCfg.Access.Issuer),
	)

	manager, err := easyjwt.NewEd25519ManagerBuilder[ginpkg.AuthUser](jwtCfg.Private, jwtCfg.Public).
		ClaimsConfig(claimsCfg).
		Build()
	if err != nil {
		panic(err)
	}
	return manager
}

// InitRefreshTokenManager 创建用于 Refresh Token 的 Manager
func InitRefreshTokenManager() easyjwt.Manager[ginpkg.AuthUser] {
	jwtCfg := loadJwtConfig()

	claimsCfg := easyjwt.NewClaimsConfig(
		time.Duration(jwtCfg.Refresh.Expiration)*time.Second,
		easyjwt.WithIssuer(jwtCfg.Refresh.Issuer),
	)

	manager, err := easyjwt.NewEd25519ManagerBuilder[ginpkg.AuthUser](jwtCfg.Private, jwtCfg.Public).
		ClaimsConfig(claimsCfg).
		Build()
	if err != nil {
		panic(err)
	}
	return manager
}
