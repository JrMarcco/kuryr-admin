package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ pkggin.RouteRegistry = (*BizConfigHandler)(nil)

type BizConfigHandler struct {
	svc service.BizConfigService
}

func (h *BizConfigHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/biz_config")

	v1.Handle(http.MethodPost, "/save", pkggin.B(h.Save))
	v1.Handle(http.MethodGet, "/find", pkggin.Q(h.Find))
}

type saveBizConfigReq struct {
	BizId          uint64          `json:"biz_id"`
	ChannelConfig  *channelConfig  `json:"channel_config,omitempty"`
	QuotaConfig    *quotaConfig    `json:"quota_config,omitempty"`
	CallbackConfig *callbackConfig `json:"callback_config,omitempty"`
	RateLimit      int32           `json:"rate_limit"`
}

type channelConfig struct {
	Channels          []channelItem `json:"channels"`
	RetryPolicyConfig *retryConfig  `json:"retry_policy_config,omitempty"`
}

type channelItem struct {
	Channel  int32 `json:"channel"`
	Priority int32 `json:"priority"`
	Enabled  bool  `json:"enabled"`
}

type quotaConfig struct {
	Daily   quota `json:"daily,omitempty"`
	Monthly quota `json:"monthly,omitempty"`
}

type quota struct {
	Sms   int32 `json:"sms"`
	Email int32 `json:"email"`
}

type callbackConfig struct {
	ServiceName       string       `json:"service_name"`
	RetryPolicyConfig *retryConfig `json:"retry_policy_config,omitempty"`
}

type retryConfig struct {
	InitialInterval int32 `json:"initial_interval"`
	MaxInterval     int32 `json:"max_interval"`
	MaxRetryTimes   int32 `json:"max_retry_times"`
}

func (h *BizConfigHandler) Save(ctx *gin.Context, req saveBizConfigReq) (pkggin.R, error) {
	if req.BizId == 0 {
		return pkggin.R{
			Code: http.StatusBadRequest,
			Msg:  "invalid biz_id, must be greater than 0",
		}, nil
	}

	if req.RateLimit < 0 {
		return pkggin.R{
			Code: http.StatusBadRequest,
			Msg:  "invalid rate_limit, must be greater than or equal to 0",
		}, nil
	}

	// 构建 domain.BizConfig
	bizConfig := domain.BizConfig{
		BizId:     req.BizId,
		RateLimit: req.RateLimit,
	}

	// 转换 ChannelConfig
	if req.ChannelConfig != nil {
		channelItems := make([]domain.ChannelItem, len(req.ChannelConfig.Channels))
		for i, item := range req.ChannelConfig.Channels {
			channelItems[i] = domain.ChannelItem{
				Channel:  item.Channel,
				Priority: item.Priority,
				Enabled:  item.Enabled,
			}
		}
		bizConfig.ChannelConfig = &domain.ChannelConfig{
			Channels: channelItems,
			RetryPolicyConfig: &domain.RetryConfig{
				InitialInterval: time.Duration(req.ChannelConfig.RetryPolicyConfig.InitialInterval) * time.Millisecond,
				MaxInterval:     time.Duration(req.ChannelConfig.RetryPolicyConfig.MaxInterval) * time.Millisecond,
				MaxRetryTimes:   req.ChannelConfig.RetryPolicyConfig.MaxRetryTimes,
			},
		}
	}

	// 转换 Quota
	if req.QuotaConfig != nil {
		bizConfig.QuotaConfig = &domain.QuotaConfig{
			Daily: &domain.Quota{
				Sms:   req.QuotaConfig.Daily.Sms,
				Email: req.QuotaConfig.Daily.Email,
			},
			Monthly: &domain.Quota{
				Sms:   req.QuotaConfig.Monthly.Sms,
				Email: req.QuotaConfig.Monthly.Email,
			},
		}
	}

	// 转换 CallbackConfig
	if req.CallbackConfig != nil {
		bizConfig.CallbackConfig = &domain.CallbackConfig{
			ServiceName: req.CallbackConfig.ServiceName,
			RetryPolicyConfig: &domain.RetryConfig{
				InitialInterval: time.Duration(req.CallbackConfig.RetryPolicyConfig.InitialInterval) * time.Millisecond,
				MaxInterval:     time.Duration(req.CallbackConfig.RetryPolicyConfig.MaxInterval) * time.Millisecond,
				MaxRetryTimes:   req.CallbackConfig.RetryPolicyConfig.MaxRetryTimes,
			},
		}
	}

	err := h.svc.Save(ctx, bizConfig)
	if err != nil {
		return pkggin.R{}, err
	}
	return pkggin.R{Code: http.StatusOK}, nil
}

type getBizConfigReq struct {
	BizId uint64 `json:"biz_id" form:"biz_id"`
}

func (h *BizConfigHandler) Find(ctx *gin.Context, req getBizConfigReq) (pkggin.R, error) {
	if req.BizId == 0 {
		return pkggin.R{
			Code: http.StatusBadRequest,
			Msg:  "invalid biz_id, must be greater than 0",
		}, nil
	}

	bizConfig, err := h.svc.FindByBizId(ctx, req.BizId)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return pkggin.R{Code: http.StatusOK}, nil
		}
		return pkggin.R{}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
		Data: bizConfig,
	}, nil
}

func NewBizConfigHandler(svc service.BizConfigService) *BizConfigHandler {
	return &BizConfigHandler{
		svc: svc,
	}
}
