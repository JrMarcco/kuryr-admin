package repository

import (
	"context"
	"strings"

	"github.com/JrMarcco/easy-kit/slice"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	"github.com/JrMarcco/kuryr-admin/internal/errs"
	pkggorm "github.com/JrMarcco/kuryr-admin/internal/pkg/gorm"
	"github.com/JrMarcco/kuryr-admin/internal/repository/dao"
	"github.com/JrMarcco/kuryr-admin/internal/search"
	"gorm.io/gorm"
)

type BizRepo interface {
	SaveWithTx(ctx context.Context, tx *gorm.DB, bi domain.BizInfo) (domain.BizInfo, error)
	DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint64) error

	Search(ctx context.Context, criteria search.BizSearchCriteria, param *pkggorm.PaginationParam) (*pkggorm.PaginationResult[domain.BizInfo], error)
	FindById(ctx context.Context, id uint64) (domain.BizInfo, error)
}

var _ BizRepo = (*DefaultBizRepo)(nil)

type DefaultBizRepo struct {
	dao dao.BizDao
}

func (r *DefaultBizRepo) SaveWithTx(ctx context.Context, tx *gorm.DB, bi domain.BizInfo) (domain.BizInfo, error) {
	entity, err := r.dao.SaveWithTx(ctx, tx, r.toEntity(bi))
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

func (r *DefaultBizRepo) DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	return r.dao.DeleteWithTx(ctx, tx, id)
}

func (r *DefaultBizRepo) Search(
	ctx context.Context, criteria search.BizSearchCriteria, param *pkggorm.PaginationParam,
) (*pkggorm.PaginationResult[domain.BizInfo], error) {
	res, err := r.dao.Search(ctx, criteria, param)
	if err != nil {
		return nil, err
	}

	if res.Total == 0 {
		return pkggorm.NewPaginationResult([]domain.BizInfo{}, 0), nil
	}

	records := slice.Map(res.Records, func(idx int, src dao.BizInfo) domain.BizInfo {
		return r.toDomain(src)
	})
	return &pkggorm.PaginationResult[domain.BizInfo]{
		Total:   res.Total,
		Records: records,
	}, nil
}

func (r *DefaultBizRepo) FindById(ctx context.Context, id uint64) (domain.BizInfo, error) {
	entity, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.BizInfo{}, err
	}
	return r.toDomain(entity), nil
}

func (r *DefaultBizRepo) toDomain(entity dao.BizInfo) domain.BizInfo {
	return domain.BizInfo{
		Id:           entity.Id,
		BizType:      domain.BizType(entity.BizType),
		BizKey:       entity.BizKey,
		BizSecret:    entity.BizSecret[:3] + "****" + entity.BizSecret[len(entity.BizSecret)-3:],
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
		BizType:      string(bi.BizType),
		BizKey:       bi.BizKey,
		BizSecret:    bi.BizSecret,
		BizName:      bi.BizName,
		Contact:      bi.Contact,
		ContactEmail: bi.ContactEmail,
		CreatedAt:    bi.CreatedAt,
		UpdatedAt:    bi.UpdatedAt,
	}
}

func NewBizRepo(dao dao.BizDao) *DefaultBizRepo {
	return &DefaultBizRepo{dao: dao}
}
