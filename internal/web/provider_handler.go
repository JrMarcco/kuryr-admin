package web

import (
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ pkggin.RouteRegistry = (*ProviderHandler)(nil)

type ProviderHandler struct {
	svc service.ProviderService
}

func (h *ProviderHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/provider")

	v1.Handle(http.MethodPost, "/save", pkggin.B(h.Save))
	v1.Handle(http.MethodGet, "/list", pkggin.W(h.List))
	v1.Handle(http.MethodGet, "/find_by_channel", pkggin.Q(h.FindByChannel))
}

type saveProviderReq struct {
	ProviderName string `json:"provider_name"`
	Channel      int32  `json:"channel"` // sms = 1 / email = 2

	Endpoint string `json:"endpoint"`
	RegionId string `json:"region_id"`

	AppId     string `json:"app_id"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`

	Weight     int32 `json:"weight"`
	QpsLimit   int32 `json:"qps_limit"`
	DailyLimit int32 `json:"daily_limit"`

	AuditCallbackUrl string `json:"audit_callback_url"`
}

func (h *ProviderHandler) Save(ctx *gin.Context, req saveProviderReq) (pkggin.R, error) {
	provider := domain.Provider{
		ProviderName:     req.ProviderName,
		Channel:          req.Channel,
		Endpoint:         req.Endpoint,
		RegionId:         req.RegionId,
		AppId:            req.AppId,
		ApiKey:           req.ApiKey,
		ApiSecret:        req.ApiSecret,
		Weight:           req.Weight,
		QpsLimit:         req.QpsLimit,
		DailyLimit:       req.DailyLimit,
		AuditCallbackUrl: req.AuditCallbackUrl,
	}

	err := h.svc.Save(ctx, provider)
	if err != nil {
		return pkggin.R{}, err
	}
	return pkggin.R{Code: http.StatusOK}, nil
}

type listProviderResp struct {
	Records []domain.Provider `json:"records"`
}

func (h *ProviderHandler) List(ctx *gin.Context) (pkggin.R, error) {
	res, err := h.svc.List(ctx)
	if err != nil {
		return pkggin.R{}, err
	}

	return pkggin.R{
		Code: http.StatusOK,
		Data: listProviderResp{
			Records: res,
		},
	}, nil
}

type findByChannelReq struct {
	Channel int32 `json:"channel" form:"channel"`
}

func (h *ProviderHandler) FindByChannel(ctx *gin.Context, req findByChannelReq) (pkggin.R, error) {
	res, err := h.svc.FindByChannel(ctx, req.Channel)
	if err != nil {
		return pkggin.R{}, err
	}

	return pkggin.R{
		Code: http.StatusOK,
		Data: res,
	}, nil
}

func NewProviderHandler(svc service.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		svc: svc,
	}
}
