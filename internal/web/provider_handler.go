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
	v1.Handle(http.MethodGet, "/list", pkggin.Q(h.List))
}

type saveProviderReq struct {
	ProviderName string `json:"provider_name"`
	Channel      int32  `json:"channel"` // sms = 1 / email = 2

	Endpoint string `json:"endpoint"`
	RegionId string `json:"region_id"`

	AppId     string `json:"app_id"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`

	Weight     int `json:"weight"`
	QpsLimit   int `json:"qps_limit"`
	DailyLimit int `json:"daily_limit"`

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
		return pkggin.R{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, err
	}
	return pkggin.R{Code: http.StatusOK}, nil
}

type listProviderReq struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
}

type listProviderResp struct {
	Total   int64             `json:"total"`
	Content []domain.Provider `json:"content"`
}

func (h *ProviderHandler) List(ctx *gin.Context, req listProviderReq) (pkggin.R, error) {
	providers, err := h.svc.List(ctx)
	if err != nil {
		return pkggin.R{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, err
	}

	return pkggin.R{
		Code: http.StatusOK,
		Data: listProviderResp{
			Total:   int64(len(providers)),
			Content: providers,
		},
	}, nil
}

func NewProviderHandler(svc service.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		svc: svc,
	}
}
