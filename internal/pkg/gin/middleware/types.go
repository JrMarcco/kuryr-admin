package middleware

import "github.com/gin-gonic/gin"

const (
	headerNameJwtToken     = "x-jwt-token"
	headerNameRefreshToken = "x-refresh-token"

	paramNameAuthUser = "auth-user"
)

type Builder interface {
	Build() gin.HandlerFunc
}
