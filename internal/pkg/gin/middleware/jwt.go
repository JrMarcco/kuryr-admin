package middleware

import (
	"net/http"

	eawsyjwt "github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/easy-kit/set"
	pkggin "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	ijwt "github.com/JrMarcco/kuryr-admin/internal/web/jwt"
	"github.com/gin-gonic/gin"
)

var _ Builder = (*JwtBuilder)(nil)

type JwtBuilder struct {
	ijwt.Handler
	atManager eawsyjwt.Manager[pkggin.AuthUser] // 这里是 access token manager
	ignores   set.Set[string]
}

func (b *JwtBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if b.ignores != nil && b.ignores.Exist(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		token := b.ExtractAccessToken(ctx)
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

		ctx.Set(pkggin.ContextKeyAuthUser, au)
		ctx.Next()
	}
}

func NewJwtBuilder(
	handler ijwt.Handler, atManager eawsyjwt.Manager[pkggin.AuthUser], ignores set.Set[string],
) *JwtBuilder {
	return &JwtBuilder{
		Handler:   handler,
		atManager: atManager,
		ignores:   ignores,
	}
}
