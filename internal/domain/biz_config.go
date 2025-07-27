package domain

import "time"

// BizConfig 业务方配置领域对象。
type BizConfig struct {
	Id             uint64          `json:"id"`             // 对应 ( biz_info.id )
	OwnerType      BizType         `json:"owner_type"`     // 业务类型
	ChannelConfig  *ChannelConfig  `json:"channel_config"` // 渠道配置
	QuotaConfig    *QuotaConfig    `json:"quota_config"`   // 配额配置
	CallbackConfig *CallbackConfig `json:"callback_config"`
	RateLimit      int             `json:"rate_limit"`
}

type RetryConfig struct {
	InitialInterval time.Duration `json:"initial_interval"`
	MaxInterval     time.Duration `json:"max_interval"`
	MaxRetryTimes   int32         `json:"max_retry_times"`
}

type ChannelItem struct {
	Channel  string `json:"channel"`
	Priority int    `json:"priority"`
	Enabled  bool   `json:"enabled"`
}

type ChannelConfig struct {
	Channels          []ChannelItem `json:"channels"`
	RetryPolicyConfig *RetryConfig  `json:"retry_policy_config"`
}

type Quota struct {
	SMS   int32 `json:"sms"`
	Email int32 `json:"email"`
}

type QuotaConfig struct {
	Daily   *Quota `json:"daily"`
	Monthly *Quota `json:"monthly"`
}

type CallbackConfig struct {
	ServiceName       string       `json:"service_name"`
	RetryPolicyConfig *RetryConfig `json:"retry_policy_config"`
}
