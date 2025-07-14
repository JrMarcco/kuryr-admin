package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
	"go.uber.org/fx"
)

var RepoFxOpt = fx.Options(
	// dao
	fx.Provide(
		fx.Annotate(
			dao.NewUserDAO,
			fx.As(new(dao.UserDAO)),
		),
	),
	// cache

	// repo
	fx.Provide(
		fx.Annotate(
			repository.NewUserRepo,
			fx.As(new(repository.UserRepo)),
		),
	),
)
