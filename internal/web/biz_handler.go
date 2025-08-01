package web

import (
	"net/http"
	"strconv"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	pkggorm "github.com/JrMarcco/kuryr-admin/internal/pkg/gorm"
	"github.com/JrMarcco/kuryr-admin/internal/search"
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
	v1.Handle(http.MethodGet, "/search", pkggin.QU(h.Search))
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
		return pkggin.R{}, err
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
		return pkggin.R{}, err
	}
	return pkggin.R{Code: http.StatusOK}, nil
}

type searchBizReq struct {
	BizName string `json:"biz_name" form:"biz_name"`
	*pkggorm.PaginationParam
}

// Search 分页查询业务方信息
func (h *BizHandler) Search(ctx *gin.Context, req searchBizReq, au pkggin.AuthUser) (pkggin.R, error) {
	var res *pkggorm.PaginationResult[domain.BizInfo]
	var err error

	switch au.UserType {
	case domain.UserTypeAdmin:
		criteria := search.BizSearchCriteria{
			BizName: req.BizName,
		}
		res, err = h.svc.Search(ctx, criteria, req.PaginationParam)
	case domain.UserTypeOperator:
		var bizInfo domain.BizInfo
		bizInfo, err = h.svc.FindById(ctx, au.Bid)
		res = pkggorm.NewPaginationResult([]domain.BizInfo{bizInfo}, 1)
	default:
		return pkggin.R{}, errs.ErrUnknownUser
	}

	if err != nil {
		return pkggin.R{}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
		Data: res,
	}, nil
}

func NewBizHandler(svc service.BizService) *BizHandler {
	return &BizHandler{svc: svc}
}
