package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	ginpkg "github.com/JrMarcco/kuryr-admin/internal/pkg/gin"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var _ Handler = (*RedisHandler)(nil)

type RedisHandler struct {
	rc         redis.Cmdable
	expiration time.Duration
}

func (h *RedisHandler) ExtractAccessToken(ctx *gin.Context) string {
	token := ctx.GetHeader(ginpkg.HeaderNameAccessToken)
	if token != "" {
		return strings.TrimPrefix(token, "Bearer ")
	}
	return ""
}

func (h *RedisHandler) CreateSession(ctx *gin.Context, sid string, uid uint64) error {
	return h.rc.Set(ctx, h.redisKey(sid), uid, h.expiration).Err()
}

func (h *RedisHandler) CheckSession(ctx *gin.Context, sid string, uid uint64) error {
	storedUid, err := h.rc.Get(ctx, h.redisKey(sid)).Uint64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("has been logged out")
		}
		return err
	}
	if storedUid != uid {
		return errors.New("session user mismatch")
	}
	return nil
}

func (h *RedisHandler) RefreshSession(ctx *gin.Context, sid string) error {
	return h.rc.Expire(ctx, h.redisKey(sid), h.expiration).Err()
}

func (h *RedisHandler) ClearSession(ctx *gin.Context, sid string) error {
	return h.rc.Del(ctx, h.redisKey(sid)).Err()
}

func (h *RedisHandler) redisKey(sid string) string {
	return fmt.Sprintf("user:sid:%s", sid)
}

func NewRedisHandler(rc redis.Cmdable, expiration time.Duration) *RedisHandler {
	return &RedisHandler{
		rc:         rc,
		expiration: expiration,
	}
}
