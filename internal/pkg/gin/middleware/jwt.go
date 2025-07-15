package middleware

import (
	"net/http"
	"strings"

	eawsyjwt "github.com/JrMarcco/easy-kit/jwt"
	"github.com/JrMarcco/easy-kit/set"
	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var _ Builder = (*JwtBuilder)(nil)

type JwtBuilder struct {
	rc        redis.Cmdable
	atManager eawsyjwt.Manager[ginpkg.AuthUser] // 这里是 access token manager
	ignores   set.Set[string]
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
		ctx.Set(ginpkg.ParamNameAuthUser, au)
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
	rc redis.Cmdable, atManager eawsyjwt.Manager[ginpkg.AuthUser], ignores set.Set[string],
) *JwtBuilder {
	return &JwtBuilder{
		rc:        rc,
		atManager: atManager,
		ignores:   ignores,
	}
}
