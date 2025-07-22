package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	configv1 "github.com/JrMarcco/kuryr-api/api/config/v1"
)

type BizConfigService interface {
	Create(ctx context.Context, bizConfig domain.BizConfig) error
}

var _ BizConfigService = (*DefaultBizConfigService)(nil)

type DefaultBizConfigService struct {
	remoteSvcName string
	grpcClients   *client.Manager[configv1.BizConfigServiceClient]
}

func (s *DefaultBizConfigService) Create(ctx context.Context, bizConfig domain.BizConfig) error {
	grpcClient, err := s.grpcClients.Get(s.remoteSvcName)
	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	resp, err := grpcClient.Save(ctx, &configv1.SaveRequest{})
	cancel()
	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to save biz config: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("[kuryr-admin] failed to save biz config: [%s]", resp.ErrMsg)
	}
	return nil
}

func NewDefaultBizConfigService(grpcClients *client.Manager[configv1.BizConfigServiceClient]) *DefaultBizConfigService {
	return &DefaultBizConfigService{
		remoteSvcName: "kuryr",
		grpcClients:   grpcClients,
	}
}
