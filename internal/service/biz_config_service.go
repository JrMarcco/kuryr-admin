package service

import (
	"context"
	"fmt"

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

	// 构建 gRPC 请求参数
	pb := &configv1.BizConfig{
		BizId:     bizConfig.OwnerId,
		RateLimit: int32(bizConfig.RateLimit),
	}

	if bizConfig.ChannelConfig != nil {
		pb.ChannelConfig = s.convertToPbChannel(bizConfig.ChannelConfig)
	}
	if bizConfig.QuotaConfig != nil {
		pb.QuotaConfig = s.convertToQuota(bizConfig.QuotaConfig)
	}

	if bizConfig.CallbackConfig != nil {
		pb.CallbackConfig = s.convertToPbCallback(bizConfig.CallbackConfig)
	}

	resp, err := grpcClient.Save(ctx, &configv1.SaveRequest{Config: pb})
	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to save biz config: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("[kuryr-admin] failed to save biz config: [%s]", resp.ErrMsg)
	}
	return nil
}

// convertToPbChannel 渠道配置 proto buf
func (s *DefaultBizConfigService) convertToPbChannel(config *domain.ChannelConfig) *configv1.ChannelConfig {
	pbItems := make([]*configv1.ChannelItem, len(config.Channels))
	for i, item := range config.Channels {
		pbItems[i] = &configv1.ChannelItem{
			Channel:  item.Channel,
			Priority: int32(item.Priority),
			Enabled:  item.Enabled,
		}
	}

	return &configv1.ChannelConfig{
		Items:       pbItems,
		RetryPolicy: s.convertToPbRetry(config.RetryPolicyConfig),
	}
}

// convertToQuota 配额配置 proto buf
func (s *DefaultBizConfigService) convertToQuota(config *domain.QuotaConfig) *configv1.QuotaConfig {
	quota := &configv1.QuotaConfig{}
	if config.Daily != nil {
		quota.Daily = &configv1.Quota{
			Sms:   config.Daily.SMS,
			Email: config.Daily.Email,
		}
	}
	if config.Monthly != nil {
		quota.Monthly = &configv1.Quota{
			Sms:   config.Monthly.SMS,
			Email: config.Monthly.Email,
		}
	}
	return quota
}

// convertToPbCallback 回调配置 proto buf
func (s *DefaultBizConfigService) convertToPbCallback(config *domain.CallbackConfig) *configv1.CallbackConfig {
	return &configv1.CallbackConfig{
		ServiceName: config.ServiceName,
		RetryPolicy: s.convertToPbRetry(config.RetryPolicyConfig),
	}
}

// convertToPbRetry 重试机制 proto buf
func (s *DefaultBizConfigService) convertToPbRetry(config *domain.RetryConfig) *configv1.RetryPolicyConfig {
	return &configv1.RetryPolicyConfig{
		InitIntervalMs: int32(config.InitialInterval.Milliseconds()),
		MaxIntervalMs:  int32(config.MaxInterval.Milliseconds()),
		MaxRetryTimes:  config.MaxRetryTimes,
	}
}

func NewDefaultBizConfigService(grpcClients *client.Manager[configv1.BizConfigServiceClient]) *DefaultBizConfigService {
	return &DefaultBizConfigService{
		remoteSvcName: "kuryr",
		grpcClients:   grpcClients,
	}
}
