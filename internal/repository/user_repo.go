package repository

import (
	"context"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (domain.SysUser, error)
}

var _ UserRepo = (*DefaultUserRepo)(nil)

type DefaultUserRepo struct {
	userDAO dao.UserDAO
}

func (r *DefaultUserRepo) FindByEmail(ctx context.Context, email string) (domain.SysUser, error) {
	entity, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return domain.SysUser{}, err
	}
	return r.toDomain(entity), nil
}

func (r *DefaultUserRepo) toDomain(entity dao.SysUser) domain.SysUser {
	return domain.SysUser{
		Id:        entity.Id,
		Email:     entity.Email,
		Password:  entity.Password,
		RealName:  entity.RealName,
		UserType:  domain.UserType(entity.UserType),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func NewUserRepo(userDAO dao.UserDAO) *DefaultUserRepo {
	return &DefaultUserRepo{userDAO: userDAO}
}
