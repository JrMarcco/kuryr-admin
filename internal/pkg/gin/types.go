package gin

import "github.com/gin-gonic/gin"

type Registry interface {
	RegisterRoutes(engine *gin.Engine)
}

type R struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
