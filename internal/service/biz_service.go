package service

import (
	"context"

	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/repository"
)

type BizService interface {
	Count(ctx context.Context) (int64, error)

	List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error)
	FindById(ctx context.Context, id uint64) (domain.BizInfo, error)
}

var _ BizService = (*DefaultBizService)(nil)

type DefaultBizService struct {
	bizRepo repository.BizRepo
}

func (s *DefaultBizService) Count(ctx context.Context) (int64, error) {
	return s.bizRepo.Count(ctx)
}

func (s *DefaultBizService) List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error) {
	return s.bizRepo.List(ctx, offset, limit)
}

func (s *DefaultBizService) FindById(ctx context.Context, id uint64) (domain.BizInfo, error) {
	return s.bizRepo.FindById(ctx, id)
}

func NewBizService(bizRepo repository.BizRepo) *DefaultBizService {
	return &DefaultBizService{bizRepo: bizRepo}
}
