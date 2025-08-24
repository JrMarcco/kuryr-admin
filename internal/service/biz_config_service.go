package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	commonv1 "github.com/JrMarcco/kuryr-api/api/go/common/v1"
	configv1 "github.com/JrMarcco/kuryr-api/api/go/config/v1"
	"gorm.io/gorm"
)

type BizConfigService interface {
	Save(ctx context.Context, bizConfig domain.BizConfig) error
	Delete(ctx context.Context, id uint64) error
	FindByBizId(ctx context.Context, id uint64) (domain.BizConfig, error)
}

var _ BizConfigService = (*DefaultBizConfigService)(nil)

type DefaultBizConfigService struct {
	grpcServerName string
	grpcClients    *client.Manager[configv1.BizConfigServiceClient]
}

func (s *DefaultBizConfigService) Save(ctx context.Context, bizConfig domain.BizConfig) error {

	bi, err := s.bizRepo.FindById(ctx, bizConfig.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("[kuryr-admin] failed to find biz info [ %d ]", bizConfig.Id)
		}
		return fmt.Errorf("[kuryr-admin] failed to find biz info: %w", err)
	}

	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	// 构建 grpc 请求
	pb := &configv1.BizConfig{
		BizId:     bi.Id,
		OwnerType: string(bizConfig.OwnerType),
		RateLimit: bizConfig.RateLimit,
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

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	resp, err := grpcClient.Save(ctx, &configv1.SaveRequest{Config: pb})
	cancel()

	if err != nil {
		return fmt.Errorf("[kuryr-admin] failed to save biz config: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("[kuryr-admin] failed to save biz config: [ %s ]", resp.ErrMsg)
	}
	return nil
}

func (s *DefaultBizConfigService) Delete(ctx context.Context, id uint64) error {
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
		return fmt.Errorf("[kuryr-admin] failed to delete biz config: [ %s ]", resp.ErrMsg)
	}
	return nil
}

// convertToPbChannel 渠道配置 proto buf
func (s *DefaultBizConfigService) convertToPbChannel(config *domain.ChannelConfig) *configv1.ChannelConfig {
	pbItems := make([]*configv1.ChannelItem, len(config.Channels))
	for i, item := range config.Channels {
		pbItems[i] = &configv1.ChannelItem{
			Channel:  commonv1.Channel(item.Channel),
			Priority: item.Priority,
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
			Sms:   config.Daily.Sms,
			Email: config.Daily.Email,
		}
	}
	if config.Monthly != nil {
		quota.Monthly = &configv1.Quota{
			Sms:   config.Monthly.Sms,
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

func (s *DefaultBizConfigService) FindByBizId(ctx context.Context, id uint64) (domain.BizConfig, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.BizConfig{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	resp, err := grpcClient.FindById(ctx, &configv1.FindByIdRequest{Id: id})
	cancel()

	if err != nil {
		return domain.BizConfig{}, fmt.Errorf("[kuryr-admin] failed to get biz config: %w", err)
	}
	if resp.ErrCode == commonv1.ErrCode_BIZ_CONFIG_NOT_FOUND {
		return domain.BizConfig{}, errs.ErrRecordNotFound
	}
	return s.pbToDomain(resp.Config), nil
}

func (s *DefaultBizConfigService) pbToDomain(pb *configv1.BizConfig) domain.BizConfig {
	bizConfig := domain.BizConfig{
		Id:        pb.BizId,
		OwnerType: domain.BizType(pb.OwnerType),
		RateLimit: pb.RateLimit,
	}

	if pb.ChannelConfig != nil {
		channelConfig := &domain.ChannelConfig{
			Channels: make([]domain.ChannelItem, len(pb.ChannelConfig.Items)),
		}

		for index, item := range pb.ChannelConfig.Items {
			channelConfig.Channels[index] = domain.ChannelItem{
				Channel:  int32(item.Channel),
				Priority: item.Priority,
				Enabled:  item.Enabled,
			}
		}

		if pb.ChannelConfig.RetryPolicy != nil {
			retryPolicyConfig := s.convertRetry(pb.ChannelConfig.RetryPolicy)
			channelConfig.RetryPolicyConfig = retryPolicyConfig
		}
		bizConfig.ChannelConfig = channelConfig
	}

	if pb.QuotaConfig != nil {
		quotaConfig := &domain.QuotaConfig{}
		if pb.QuotaConfig.Daily != nil {
			dailyQuota := pb.QuotaConfig.Daily
			quotaConfig.Daily = &domain.Quota{
				Sms:   dailyQuota.Sms,
				Email: dailyQuota.Email,
			}
		}
		if pb.QuotaConfig.Monthly != nil {
			monthlyQuota := pb.QuotaConfig.Monthly
			quotaConfig.Monthly = &domain.Quota{
				Sms:   monthlyQuota.Sms,
				Email: monthlyQuota.Email,
			}
		}
		bizConfig.QuotaConfig = quotaConfig
	}

	if pb.CallbackConfig != nil {
		callbackConfig := &domain.CallbackConfig{
			ServiceName: pb.CallbackConfig.ServiceName,
		}

		if pb.CallbackConfig.RetryPolicy != nil {
			retryPolicyConfig := s.convertRetry(pb.CallbackConfig.RetryPolicy)
			callbackConfig.RetryPolicyConfig = retryPolicyConfig
		}
		bizConfig.CallbackConfig = callbackConfig
	}
	return bizConfig
}

func (s *DefaultBizConfigService) convertRetry(pbRetry *configv1.RetryPolicyConfig) *domain.RetryConfig {
	return &domain.RetryConfig{
		InitialInterval: time.Duration(pbRetry.InitIntervalMs),
		MaxInterval:     time.Duration(pbRetry.MaxIntervalMs),
		MaxRetryTimes:   pbRetry.MaxRetryTimes,
	}
}

func NewDefaultBizConfigService(
	grpcServerName string, grpcClients *client.Manager[configv1.BizConfigServiceClient],
) *DefaultBizConfigService {
	return &DefaultBizConfigService{
		grpcServerName: grpcServerName,
		grpcClients:    grpcClients,
	}
}
