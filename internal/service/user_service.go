package service

import (
	"context"

	"github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (accessToken, refreshToken string, err error)
	RefreshToken(ctx context.Context, rt string) (accessToken, refreshToken string, err error)
}

var _ UserService = (*DefaultUserService)(nil)

type DefaultUserService struct {
	userRepo repository.UserRepo

	atManager jwt.Manager[domain.AuthUser] // access token manager
	stManager jwt.Manager[domain.AuthUser] // refresh token manager
}

func (s *DefaultUserService) Login(ctx context.Context, username string, password string) (accessToken, refreshToken string, err error) {
	u, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", "", err
	}

	return s.generateToken(domain.AuthUser{
		SSId: uuid.NewString(),
		Id:   u.Id,
	})
}

func (s *DefaultUserService) RefreshToken(ctx context.Context, rt string) (accessToken, refreshToken string, err error) {
	decrypt, err := s.stManager.Decrypt(rt)
	if err != nil {
		return "", "", err
	}
	au := decrypt.Data
	return s.generateToken(au)
}

func (s *DefaultUserService) generateToken(au domain.AuthUser) (accessToken, refreshToken string, err error) {
	// access token
	at, err := s.atManager.Encrypt(au)
	if err != nil {
		return "", "", err
	}
	// refresh token
	rt, err := s.stManager.Encrypt(au)
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}

func NewUserService(
	userRepo repository.UserRepo, atManager, stManager jwt.Manager[domain.AuthUser],
) *DefaultUserService {
	return &DefaultUserService{
		userRepo:  userRepo,
		atManager: atManager,
		stManager: stManager,
	}
}
