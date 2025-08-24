package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
	"go.uber.org/fx"
)

var RepoFxOpt = fx.Module(
	"repository",
	// dao
	fx.Provide(
		// user dao
		fx.Annotate(
			dao.NewUserDAO,
			fx.As(new(dao.UserDao)),
		),
	),
	// cache

	// repo
	fx.Provide(
		// user repo
		fx.Annotate(
			repository.NewUserRepo,
			fx.As(new(repository.UserRepo)),
		),
	),
)
