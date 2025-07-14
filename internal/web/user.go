package web

import (
	"net/http"

	"github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service"
	"github.com/gin-gonic/gin"
)

var _ ginpkg.Registry = (*UserHandler)(nil)

type UserHandler struct {
	jwtManager jwt.Manager[domain.AuthUser]
	userSvc    service.UserService
}

func (h *UserHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/user")

	v1.Handle(http.MethodPost, "/login", ginpkg.B[LoginReq](h.Login))
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(ctx *gin.Context, req LoginReq) (ginpkg.R, error) {
	user, err := h.userSvc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		}, err
	}

	au := domain.AuthUser{Id: user.Id}
	token, err := h.jwtManager.Encrypt(au)
	if err != nil {
		return ginpkg.R{
			Code: http.StatusUnauthorized,
			Msg:  err.Error(),
		}, err
	}
	return ginpkg.R{
		Code: http.StatusOK,
		Data: token,
	}, nil
}

func NewUserHandler(jwtManager jwt.Manager[domain.AuthUser], userSvc service.UserService) *UserHandler {
	return &UserHandler{
		jwtManager: jwtManager,
		userSvc:    userSvc,
	}
}
