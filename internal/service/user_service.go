package service

import (
	"context"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (domain.SysUser, error)
}

var _ UserService = (*DefaultUserService)(nil)

type DefaultUserService struct {
	userRepo repository.UserRepo
}

func (d *DefaultUserService) Login(ctx context.Context, username string, password string) (domain.SysUser, error) {
	u, err := d.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return domain.SysUser{}, errs.ErrInvalidUser
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return domain.SysUser{}, errs.ErrInvalidUser
	}

	return u, nil
}

func NewUserService(userRepo repository.UserRepo) *DefaultUserService {
	return &DefaultUserService{userRepo: userRepo}
}
