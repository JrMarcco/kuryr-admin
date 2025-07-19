package repository

import (
	"context"
	"strings"

	"github.com/JrMarcco/easy-kit/slice"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
	"gorm.io/gorm"
)

type BizRepo interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, bi domain.BizInfo) (domain.BizInfo, error)

	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, offset, limit int) ([]domain.BizInfo, error)
	FindById(ctx context.Context, id uint64) (domain.BizInfo, error)
}

var _ BizRepo = (*DefaultBizRepo)(nil)

type DefaultBizRepo struct {
	bizDAO dao.BizDAO
}

func (r *DefaultBizRepo) CreateWithTx(ctx context.Context, tx *gorm.DB, bi domain.BizInfo) (domain.BizInfo, error) {
	entity, err := r.bizDAO.SaveWithTx(ctx, tx, r.toEntity(bi))
	if err != nil {
		if isUniqueConstraintError(err) {
			if strings.Contains(err.Error(), "biz_key") {
				return domain.BizInfo{}, errs.ErrBizKeyConflict
			}
		}
		return domain.BizInfo{}, err
	}
	return r.toDomain(entity), nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())

	// postgresql 唯一键冲突错误关键词
	postgresKeywords := []string{
		"unique constraint",
		"duplicate key",
		"violates unique constraint",
	}

	// MySQL 唯一键冲突错误关键词
	mysqlKeywords := []string{
		"duplicate entry",
		"unique constraint",
	}

	keywords := append(postgresKeywords, mysqlKeywords...)
	for _, keyword := range keywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
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
		Id:           entity.Id,
		BizKey:       entity.BizKey,
		BizSecret:    entity.BizSecret,
		BizName:      entity.BizName,
		Contact:      entity.Contact,
		ContactEmail: entity.ContactEmail,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

func (r *DefaultBizRepo) toEntity(bi domain.BizInfo) dao.BizInfo {
	return dao.BizInfo{
		Id:           bi.Id,
		BizKey:       bi.BizKey,
		BizSecret:    bi.BizSecret,
		BizName:      bi.BizName,
		Contact:      bi.Contact,
		ContactEmail: bi.ContactEmail,
		CreatedAt:    bi.CreatedAt,
		UpdatedAt:    bi.UpdatedAt,
	}
}

func NewBizRepo(bizDAO dao.BizDAO) *DefaultBizRepo {
	return &DefaultBizRepo{bizDAO: bizDAO}
}
