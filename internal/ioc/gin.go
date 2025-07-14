package ioc

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var GinFxOpt = fx.Provide(InitGin)

func InitGin() *gin.Engine {
	return gin.Default()
}
