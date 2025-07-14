package middleware

import (
	"github.com/gin-gonic/gin"
)

var _ Builder = (*AccessLogBuilder)(nil)

type AccessLogBuilder struct{}

func (b *AccessLogBuilder) Build() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func NewAccessLogBuilder() *AccessLogBuilder {
	return &AccessLogBuilder{}
}
