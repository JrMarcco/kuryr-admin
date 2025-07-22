package service

import (
	"context"

	"github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accountTypeEmail  = "email"
	accountTypeMobile = "mobile"

	verifyTypePasswd = "passwd"
	verifyTypeCode   = "code"
)

type UserService interface {
	LoginWithType(ctx context.Context, account string, credential string, accountType, VerifyType string) (pkggin.AuthUser, error)
	GenerateToken(ctx context.Context, au pkggin.AuthUser) (accessToken, refreshToken string, err error)
	VerifyRefreshToken(ctx context.Context, token string) (pkggin.AuthUser, error)
}

var _ UserService = (*JwtUserService)(nil)

type JwtUserService struct {
	repo repository.UserRepo

	atManager jwt.Manager[pkggin.AuthUser] // access token manager
	stManager jwt.Manager[pkggin.AuthUser] // refresh token manager
}

func (s *JwtUserService) LoginWithType(
	ctx context.Context, account string, credential string, accountType, VerifyType string,
) (pkggin.AuthUser, error) {
	var (
		u   domain.SysUser
		err error
	)
	switch accountType {
	case accountTypeEmail:
		u, err = s.repo.FindByEmail(ctx, account)
	default:
		return pkggin.AuthUser{}, errs.ErrInvalidAccountType
	}

	if err != nil {
		return pkggin.AuthUser{}, errs.ErrInvalidUser
	}

	switch VerifyType {
	case verifyTypePasswd:
		err = s.verifyPasswd(u, credential)
	default:
		return pkggin.AuthUser{}, errs.ErrInvalidVerifyType
	}
	if err != nil {
		return pkggin.AuthUser{}, err
	}

	return pkggin.AuthUser{
		Sid:      uuid.NewString(),
		Bid:      u.BizId,
		Uid:      u.Id,
		UserType: u.UserType,
	}, nil
}

func (s *JwtUserService) verifyPasswd(u domain.SysUser, credential string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(credential)); err != nil {
		return errs.ErrInvalidUser
	}
	return nil
}

func (s *JwtUserService) GenerateToken(_ context.Context, au pkggin.AuthUser) (accessToken, refreshToken string, err error) {
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

func (s *JwtUserService) VerifyRefreshToken(_ context.Context, token string) (pkggin.AuthUser, error) {
	decrypt, err := s.stManager.Decrypt(token)
	if err != nil {
		return pkggin.AuthUser{}, err
	}
	return decrypt.Data, nil
}

func NewJwtUserService(
	repo repository.UserRepo, atManager, stManager jwt.Manager[pkggin.AuthUser],
) *JwtUserService {
	return &JwtUserService{
		repo:      repo,
		atManager: atManager,
		stManager: stManager,
	}
}
