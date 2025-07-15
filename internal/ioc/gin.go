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
	//
	// 注意：
	//	这里需要保证 fx 在 ioc.RegisterRoutes 之前完成。
	//	ioc.RegisterRoutes 在 engine.Use(middlewares...) 之前调用，
	//	会导致这里注册的路由“错过”这里注册的 middleware ，即导致 middleware 失效。
	middlewares := make([]gin.HandlerFunc, 0, len(mbs))
	for _, mb := range mbs {
		middlewares = append(middlewares, mb.Build())
	}
	engine.Use(middlewares...)

	return engine
}
