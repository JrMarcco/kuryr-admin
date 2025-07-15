package cache

import "context"

type SessionCache interface {
	Set(ctx context.Context, sid string, uid uint64) error
	Exists(ctx context.Context, sid string) (bool, error)
	Refresh(ctx context.Context, sid string) error
}
