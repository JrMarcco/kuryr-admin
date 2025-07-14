package ioc

import (
	"github.com/JrMarcco/kuryr-admin/internal/pkg/gin/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var GinFxOpt = fx.Provide(
	InitMiddlewares,
	InitGin,
)

// InitMiddlewares 提供一个用于创建有序中间件切片的函数
func InitMiddlewares(
	corsBuilder *middleware.CorsBuilder,
	jwtBuilder *middleware.JwtBuilder,
) []middleware.Builder {
	// 按顺序排列中间件
	return []middleware.Builder{
		corsBuilder,
		jwtBuilder,
	}
}

func InitGin(mbs []middleware.Builder) *gin.Engine {
	engine := gin.Default()

	// 注册中间件
	middlewares := make([]gin.HandlerFunc, 0, len(mbs))
	for _, mb := range mbs {
		middlewares = append(middlewares, mb.Build())
	}
	engine.Use(middlewares...)

	return engine
}
