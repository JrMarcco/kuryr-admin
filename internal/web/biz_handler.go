package web

import (
	"net/http"
	"strconv"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ pkggin.RouteRegistry = (*BizHandler)(nil)

// BizHandler 业务方信息 web handler。
type BizHandler struct {
	svc service.BizService
}

func (h *BizHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/biz")

	v1.Handle(http.MethodPost, "/save", pkggin.BU(h.Save))
	v1.Handle(http.MethodDelete, "/delete", pkggin.QU(h.Delete))
	v1.Handle(http.MethodGet, "/list", pkggin.QU(h.List))
}

type createBizReq struct {
	BizType      string `json:"biz_type"`
	BizKey       string `json:"biz_key"`
	BizName      string `json:"biz_name"`
	Contact      string `json:"contact"`
	ContactEmail string `json:"contact_email"`
}

// Save 新建业务方信息。
func (h *BizHandler) Save(ctx *gin.Context, req createBizReq, au pkggin.AuthUser) (pkggin.R, error) {
	if au.UserType != domain.UserTypeAdmin {
		return pkggin.R{
			Code: http.StatusForbidden,
			Msg:  "[kuryr-admin] only admin can save biz",
		}, nil
	}

	bi := domain.BizInfo{
		BizType:      domain.BizType(req.BizType),
		BizKey:       req.BizKey,
		BizName:      req.BizName,
		Contact:      req.Contact,
		ContactEmail: req.ContactEmail,
		CreatorId:    au.Uid,
	}
	bi, err := h.svc.Save(ctx, bi)
	if err != nil {
		return pkggin.R{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
		Data: strconv.FormatUint(bi.Id, 10),
	}, nil
}

type deleteBizReq struct {
	BizId uint64 `json:"biz_id" form:"biz_id"`
}

func (h *BizHandler) Delete(ctx *gin.Context, req deleteBizReq, au pkggin.AuthUser) (pkggin.R, error) {
	if au.UserType != domain.UserTypeAdmin {
		return pkggin.R{
			Code: http.StatusForbidden,
			Msg:  "[kuryr-admin] only admin can delete biz",
		}, nil
	}
	err := h.svc.Delete(ctx, req.BizId)
	if err != nil {
		return pkggin.R{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, err
	}
	return pkggin.R{Code: http.StatusOK}, nil
}

type listBizReq struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
}

type listBizResp struct {
	Total   int64            `json:"total"`
	Content []domain.BizInfo `json:"content"`
}

// List 分页查询业务方信息
func (h *BizHandler) List(ctx *gin.Context, req listBizReq, au pkggin.AuthUser) (pkggin.R, error) {
	var (
		list  []domain.BizInfo
		total int64
		err   error
	)

	switch au.UserType {
	case domain.UserTypeAdmin:
		list, err = h.svc.List(ctx, req.Offset, req.Limit)
		if err == nil {
			total, err = h.svc.Count(ctx)
		}
	case domain.UserTypeOperator:
		var bizInfo domain.BizInfo
		bizInfo, err = h.svc.FindById(ctx, au.Bid)
		list = append(list, bizInfo)
		total = 1
	default:
		return pkggin.R{}, errs.ErrUnknownUser
	}

	if err != nil {
		return pkggin.R{}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
		Data: listBizResp{
			Total:   total,
			Content: list,
		},
	}, nil
}

func NewBizHandler(svc service.BizService) *BizHandler {
	return &BizHandler{svc: svc}
}
