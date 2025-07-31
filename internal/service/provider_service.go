package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/easy-kit/slice"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	pkggorm "github.com/JrMarcco/kuryr-admin/internal/pkg/gorm"
	"github.com/JrMarcco/kuryr-admin/internal/search"
	commonv1 "github.com/JrMarcco/kuryr-api/api/common/v1"
	providerv1 "github.com/JrMarcco/kuryr-api/api/provider/v1"
)

type ProviderService interface {
	Save(ctx context.Context, provider domain.Provider) error
	Search(ctx context.Context, criteria search.ProviderCriteria, param *pkggorm.PaginationParam) (*pkggorm.PaginationResult[domain.Provider], error)
}

var _ ProviderService = (*DefaultProviderService)(nil)

type DefaultProviderService struct {
	grpcServerName string
	grpcClients    *client.Manager[providerv1.ProviderServiceClient]
}

func (s *DefaultProviderService) Save(ctx context.Context, provider domain.Provider) error {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	resp, err := grpcClient.Save(ctx, &providerv1.SaveRequest{Provider: s.domainToPb(provider)})
	cancel()

	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to save provider: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("[kuryr-admin] failed to save provider: [ %s ]", resp.ErrMsg)
	}
	return nil
}

func (s *DefaultProviderService) domainToPb(provider domain.Provider) *providerv1.Provider {
	return &providerv1.Provider{
		Id:               provider.Id,
		ProviderName:     provider.ProviderName,
		Channel:          commonv1.Channel(provider.Channel),
		Endpoint:         provider.Endpoint,
		RegionId:         provider.RegionId,
		AppId:            provider.AppId,
		ApiKey:           provider.ApiKey,
		ApiSecret:        provider.ApiSecret,
		Weight:           provider.Weight,
		QpsLimit:         provider.QpsLimit,
		DailyLimit:       provider.DailyLimit,
		AuditCallbackUrl: provider.AuditCallbackUrl,
	}
}

func (s *DefaultProviderService) Search(ctx context.Context, criteria search.ProviderCriteria, param *pkggorm.PaginationParam) (*pkggorm.PaginationResult[domain.Provider], error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	if param == nil {
		param = &pkggorm.PaginationParam{
			Offset: 0,
			Limit:  10,
		}
	}

	req := &providerv1.SearchRequest{
		ProviderName: criteria.ProviderName,
		Channel:      commonv1.Channel(criteria.Channel),
		Offset:       int32(param.Offset),
		Limit:        int32(param.Limit),
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	resp, err := grpcClient.Search(ctx, req)
	cancel()

	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to search providers: %w", err)
	}

	if resp.Total == 0 {
		return pkggorm.NewPaginationResult[domain.Provider]([]domain.Provider{}, 0), nil
	}

	providers := slice.Map(resp.Providers, func(_ int, src *providerv1.Provider) domain.Provider {
		return s.pbToDomain(src)
	})
	return pkggorm.NewPaginationResult(providers, resp.Total), nil
}

func (s *DefaultProviderService) pbToDomain(pb *providerv1.Provider) domain.Provider {
	return domain.Provider{
		Id:               pb.Id,
		ProviderName:     pb.ProviderName,
		Channel:          int32(pb.Channel),
		Endpoint:         pb.Endpoint,
		RegionId:         pb.RegionId,
		AppId:            pb.AppId,
		ApiKey:           pb.ApiKey,
		ApiSecret:        pb.ApiSecret,
		Weight:           pb.Weight,
		QpsLimit:         pb.QpsLimit,
		DailyLimit:       pb.DailyLimit,
		AuditCallbackUrl: pb.AuditCallbackUrl,
		ActiveStatus:     pb.ActiveStatus,
	}
}

func NewDefaultProviderService(
	grpcServerName string, grpcClients *client.Manager[providerv1.ProviderServiceClient],
) *DefaultProviderService {
	return &DefaultProviderService{
		grpcServerName: grpcServerName,
		grpcClients:    grpcClients,
	}
}
