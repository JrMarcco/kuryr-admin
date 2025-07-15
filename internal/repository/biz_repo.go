package repository

import (
	"context"

	"github.com/JrMarcco/easy-kit/slice"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
)

type BizRepo interface {
	Count(ctx context.Context) (int64, error)

	List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error)
	FindById(ctx context.Context, id uint64) (domain.BizInfo, error)
}

var _ BizRepo = (*DefaultBizRepo)(nil)

type DefaultBizRepo struct {
	bizDAO dao.BizDAO
}

func (r *DefaultBizRepo) Count(ctx context.Context) (int64, error) {
	return r.bizDAO.Count(ctx)
}

func (r *DefaultBizRepo) List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error) {
	entities, err := r.bizDAO.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(entities, func(_ int, src dao.BizInfo) domain.BizInfo {
		return r.toDomain(src)
	}), nil
}

func (r *DefaultBizRepo) FindById(ctx context.Context, id uint64) (domain.BizInfo, error) {
	entity, err := r.bizDAO.FindById(ctx, id)
	if err != nil {
		return domain.BizInfo{}, err
	}
	return r.toDomain(entity), nil
}

func (r *DefaultBizRepo) toDomain(entity dao.BizInfo) domain.BizInfo {
	return domain.BizInfo{
		Id:        entity.Id,
		BizKey:    entity.BizKey,
		BizSecret: entity.BizSecret,
		BizName:   entity.BizName,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func NewBizRepo(bizDAO dao.BizDAO) *DefaultBizRepo {
	return &DefaultBizRepo{bizDAO: bizDAO}
}
