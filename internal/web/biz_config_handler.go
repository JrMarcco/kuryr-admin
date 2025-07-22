package web

import (
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ pkggin.RouteRegistry = (*BizConfigHandler)(nil)

type BizConfigHandler struct {
	svc service.BizConfigService
}

func (b *BizConfigHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/biz_config")

	v1.Handle(http.MethodPost, "/create", pkggin.BU[saveReq](b.Create))
}

type saveReq struct {
}

func (b *BizConfigHandler) Create(ctx *gin.Context, req saveReq, au pkggin.AuthUser) (pkggin.R, error) {
	err := b.svc.Create(ctx, domain.BizConfig{})
	if err != nil {
		return pkggin.R{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
	}, nil
}

func NewBizConfigHandler(svc service.BizConfigService) *BizConfigHandler {
	return &BizConfigHandler{
		svc: svc,
	}
}
