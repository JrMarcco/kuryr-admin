package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/kuryr-admin/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

var _ cache.SessionCache = (*RSessionCache)(nil)

type RSessionCache struct {
	rc         redis.Cmdable
	expiration time.Duration
}

func (c *RSessionCache) Set(ctx context.Context, sid string, uid uint64) error {
	return c.rc.Set(ctx, c.key(sid), uid, c.expiration).Err()
}

func (c *RSessionCache) Exists(ctx context.Context, sid string) (bool, error) {
	res, err := c.rc.Exists(ctx, c.key(sid)).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (c *RSessionCache) Refresh(ctx context.Context, sid string) error {
	return c.rc.Expire(ctx, c.key(sid), c.expiration).Err()

}

func (c *RSessionCache) Clear(ctx context.Context, sid string) error {
	return c.rc.Del(ctx, c.key(sid)).Err()
}

func (c *RSessionCache) key(sid string) string {
	return fmt.Sprintf("user:sid:%s", sid)
}

func NewRSessionCache(rc redis.Cmdable, expiration time.Duration) *RSessionCache {
	return &RSessionCache{
		rc:         rc,
		expiration: expiration,
	}
}
