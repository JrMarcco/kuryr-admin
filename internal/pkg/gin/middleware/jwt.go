package middleware

import (
	"net/http"
	"strings"

	eawsyjwt "github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/easy-kit/set"
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/JrMarcco/kuryr-admin/internal/service/session"
	"github.com/gin-gonic/gin"
)

var _ Builder = (*JwtBuilder)(nil)

type JwtBuilder struct {
	sessionSvc session.Service
	atManager  eawsyjwt.Manager[ginpkg.AuthUser] // 这里是 access token manager
	ignores    set.Set[string]
}

func (b *JwtBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if b.ignores != nil && b.ignores.Exist(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		token := b.ExtractToken(ctx)
		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		decrypted, err := b.atManager.Decrypt(token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		au := decrypted.Data
		// 检查 session
		// 注意：
		//	这里是可选的，只依赖 RefreshToken 的检测通常就足够了。
		//err = b.sessionSvc.Check(ctx, au.Sid)
		//if err != nil {
		//	// 系统错误或用户已经退出登录
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		ctx.Set(ginpkg.ContextKeyAuthUser, au)
		ctx.Next()
	}
}

func (b *JwtBuilder) ExtractToken(ctx *gin.Context) string {
	token := ctx.GetHeader(ginpkg.HeaderNameAccessToken)
	if token != "" {
		return strings.TrimPrefix(token, "Bearer ")
	}
	return ""
}

func NewJwtBuilder(
	sessionSvc session.Service, atManager eawsyjwt.Manager[ginpkg.AuthUser], ignores set.Set[string],
) *JwtBuilder {
	return &JwtBuilder{
		sessionSvc: sessionSvc,
		atManager:  atManager,
		ignores:    ignores,
	}
}
