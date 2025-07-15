package service

import (
	"context"

	"github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (ginpkg.AuthUser, error)
	GenerateToken(ctx context.Context, au ginpkg.AuthUser) (accessToken, refreshToken string, err error)
	VerifyRefreshToken(ctx context.Context, token string) (ginpkg.AuthUser, error)
}

var _ UserService = (*JwtUserService)(nil)

type JwtUserService struct {
	userRepo repository.UserRepo

	atManager jwt.Manager[ginpkg.AuthUser] // access token manager
	stManager jwt.Manager[ginpkg.AuthUser] // refresh token manager
}

func (s *JwtUserService) Login(ctx context.Context, username string, password string) (ginpkg.AuthUser, error) {
	u, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return ginpkg.AuthUser{}, errs.ErrInvalidUser
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return ginpkg.AuthUser{}, errs.ErrInvalidUser
	}

	return ginpkg.AuthUser{
		Uid:      u.Id,
		Sid:      uuid.NewString(),
		UserType: u.UserType,
	}, nil
}

func (s *JwtUserService) GenerateToken(_ context.Context, au ginpkg.AuthUser) (accessToken, refreshToken string, err error) {
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

func (s *JwtUserService) VerifyRefreshToken(_ context.Context, token string) (ginpkg.AuthUser, error) {
	decrypt, err := s.stManager.Decrypt(token)
	if err != nil {
		return ginpkg.AuthUser{}, err
	}
	return decrypt.Data, nil
}

func NewJwtUserService(
	userRepo repository.UserRepo, atManager, stManager jwt.Manager[ginpkg.AuthUser],
) *JwtUserService {
	return &JwtUserService{
		userRepo:  userRepo,
		atManager: atManager,
		stManager: stManager,
	}
}
