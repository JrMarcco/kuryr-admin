package gin

import (
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/gin-gonic/gin"
)

const (
	HeaderNameAccessToken = "x-access-token"
	ContextKeyAuthUser    = "auth-user"
)

// RouteRegistry 路由注册器。
// Handler 需要实现这个接口并在 RegisterRoutes 方法内注册路由。
type RouteRegistry interface {
	RegisterRoutes(engine *gin.Engine)
}

// R 接口统一返回
type R struct {
	Code int    `json:"code"` // 使用 http.StatusXxx
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type AuthUser struct {
	Sid      string          `json:"sid"`       // session id
	Bid      uint64          `json:"bid"`       // biz id
	Uid      uint64          `json:"uid"`       // user id
	UserType domain.UserType `json:"user_type"` // user type
}
