package repository

import (
	"context"

	"github.com/JrMarcco/kuryr-admin/internal/repository/cache"
)

type SessionRepo interface {
	Check(ctx context.Context, sid string) (bool, error)
	Create(ctx context.Context, sid string, uid uint64) error
	Refresh(ctx context.Context, sid string) error
	Clear(ctx context.Context, sid string) error
}

var _ SessionRepo = (*DefaultSessionRepo)(nil)

type DefaultSessionRepo struct {
	sessionCache cache.SessionCache
}

func (r *DefaultSessionRepo) Check(ctx context.Context, sid string) (bool, error) {
	return r.sessionCache.Exists(ctx, sid)
}

func (r *DefaultSessionRepo) Create(ctx context.Context, sid string, uid uint64) error {
	return r.sessionCache.Set(ctx, sid, uid)
}

func (r *DefaultSessionRepo) Refresh(ctx context.Context, sid string) error {
	return r.sessionCache.Refresh(ctx, sid)
}

func (r *DefaultSessionRepo) Clear(ctx context.Context, sid string) error {
	return r.sessionCache.Clear(ctx, sid)
}

func NewDefaultSessionRepo(sessionCache cache.SessionCache) *DefaultSessionRepo {
	return &DefaultSessionRepo{sessionCache: sessionCache}
}
