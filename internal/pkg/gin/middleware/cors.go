package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var _ Builder = (*CorsBuilder)(nil)

type CorsBuilder struct {
	allowCredentials bool
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	maxAge           time.Duration

	allowOriginFunc func(origin string) bool
}

func (b *CorsBuilder) Build() gin.HandlerFunc {
	corsCfg := cors.Config{
		AllowCredentials: b.allowCredentials,
		AllowMethods:     b.allowMethods,
		AllowHeaders:     b.allowHeaders,
		ExposeHeaders:    b.exposeHeaders,
		MaxAge:           b.maxAge,

		AllowOriginFunc: b.allowOriginFunc,
	}

	return cors.New(corsCfg)
}

func (b *CorsBuilder) AllowCredentials(allowCredentials bool) *CorsBuilder {
	b.allowCredentials = allowCredentials
	return b
}

func (b *CorsBuilder) AllowMethods(allowMethods []string) *CorsBuilder {
	b.allowMethods = allowMethods
	return b
}

func (b *CorsBuilder) AllowHeaders(allowHeaders []string) *CorsBuilder {
	b.allowHeaders = allowHeaders
	return b
}

func (b *CorsBuilder) ExposeHeaders(exposeHeaders []string) *CorsBuilder {
	b.exposeHeaders = exposeHeaders
	return b
}

func (b *CorsBuilder) MaxAge(maxAge time.Duration) *CorsBuilder {
	b.maxAge = maxAge
	return b
}

func (b *CorsBuilder) AllowOriginFunc(allowOriginFunc func(origin string) bool) *CorsBuilder {
	b.allowOriginFunc = allowOriginFunc
	return b
}

func NewCorsBuilder() *CorsBuilder {
	return &CorsBuilder{
		allowCredentials: true,
		allowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		allowHeaders:     []string{"Content-Length", "Content-Type", "Authorization"},
		exposeHeaders:    []string{headerNameJwtToken, headerNameRefreshToken},
		maxAge:           12 * time.Hour,
		allowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
	}
}
