package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/easy-kit/slice"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	commonv1 "github.com/JrMarcco/kuryr-api/api/common/v1"
	providerv1 "github.com/JrMarcco/kuryr-api/api/provider/v1"
)

type ProviderService interface {
	Save(ctx context.Context, provider domain.Provider) error
	List(ctx context.Context) ([]domain.Provider, error)
	FindByChannel(ctx context.Context, channel int32) ([]domain.Provider, error)
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

func (s *DefaultProviderService) List(ctx context.Context) ([]domain.Provider, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	resp, err := grpcClient.List(ctx, &providerv1.ListRequest{})
	cancel()

	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to list providers: %w", err)
	}

	if len(resp.Providers) == 0 {
		return []domain.Provider{}, nil
	}

	providers := slice.Map(resp.Providers, func(_ int, src *providerv1.Provider) domain.Provider {
		return s.pbToDomain(src)
	})
	return providers, nil
}

func (s *DefaultProviderService) FindByChannel(ctx context.Context, channel int32) ([]domain.Provider, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	resp, err := grpcClient.FindByChannel(ctx, &providerv1.FindByChannelRequest{Channel: commonv1.Channel(channel)})
	cancel()

	if err != nil {
		return nil, fmt.Errorf("[kuryr-admin] failed to find providers by channel: %w", err)
	}

	if len(resp.Providers) == 0 {
		return []domain.Provider{}, nil
	}

	providers := slice.Map(resp.Providers, func(_ int, src *providerv1.Provider) domain.Provider {
		return s.pbToDomain(src)
	})
	return providers, nil
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
