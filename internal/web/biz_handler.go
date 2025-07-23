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

	v1.Handle(http.MethodPost, "/create", pkggin.BU(h.Create))
	v1.Handle(http.MethodGet, "/list", pkggin.QU(h.List))
}

type createBizReq struct {
	BizKey       string `json:"biz_key"`
	BizName      string `json:"biz_name"`
	Contact      string `json:"contact"`
	ContactEmail string `json:"contact_email"`
}

// Create 新建业务方信息。
func (h *BizHandler) Create(ctx *gin.Context, req createBizReq, au pkggin.AuthUser) (pkggin.R, error) {
	bi := domain.BizInfo{
		BizKey:       req.BizKey,
		BizName:      req.BizName,
		Contact:      req.Contact,
		ContactEmail: req.ContactEmail,
		CreatorId:    au.Uid,
	}
	bi, err := h.svc.Create(ctx, bi)
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

type listBizReq struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
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
