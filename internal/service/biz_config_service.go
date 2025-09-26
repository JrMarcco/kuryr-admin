package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	commonv1 "github.com/JrMarcco/kuryr-api/api/go/common/v1"
	configv1 "github.com/JrMarcco/kuryr-api/api/go/config/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type BizConfigService interface {
	Save(ctx context.Context, bizConfig domain.BizConfig) (domain.BizConfig, error)
	FindByBizId(ctx context.Context, id uint64) (domain.BizConfig, error)
}

var _ BizConfigService = (*DefaultBizConfigService)(nil)

type DefaultBizConfigService struct {
	grpcServerName string
	grpcClients    *client.Manager[configv1.BizConfigServiceClient]
}

func (s *DefaultBizConfigService) Save(ctx context.Context, bizConfig domain.BizConfig) (domain.BizConfig, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.BizConfig{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	// 构建 grpc 请求
	pb := &configv1.BizConfig{
		Id:        bizConfig.Id,
		BizId:     bizConfig.BizId,
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
	defer cancel()

	if bizConfig.Id == 0 {
		// 创建配置
		resp, err := grpcClient.Save(ctx, &configv1.SaveRequest{BizConfig: pb})
		if err != nil {
			return domain.BizConfig{}, fmt.Errorf("[kuryr-admin] failed to save biz config: %w", err)
		}
		return s.pbToDomain(resp.BizConfig), nil
	}

	fieldMask := &fieldmaskpb.FieldMask{
		Paths: []string{
			configv1.FieldChannelConfig,
			configv1.FieldQuotaConfig,
			configv1.FieldCallbackConfig,
			configv1.FieldRateLimit,
		},
	}

	// 更新配置
	resp, err := grpcClient.Update(ctx, &configv1.UpdateRequest{
		FieldMask: fieldMask,
		BizConfig: pb,
	})

	if err != nil {
		return domain.BizConfig{}, fmt.Errorf("[kuryr-admin] failed to save biz config: %w", err)
	}
	return s.pbToDomain(resp.BizConfig), nil
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
	defer cancel()

	resp, err := grpcClient.FindByBizId(ctx, &configv1.FindByBizIdRequest{
		// FieldMask 传空会返回所有字段
		FieldMask: &fieldmaskpb.FieldMask{},
		BizId:     id,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.BizConfig{}, errs.ErrRecordNotFound
		}
		return domain.BizConfig{}, fmt.Errorf("[kuryr-admin] failed to get biz config: %w", err)
	}
	return s.pbToDomain(resp.BizConfig), nil
}

func (s *DefaultBizConfigService) pbToDomain(pb *configv1.BizConfig) domain.BizConfig {
	bizConfig := domain.BizConfig{
		Id:        pb.Id,
		BizId:     pb.BizId,
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
