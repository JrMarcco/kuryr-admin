package gin

import "github.com/gin-gonic/gin"

const (
	HeaderNameJwtToken     = "x-jwt-token"
	HeaderNameRefreshToken = "x-refresh-token"

	ParamNameAuthUser = "auth-user"
)

type Registry interface {
	RegisterRoutes(engine *gin.Engine)
}

type R struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
