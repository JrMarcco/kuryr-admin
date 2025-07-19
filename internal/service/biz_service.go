package service

import (
	"context"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BizService interface {
	Create(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error)

	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error)
	FindById(ctx context.Context, id uint64) (domain.BizInfo, error)
}

var _ BizService = (*DefaultBizService)(nil)

type DefaultBizService struct {
	db        *gorm.DB // db 数据库连接，用于开启事务
	bizRepo   repository.BizRepo
	userRepo  repository.UserRepo
	generator secret.Generator // biz secret 生成器
}

func (s *DefaultBizService) Create(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error) {
	bizSecret, err := s.generator.Generate(32)
	if err != nil {
		return domain.BizInfo{}, err
	}
	bi.BizSecret = bizSecret

	// 开启事务
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var innerErr error
		var res domain.BizInfo
		res, innerErr = s.bizRepo.CreateWithTx(ctx, tx, bi)
		if innerErr != nil {
			return innerErr
		}
		if bi.Id == 0 {
			// 当前为新建业务方，创建操作员
			// TODO: 这里要改成使用默认密码生成策略，然后创建成功后发送邮件通知，这里暂时写死
			var defaultPasswd []byte
			defaultPasswd, innerErr = bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if innerErr != nil {
				return innerErr
			}
			operator := domain.SysUser{
				Email:     bi.ContactEmail,
				Password:  string(defaultPasswd),
				RealName:  bi.Contact,
				UserType:  domain.UserTypeOperator,
				BizId:     res.Id,
				CreatedAt: res.CreatedAt,
				UpdatedAt: res.UpdatedAt,
			}

			_, innerErr = s.userRepo.CreateWithTx(ctx, tx, operator)
		}
		bi.Id = res.Id
		return nil
	})

	if err != nil {
		return domain.BizInfo{}, err
	}
	return bi, err
}

func (s *DefaultBizService) Count(ctx context.Context) (int64, error) {
	return s.bizRepo.Count(ctx)
}

func (s *DefaultBizService) List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error) {
	return s.bizRepo.List(ctx, offset, limit)
}

func (s *DefaultBizService) FindById(ctx context.Context, id uint64) (domain.BizInfo, error) {
	return s.bizRepo.FindById(ctx, id)
}

func NewBizService(
	db *gorm.DB, bizRepo repository.BizRepo, userRepo repository.UserRepo, generator secret.Generator,
) *DefaultBizService {
	return &DefaultBizService{
		db:        db,
		bizRepo:   bizRepo,
		userRepo:  userRepo,
		generator: generator,
	}
}
