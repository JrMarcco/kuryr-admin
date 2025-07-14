package repository

import (
	"context"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
)

type UserRepo interface {
	FindByUsername(ctx context.Context, username string) (domain.SysUser, error)
}

var _ UserRepo = (*DefaultUserRepo)(nil)

type DefaultUserRepo struct {
	userDAO dao.UserDAO
}

func (r *DefaultUserRepo) FindByUsername(ctx context.Context, username string) (domain.SysUser, error) {
	entity, err := r.userDAO.FindByUsername(ctx, username)
	if err != nil {
		return domain.SysUser{}, err
	}
	return r.toDomain(entity), nil
}

func (r *DefaultUserRepo) toDomain(entity dao.SysUser) domain.SysUser {
	return domain.SysUser{
		Id:        entity.Id,
		Username:  entity.Username,
		Password:  entity.Password,
		Email:     entity.Email,
		UserType:  domain.UserType(entity.UserType),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func NewUserRepo(userDAO dao.UserDAO) *DefaultUserRepo {
	return &DefaultUserRepo{userDAO: userDAO}
}
