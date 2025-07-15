package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
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

	// repo
	fx.Provide(
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
