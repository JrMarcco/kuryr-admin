package session

import (
	"context"
	"errors"

	"github.com/JrMarcco/kuryr-admin/internal/repository"
)

type Service interface {
	Check(ctx context.Context, sid string) error
	Create(ctx context.Context, sid string, uid uint64) error
	Refresh(ctx context.Context, sid string) error
	Clear(ctx context.Context, sid string) error
}

var _ Service = (*RedisSessionService)(nil)

type RedisSessionService struct {
	sessionRepo repository.SessionRepo
}

func (s *RedisSessionService) Check(ctx context.Context, sid string) error {
	res, err := s.sessionRepo.Check(ctx, sid)
	if err != nil {
		return err
	}
	if !res {
		return errors.New("has been logged out")
	}
	return nil
}

func (s *RedisSessionService) Create(ctx context.Context, sid string, uid uint64) error {
	return s.sessionRepo.Create(ctx, sid, uid)
}

func (s *RedisSessionService) Refresh(ctx context.Context, sid string) error {
	return s.sessionRepo.Refresh(ctx, sid)
}

func (s *RedisSessionService) Clear(ctx context.Context, sid string) error {
	return s.sessionRepo.Clear(ctx, sid)
}
func NewRedisSessionService(sessionRepo repository.SessionRepo) *RedisSessionService {
	return &RedisSessionService{sessionRepo: sessionRepo}
}
