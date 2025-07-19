package jwt

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	ExtractAccessToken(ctx *gin.Context) string
	CreateSession(ctx *gin.Context, sid string, uid uint64) error
	CheckSession(ctx *gin.Context, sid string, uid uint64) error
	RefreshSession(ctx *gin.Context, sid string) error
	ClearSession(ctx *gin.Context, sid string) error
}
