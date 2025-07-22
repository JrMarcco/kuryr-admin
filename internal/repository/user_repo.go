package repository

import (
	"context"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, u domain.SysUser) (domain.SysUser, error)
	FindByEmail(ctx context.Context, email string) (domain.SysUser, error)
}

var _ UserRepo = (*DefaultUserRepo)(nil)

type DefaultUserRepo struct {
	dao dao.UserDao
}

func (r *DefaultUserRepo) CreateWithTx(ctx context.Context, tx *gorm.DB, u domain.SysUser) (domain.SysUser, error) {
	eu, err := r.dao.CreateWithTx(ctx, tx, dao.SysUser{
		Id:        u.Id,
		Email:     u.Email,
		Password:  u.Password,
		RealName:  u.RealName,
		UserType:  u.UserType.String(),
		BizId:     u.BizId,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
	if err != nil {
		return domain.SysUser{}, err
	}
	return r.toDomain(eu), nil
}

func (r *DefaultUserRepo) FindByEmail(ctx context.Context, email string) (domain.SysUser, error) {
	eu, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.SysUser{}, err
	}
	return r.toDomain(eu), nil
}

func (r *DefaultUserRepo) toDomain(eu dao.SysUser) domain.SysUser {
	return domain.SysUser{
		Id:        eu.Id,
		Email:     eu.Email,
		Password:  eu.Password,
		RealName:  eu.RealName,
		UserType:  domain.UserType(eu.UserType),
		BizId:     eu.BizId,
		CreatedAt: eu.CreatedAt,
		UpdatedAt: eu.UpdatedAt,
	}
}

func NewUserRepo(dao dao.UserDao) *DefaultUserRepo {
	return &DefaultUserRepo{dao: dao}
}
