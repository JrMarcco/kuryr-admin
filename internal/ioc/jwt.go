package ioc

import (
	"time"

	easyjwt "github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var JwtManagerOpt = fx.Provide(InitJwtManager)

func InitJwtManager() easyjwt.Manager[domain.AuthUser] {
	type jwtConfig struct {
		Expiration int    `mapstructure:"expiration"`
		Private    string `mapstructure:"private"`
		Public     string `mapstructure:"public"`
	}

	jwtCfg := &jwtConfig{}
	if err := viper.UnmarshalKey("jwt", jwtCfg); err != nil {
		panic(err)
	}

	claimsCfg := easyjwt.NewClaimsConfig(
		time.Duration(jwtCfg.Expiration),
		easyjwt.WithIssuer("kuryr-admin"),
	)

	manager, err := easyjwt.NewEd25519ManagerBuilder[domain.AuthUser](jwtCfg.Private, jwtCfg.Public).
		ClaimsConfig(claimsCfg).
		Build()
	if err != nil {
		panic(err)
	}
	return manager
}
