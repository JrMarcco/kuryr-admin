package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BizInfo struct {
	Id           uint64 `gorm:"column:id"`
	BizType      string `gorm:"column:biz_type"`
	BizKey       string `gorm:"column:biz_key"`
	BizSecret    string `gorm:"column:biz_secret"`
	BizName      string `gorm:"column:biz_name"`
	Contact      string `gorm:"column:contact"`
	ContactEmail string `gorm:"column:contact_email"`
	CreatorId    uint64 `gorm:"column:creator_id"`
	CreatedAt    int64  `gorm:"column:created_at"`
	UpdatedAt    int64  `gorm:"column:updated_at"`
}

func (BizInfo) TableName() string {
	return "biz_info"
}

type BizDao interface {
	SaveWithTx(ctx context.Context, tx *gorm.DB, entity BizInfo) (BizInfo, error)
	DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint64) error

	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, offset, limit int) ([]BizInfo, error)
	FindById(ctx context.Context, id uint64) (BizInfo, error)
}

var _ BizDao = (*DefaultBizDao)(nil)

type DefaultBizDao struct {
	db *gorm.DB
}

func (d *DefaultBizDao) SaveWithTx(ctx context.Context, tx *gorm.DB, entity BizInfo) (BizInfo, error) {
	now := time.Now().UnixMilli()
	entity.CreatedAt = now
	entity.UpdatedAt = now

	// 这里使用 upsert
	err := tx.WithContext(ctx).Model(&BizInfo{}).
		Clauses(clause.OnConflict{
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

func (d *DefaultBizDao) DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	return tx.WithContext(ctx).Model(&BizInfo{}).
		Where("id = ?", id).
		Delete(&BizInfo{}).Error
}

func (d *DefaultBizDao) Count(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&BizInfo{}).Count(&count).Error
	return count, err
}

func (d *DefaultBizDao) List(ctx context.Context, offset, limit int) ([]BizInfo, error) {
	var bis []BizInfo
	err := d.db.WithContext(ctx).Model(&BizInfo{}).
		Offset(offset).
		Limit(limit).
		Find(&bis).Error
	return bis, err
}

func (d *DefaultBizDao) FindById(ctx context.Context, id uint64) (BizInfo, error) {
	var bi BizInfo
	err := d.db.WithContext(ctx).Model(&BizInfo{}).
		Where("id = ?", id).
		First(&bi).Error
	return bi, err
}

func NewBizDAO(db *gorm.DB) *DefaultBizDao {
	return &DefaultBizDao{db: db}
}
