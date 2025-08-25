package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/easy-kit/slice"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	pkggorm "github.com/JrMarcco/kuryr-admin/internal/pkg/gorm"
	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/search"
	businessv1 "github.com/JrMarcco/kuryr-api/api/go/business/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type BizService interface {
	Save(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error)
	Update(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error)
	Delete(ctx context.Context, id uint64) error

	Search(ctx context.Context, criteria search.BizSearchCriteria, param *pkggorm.PaginationParam) (*pkggorm.PaginationResult[domain.BizInfo], error)
	FindById(ctx context.Context, id uint64) (domain.BizInfo, error)
}

var _ BizService = (*DefaultBizService)(nil)

type DefaultBizService struct {
	grpcServerName string
	grpcClients    *client.Manager[businessv1.BusinessServiceClient]

	userRepo repository.UserRepo

	passwdGenerator secret.Generator
	logger          *zap.Logger
}

func (s *DefaultBizService) Save(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := grpcClient.Save(ctx, &businessv1.SaveRequest{BusinessInfo: s.domainToPb(bi)})

	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to save biz info: %w", err)
	}

	rtnBi := s.pbToDomain(resp.BusinessInfo)

	// 创建操作员
	passwd, err := s.passwdGenerator.Generate(16)
	if err != nil {
		s.logger.Error("[kuryr-admin] failed to generate password", zap.Error(err))
		return rtnBi, nil
	}

	user := domain.SysUser{
		Email:    bi.ContactEmail,
		Password: passwd,
		RealName: bi.Contact,
		UserType: domain.UserTypeOperator,
		BizId:    bi.Id,
	}

	_, err = s.userRepo.Save(ctx, user)
	if err != nil {
		s.logger.Error("[kuryr-admin] failed to save user", zap.Error(err))
		return rtnBi, nil
	}

	return rtnBi, nil
}

func (s *DefaultBizService) Update(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	fieldMask := &fieldmaskpb.FieldMask{
		Paths: []string{
			businessv1.FieldBizName,
			businessv1.FieldContact,
			businessv1.FieldContactEmail,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := grpcClient.Update(ctx, &businessv1.UpdateRequest{
		FieldMask:    fieldMask,
		BusinessInfo: s.domainToPb(bi),
	})

	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to update biz info: %w", err)
	}

	rtnBi := s.pbToDomain(resp.BusinessInfo)

	if bi.ContactEmail != "" {
		// 更新操作员信息
		user, err := s.userRepo.FindByEmail(ctx, bi.ContactEmail)
		if err != nil {
			s.logger.Error("[kuryr-admin] failed to find user", zap.Error(err))
			return rtnBi, nil
		}

		if user.BizId != rtnBi.Id {
			s.logger.Error("[kuryr-admin] user biz id not match", zap.Uint64("user_id", user.Id), zap.Uint64("biz_id", user.BizId), zap.Uint64("new_biz_id", rtnBi.Id))
			return rtnBi, nil
		}

		toUpdate := domain.SysUser{
			Id:    user.Id,
			Email: bi.ContactEmail,
		}

		if bi.Contact != "" {
			toUpdate.RealName = bi.Contact
		}

		_, err = s.userRepo.Save(ctx, toUpdate)
		if err != nil {
			s.logger.Error("[kuryr-admin] failed to save user", zap.Error(err))
			return rtnBi, nil
		}
	}

	return rtnBi, nil
}

func (s *DefaultBizService) Delete(ctx context.Context, id uint64) error {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = grpcClient.Delete(ctx, &businessv1.DeleteRequest{BizId: id})

	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to delete biz info: %w", err)
	}
	return nil
}

func (s *DefaultBizService) Search(ctx context.Context, criteria search.BizSearchCriteria, param *pkggorm.PaginationParam) (*pkggorm.PaginationResult[domain.BizInfo], error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := grpcClient.Search(ctx, &businessv1.SearchRequest{
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				businessv1.FieldBizKey,
				businessv1.FieldBizName,
				businessv1.FieldBizType,
			},
		},
		Offset:  int32(param.Offset),
		Limit:   int32(param.Limit),
		BizName: criteria.BizName,
	})

	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to search biz info: %w", err)
	}

	return &pkggorm.PaginationResult[domain.BizInfo]{
		Total: resp.Total,
		Records: slice.Map(resp.Records, func(_ int, src *businessv1.BusinessInfo) domain.BizInfo {
			return s.pbToDomain(src)
		}),
	}, nil
}

func (s *DefaultBizService) FindById(ctx context.Context, id uint64) (domain.BizInfo, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	fieldMask := &fieldmaskpb.FieldMask{
		Paths: []string{
			businessv1.FieldId,
			businessv1.FieldBizKey,
			businessv1.FieldBizName,
			businessv1.FieldBizType,
			businessv1.FieldBizSecret,
			businessv1.FieldContact,
			businessv1.FieldContactEmail,
			businessv1.FieldCreatorId,
			businessv1.FieldCreatedAt,
			businessv1.FieldUpdatedAt,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := grpcClient.FindById(ctx, &businessv1.FindByIdRequest{
		FieldMask: fieldMask,
		BizId:     id,
	})
	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to find biz info: %w", err)
	}

	bi := s.pbToDomain(resp.BusinessInfo)

	user, err := s.userRepo.FindById(ctx, bi.CreatorId)
	if err != nil {
		s.logger.Warn("[kuryr-admin] failed to find user", zap.Error(err))
	}

	user.Password = "" // 不返回密码
	bi.Creator = user

	return bi, nil
}

func (s *DefaultBizService) domainToPb(bi domain.BizInfo) *businessv1.BusinessInfo {
	return &businessv1.BusinessInfo{
		Id:           bi.Id,
		BizName:      bi.BizName,
		BizType:      string(bi.BizType),
		BizKey:       bi.BizKey,
		BizSecret:    bi.BizSecret,
		Contact:      bi.Contact,
		ContactEmail: bi.ContactEmail,
		CreatorId:    bi.CreatorId,
	}
}

func (s *DefaultBizService) pbToDomain(pb *businessv1.BusinessInfo) domain.BizInfo {
	return domain.BizInfo{
		Id:           pb.Id,
		BizType:      domain.BizType(pb.BizType),
		BizKey:       pb.BizKey,
		BizSecret:    pb.BizSecret,
		BizName:      pb.BizName,
		Contact:      pb.Contact,
		ContactEmail: pb.ContactEmail,
		CreatorId:    pb.CreatorId,
		UpdatedAt:    pb.UpdatedAt,
		CreatedAt:    pb.CreatedAt,
	}
}
func NewDefaultBizService(
	grpcServerName string,
	grpcClients *client.Manager[businessv1.BusinessServiceClient],
	userRepo repository.UserRepo,
	passwdGenerator secret.Generator,
	logger *zap.Logger,
) *DefaultBizService {
	return &DefaultBizService{
		grpcServerName: grpcServerName,
		grpcClients:    grpcClients,

		userRepo: userRepo,

		passwdGenerator: passwdGenerator,
		logger:          logger,
	}
}
