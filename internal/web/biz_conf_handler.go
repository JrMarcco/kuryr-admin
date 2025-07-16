package web

import (
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/gin-gonic/gin"
)

var _ ginpkg.RouteRegistry = (*BizConfHandler)(nil)

type BizConfHandler struct{}

func (b *BizConfHandler) RegisterRoutes(engine *gin.Engine) {
	//TODO implement me
	panic("implement me")
}
