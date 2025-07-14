package middleware

import "github.com/gin-gonic/gin"

type Builder interface {
	Build() gin.HandlerFunc
}
