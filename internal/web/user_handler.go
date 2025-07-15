package web

import (
	"log"
	"net/http"

	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/JrMarcco/kuryr-admin/internal/service/session"
	"github.com/gin-gonic/gin"
)

var _ ginpkg.RouteRegistry = (*UserHandler)(nil)

type UserHandler struct {
	userSvc    service.UserService
	sessionSvc session.Service
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
	au, err := h.userSvc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		}, err
	}

	// 创建 session
	err = h.sessionSvc.Create(ctx, au.Sid, au.Uid)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, err
	}

	at, st, err := h.userSvc.GenerateToken(ctx, au)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusInternalServerError,
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
	au, err := h.userSvc.VerifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		}, err
	}

	// 校验 session
	err = h.sessionSvc.Check(ctx, au.Sid)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		}, err
	}

	// 刷新 session 的过期时间
	err = h.sessionSvc.Create(ctx, au.Sid, au.Uid)
	if err != nil {
		// 刷新失败通常不应该中断整个流程，但需要记录日志
		log.Printf("failed to refresh session: %v", err)
	}

	// 重新生成 access token 和 refresh token
	at, st, err := h.userSvc.GenerateToken(ctx, au)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusInternalServerError,
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

func NewUserHandler(userSvc service.UserService, sessionSvc session.Service) *UserHandler {
	return &UserHandler{
		userSvc:    userSvc,
		sessionSvc: sessionSvc,
	}
}
