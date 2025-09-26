package web

import (
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ pkggin.RouteRegistry = (*TemplateHandler)(nil)

type TemplateHandler struct {
	svc service.TemplateService
}

func (h *TemplateHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/template")

	v1.Handle(http.MethodPost, "/save", pkggin.B(h.Save))
}

type saveTemplateReq struct {
	BizId            uint64 `json:"biz_id"`
	BizType          string `json:"biz_type"`
	TplName          string `json:"tpl_name"`
	TplDesc          string `json:"tpl_desc"`
	Channel          int32  `json:"channel"`
	NotificationType int32  `json:"notification_type"`
}

func (h *TemplateHandler) Save(ctx *gin.Context, req saveTemplateReq) (pkggin.R, error) {
	tpl := domain.ChannelTemplate{
		BizId:            req.BizId,
		BizType:          domain.BizType(req.BizType),
		TplName:          req.TplName,
		TplDesc:          req.TplDesc,
		Channel:          req.Channel,
		NotificationType: req.NotificationType,
	}

	_, err := h.svc.Save(ctx, tpl)
	if err != nil {
		return pkggin.R{}, err
	}

	return pkggin.R{Code: http.StatusOK}, nil
}

type saveVersionReq struct {
	TplId       uint64 `json:"tpl_id"`
	VersionName string `json:"version_name"`
	Signature   string `json:"signature"`
	Content     string `json:"content"`
	ApplyRemark string `json:"apply_remark"`
}

func (h *TemplateHandler) SaveVersion(ctx *gin.Context, req saveVersionReq) (pkggin.R, error) {
	version := domain.ChannelTemplateVersion{
		TplId:       req.TplId,
		VersionName: req.VersionName,
		Signature:   req.Signature,
		Content:     req.Content,
		ApplyRemark: req.ApplyRemark,
	}

	_, err := h.svc.SaveVersion(ctx, version)
	if err != nil {
		return pkggin.R{}, err
	}

	return pkggin.R{Code: http.StatusOK}, nil
}

func NewTemplateHandler(svc service.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		svc: svc,
	}
}
