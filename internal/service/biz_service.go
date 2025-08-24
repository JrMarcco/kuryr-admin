package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	pkggorm "github.com/JrMarcco/kuryr-admin/internal/pkg/gorm"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
	"github.com/JrMarcco/kuryr-admin/internal/search"
	businessv1 "github.com/JrMarcco/kuryr-api/api/go/business/v1"
)

type BizService interface {
	Save(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error)
	Delete(ctx context.Context, id uint64) error

	Search(ctx context.Context, criteria search.BizSearchCriteria, param *pkggorm.PaginationParam) (*pkggorm.PaginationResult[domain.BizInfo], error)
}

var _ BizService = (*DefaultBizService)(nil)

type DefaultBizService struct {
	grpcServerName string
	grpcClients    *client.Manager[businessv1.BusinessServiceClient]

	userRepo repository.UserRepo
}

func (s *DefaultBizService) Save(ctx context.Context, bi domain.BizInfo) (domain.BizInfo, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	_, err = grpcClient.Save(ctx, &businessv1.SaveRequest{BusinessInfo: s.domainToPb(bi)})
	cancel()

	if err != nil {
		return domain.BizInfo{}, fmt.Errorf("[kuryr-admin] failed to save biz info: %w", err)
	}

}

func (s *DefaultBizService) Delete(ctx context.Context, id uint64) error {
	// TODO: implement me
	panic("implement me")
}

func (s *DefaultBizService) Search(ctx context.Context, criteria search.BizSearchCriteria, param *pkggorm.PaginationParam) (*pkggorm.PaginationResult[domain.BizInfo], error) {
	// TODO: implement me
	panic("implement me")
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

func NewDefaultBizService(
	grpcServerName string, grpcClients *client.Manager[businessv1.BusinessServiceClient],
	userRepo repository.UserRepo,
) *DefaultBizService {
	return &DefaultBizService{
		grpcServerName: grpcServerName,
		userRepo:       userRepo,
		grpcClients:    grpcClients,
	}
}
