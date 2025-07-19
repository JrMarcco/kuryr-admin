package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BizInfo struct {
	Id           uint64
	BizKey       string
	BizSecret    string
	BizName      string
	Contact      string
	ContactEmail string
	CreatorId    uint64
	CreatedAt    int64
	UpdatedAt    int64
}

func (bi BizInfo) TableName() string {
	return "biz_info"
}

type BizDAO interface {
	SaveWithTx(ctx context.Context, tx *gorm.DB, entity BizInfo) (BizInfo, error)

	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, offset, limit int) ([]BizInfo, error)
	FindById(ctx context.Context, id uint64) (BizInfo, error)
}

var _ BizDAO = (*DefaultBizDAO)(nil)

type DefaultBizDAO struct {
	db *gorm.DB
}

func (d *DefaultBizDAO) SaveWithTx(ctx context.Context, tx *gorm.DB, entity BizInfo) (BizInfo, error) {
	now := time.Now().UnixMilli()
	entity.CreatedAt = now
	entity.UpdatedAt = now

	// 这里使用 upsert
	err := tx.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"biz_key":    entity.BizKey,
			"biz_name":   entity.BizName,
			"updated_at": now,
		}),
	}).Create(&entity).Error
	if err != nil {
		return BizInfo{}, err
	}
	return entity, nil
}

func (d *DefaultBizDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&BizInfo{}).Count(&count).Error
	return count, err
}

func (d *DefaultBizDAO) List(ctx context.Context, offset, limit int) ([]BizInfo, error) {
	var bis []BizInfo
	err := d.db.WithContext(ctx).Model(&BizInfo{}).
		Offset(offset).
		Limit(limit).
		Find(&bis).Error
	return bis, err
}

func (d *DefaultBizDAO) FindById(ctx context.Context, id uint64) (BizInfo, error) {
	var bi BizInfo
	err := d.db.WithContext(ctx).Model(&BizInfo{}).
		Where("id = ?", id).
		First(&bi).Error
	return bi, err
}

func NewBizDAO(db *gorm.DB) *DefaultBizDAO {
	return &DefaultBizDAO{db: db}
}
