package grpc

import (
	"time"
)

type Client struct {
	GrpcServerName string
	GrpcTimeout    time.Duration
}
