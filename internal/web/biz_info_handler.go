package web

import (
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	pkggorm "github.com/JrMarcco/kuryr-admin/internal/pkg/gorm"
	"github.com/JrMarcco/kuryr-admin/internal/search"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ pkggin.RouteRegistry = (*BizInfoHandler)(nil)

// BizInfoHandler 业务方信息 web handler。
type BizInfoHandler struct {
	svc service.BizService
}

func (h *BizInfoHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/biz_info")

	v1.Handle(http.MethodPost, "/save", pkggin.BU(h.Save))
	v1.Handle(http.MethodPut, "/update", pkggin.BU(h.Update))
	v1.Handle(http.MethodDelete, "/delete", pkggin.QU(h.Delete))
	v1.Handle(http.MethodGet, "/search", pkggin.QU(h.Search))
	v1.Handle(http.MethodGet, "/get", pkggin.QU(h.FindById))
}

type createBizReq struct {
	BizType      string `json:"biz_type"`
	BizKey       string `json:"biz_key"`
	BizName      string `json:"biz_name"`
	Contact      string `json:"contact"`
	ContactEmail string `json:"contact_email"`
}

// Save 新建业务方信息。
func (h *BizInfoHandler) Save(ctx *gin.Context, req createBizReq, au pkggin.AuthUser) (pkggin.R, error) {
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
		Data: bi,
	}, nil
}

type updateBizReq struct {
	Id      uint64 `json:"id"`
	BizName string `json:"biz_name"`
}

func (h *BizInfoHandler) Update(ctx *gin.Context, req updateBizReq, au pkggin.AuthUser) (pkggin.R, error) {
	if au.UserType != domain.UserTypeAdmin {
		return pkggin.R{
			Code: http.StatusForbidden,
			Msg:  "[kuryr-admin] only admin can update biz",
		}, nil
	}

	bi := domain.BizInfo{
		Id:      req.Id,
		BizName: req.BizName,
	}
	bi, err := h.svc.Update(ctx, bi)
	if err != nil {
		return pkggin.R{}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
		Data: bi,
	}, nil
}

type deleteBizReq struct {
	BizId uint64 `json:"biz_id" form:"biz_id"`
}

func (h *BizInfoHandler) Delete(ctx *gin.Context, req deleteBizReq, au pkggin.AuthUser) (pkggin.R, error) {
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
func (h *BizInfoHandler) Search(ctx *gin.Context, req searchBizReq, au pkggin.AuthUser) (pkggin.R, error) {
	var res *pkggorm.PaginationResult[domain.BizInfo]
	var err error

	criteria := search.BizSearchCriteria{
		BizName: req.BizName,
	}

	switch au.UserType {
	case domain.UserTypeAdmin:
		// do nothing
	case domain.UserTypeOperator:
		criteria.BizId = au.Bid
	default:
		return pkggin.R{}, errs.ErrUnknownUserType
	}

	res, err = h.svc.Search(ctx, criteria, req.PaginationParam)
	if err != nil {
		return pkggin.R{}, err
	}

	return pkggin.R{
		Code: http.StatusOK,
		Data: res,
	}, nil
}

type findBizByIdReq struct {
	BizId uint64 `json:"biz_id" form:"biz_id"`
}

func (h *BizInfoHandler) FindById(ctx *gin.Context, req findBizByIdReq, au pkggin.AuthUser) (pkggin.R, error) {
	if au.UserType != domain.UserTypeAdmin && req.BizId != au.Bid {
		return pkggin.R{
			Code: http.StatusForbidden,
			Msg:  "[kuryr-admin] operator can only find own biz info",
		}, nil
	}

	bi, err := h.svc.FindById(ctx, req.BizId)
	if err != nil {
		return pkggin.R{}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
		Data: bi,
	}, nil
}

func NewBizHandler(svc service.BizService) *BizInfoHandler {
	return &BizInfoHandler{svc: svc}
}
