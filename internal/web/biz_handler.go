package web

import (
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ ginpkg.RouteRegistry = (*BizHandler)(nil)

type BizHandler struct {
	bizSvc service.BizService
}

func (h *BizHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/biz")

	v1.Handle(http.MethodGet, "/list", ginpkg.QU[listBizReq](h.List))
}

type listBizReq struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

type listBizResp struct {
	Total   int64            `json:"total"`
	Content []domain.BizInfo `json:"content"`
}

func (h *BizHandler) List(ctx *gin.Context, req listBizReq, au ginpkg.AuthUser) (ginpkg.R, error) {
	var (
		list  []domain.BizInfo
		total int64
		err   error
	)

	switch au.UserType {
	case domain.UserTypeAdmin:
		list, err = h.bizSvc.List(ctx, req.Offset, req.Limit)
		if err == nil {
			total, err = h.bizSvc.Count(ctx)
		}
	case domain.UserTypeOperator:
		var bizInfo domain.BizInfo
		bizInfo, err = h.bizSvc.FindById(ctx, au.Id)
		list = append(list, bizInfo)
		total = 1
	default:
		return ginpkg.R{}, errs.ErrUnknownUser
	}

	if err != nil {
		return ginpkg.R{}, err
	}
	return ginpkg.R{
		Code: http.StatusOK,
		Data: listBizResp{
			Total:   total,
			Content: list,
		},
	}, nil
}

func NewBizHandler(bizSvc service.BizService) *BizHandler {
	return &BizHandler{bizSvc: bizSvc}
}
