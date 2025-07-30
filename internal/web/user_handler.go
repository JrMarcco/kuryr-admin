package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/errs"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/JrMarcco/kuryr-admin/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ pkggin.RouteRegistry = (*UserHandler)(nil)

type UserHandler struct {
	jwt.Handler
	svc    service.UserService
	logger *zap.Logger
}

func (h *UserHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/user")

	v1.Handle(http.MethodPost, "/login", pkggin.B(h.Login))
	v1.Handle(http.MethodPost, "/refresh_token", pkggin.B(h.RefreshToken))
	v1.Handle(http.MethodGet, "/logout", pkggin.W(h.Logout))
}

type loginReq struct {
	Account     string `json:"account"`
	AccountType string `json:"account_type"`
	Credential  string `json:"credential"`
	VerifyType  string `json:"verify_type"`
}

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) Login(ctx *gin.Context, req loginReq) (pkggin.R, error) {
	au, err := h.svc.LoginWithType(ctx, req.Account, req.Credential, req.AccountType, req.VerifyType)
	if err != nil {
		return pkggin.R{}, err
	}

	// 创建 session
	err = h.CreateSession(ctx, au.Sid, au.Uid)
	if err != nil {
		return pkggin.R{}, err
	}

	at, st, err := h.svc.GenerateToken(ctx, au)
	if err != nil {
		return pkggin.R{}, err
	}

	return pkggin.R{
		Code: http.StatusOK,
		Data: tokenResp{
			AccessToken:  at,
			RefreshToken: st,
		},
	}, nil
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) RefreshToken(ctx *gin.Context, req refreshTokenReq) (pkggin.R, error) {
	au, err := h.svc.VerifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return pkggin.R{}, fmt.Errorf("%w: %s", errs.ErrUnauthorized, err.Error())
	}

	// 校验 session
	err = h.CheckSession(ctx, au.Sid, au.Uid)
	if err != nil {
		return pkggin.R{}, fmt.Errorf("%w: %s", errs.ErrUnauthorized, err.Error())
	}

	// 刷新 session 的过期时间
	err = h.RefreshSession(ctx, au.Sid)
	if err != nil {
		// 刷新失败通常不应该中断整个流程，但需要记录日志
		log.Printf("failed to refresh session: %v", err)
	}

	// 重新生成 access token 和 refresh token
	at, st, err := h.svc.GenerateToken(ctx, au)
	if err != nil {
		return pkggin.R{}, err
	}
	return pkggin.R{
		Code: http.StatusOK,
		Data: tokenResp{
			AccessToken:  at,
			RefreshToken: st,
		},
	}, nil
}

func (h *UserHandler) Logout(ctx *gin.Context) (pkggin.R, error) {
	userVal, ok := ctx.Get(pkggin.ContextKeyAuthUser)
	if !ok {
		return pkggin.R{
			Code: http.StatusUnauthorized,
			Msg:  "user not logged in",
		}, nil
	}

	au, ok := userVal.(pkggin.AuthUser)
	if !ok {
		return pkggin.R{
			Code: http.StatusUnauthorized,
			Msg:  "user not logged in",
		}, nil
	}

	go func() {
		if err := h.ClearSession(ctx, au.Sid); err != nil {
			h.logger.Error("failed to clear session", zap.Error(err))
		}
	}()
	return pkggin.R{
		Code: http.StatusOK,
		Msg:  "logged out",
	}, nil
}

func NewUserHandler(handler jwt.Handler, svc service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		Handler: handler,
		svc:     svc,
		logger:  logger,
	}
}
