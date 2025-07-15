package web

import (
	"net/http"

	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ ginpkg.RouteRegistry = (*UserHandler)(nil)

type UserHandler struct {
	userSvc service.UserService
}

func (h *UserHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/user")

	v1.Handle(http.MethodPost, "/login", ginpkg.B[loginReq](h.Login))
	v1.Handle(http.MethodPost, "/refresh_token", ginpkg.B[refreshTokenReq](h.RefreshToken))
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) Login(ctx *gin.Context, req loginReq) (ginpkg.R, error) {
	at, st, err := h.userSvc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		}, err
	}
	return ginpkg.R{
		Code: http.StatusOK,
		Data: tokenResp{
			AccessToken:  at,
			RefreshToken: st,
		},
	}, nil
}

func (h *UserHandler) RefreshToken(ctx *gin.Context, req refreshTokenReq) (ginpkg.R, error) {
	at, st, err := h.userSvc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		}, err
	}
	return ginpkg.R{
		Code: http.StatusOK,
		Data: tokenResp{
			AccessToken:  at,
			RefreshToken: st,
		},
	}, nil
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}
