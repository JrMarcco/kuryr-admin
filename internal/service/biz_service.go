package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	configv1 "github.com/JrMarcco/kuryr-api/api/config/v1"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BizService interface {
	Save(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error)
	Delete(ctx context.Context, id uint64) error

	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error)
	FindById(ctx context.Context, id uint64) (domain.BizInfo, error)
}

var _ BizService = (*DefaultBizService)(nil)

type DefaultBizService struct {
	grpcServerName string

	db *gorm.DB // db 数据库连接，用于开启事务

	repo     repository.BizRepo
	userRepo repository.UserRepo

	generator secret.Generator // biz secret 生成器

	grpcClients *client.Manager[configv1.BizConfigServiceClient]
}

func (s *DefaultBizService) Save(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error) {
	bizSecret, err := s.generator.Generate(32)
	if err != nil {
		return domain.BizInfo{}, err
	}
	bi.BizSecret = bizSecret

	// 开启事务
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var innerErr error
		var res domain.BizInfo
		res, innerErr = s.repo.SaveWithTx(ctx, tx, bi)
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

			_, innerErr = s.userRepo.SaveWithTx(ctx, tx, operator)
		}
		bi.Id = res.Id
		return nil
	})

	if err != nil {
		return domain.BizInfo{}, err
	}
	return bi, err
}

func (s *DefaultBizService) Delete(ctx context.Context, id uint64) error {
	// 删除 biz config
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	resp, err := grpcClient.Delete(ctx, &configv1.DeleteRequest{Id: id})
	cancel()

	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to delete biz config: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("[kuryr-admin] failed to delete biz config: [%s]", resp.ErrMsg)
	}

	// 开启事务，删除业务以及对应操作员信息
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if innerErr := s.userRepo.DeleteByBizIdWithTx(ctx, tx, id); innerErr != nil {
			return innerErr
		}
		if innerErr := s.repo.DeleteWithTx(ctx, tx, id); innerErr != nil {
			return innerErr
		}
		return nil
	})
}

func (s *DefaultBizService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *DefaultBizService) List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *DefaultBizService) FindById(ctx context.Context, id uint64) (domain.BizInfo, error) {
	return s.repo.FindById(ctx, id)
}

func NewDefaultBizService(
	grpcServerName string,
	db *gorm.DB, bizRepo repository.BizRepo, userRepo repository.UserRepo, generator secret.Generator,
	grpcClients *client.Manager[configv1.BizConfigServiceClient],
) *DefaultBizService {
	return &DefaultBizService{
		grpcServerName: grpcServerName,
		db:             db,
		repo:           bizRepo,
		userRepo:       userRepo,
		generator:      generator,
		grpcClients:    grpcClients,
	}
}
