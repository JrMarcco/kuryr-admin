package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SysUser struct {
	Id        uint64
	Email     string
	Password  string
	RealName  string
	UserType  string
	BizId     uint64
	CreatedAt int64
	UpdatedAt int64
}

func (su SysUser) TableName() string {
	return "sys_user"
}

type UserDao interface {
	SaveWithTx(ctx context.Context, tx *gorm.DB, u SysUser) (SysUser, error)
	DeleteByBizIdWithTx(ctx context.Context, tx *gorm.DB, id uint64) error

	FindByEmail(ctx context.Context, email string) (SysUser, error)
	FindByMobile(ctx context.Context, mobile string) (SysUser, error)
}

var _ UserDao = (*DefaultUserDao)(nil)

type DefaultUserDao struct {
	db *gorm.DB
}

func (d *DefaultUserDao) SaveWithTx(ctx context.Context, tx *gorm.DB, u SysUser) (SysUser, error) {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now

	err := tx.WithContext(ctx).Model(&SysUser{}).Create(&u).Error
	return u, err
}

func (d *DefaultUserDao) DeleteByBizIdWithTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	return tx.WithContext(ctx).Model(&BizInfo{}).
		Where("biz_id = ?", id).
		Delete(&SysUser{}).Error
}

func (d *DefaultUserDao) FindByEmail(ctx context.Context, email string) (SysUser, error) {
	var su SysUser
	err := d.db.WithContext(ctx).Model(&SysUser{}).
		Where("email = ?", email).
		First(&su).Error
	return su, err
}

func (d *DefaultUserDao) FindByMobile(ctx context.Context, mobile string) (SysUser, error) {
	var su SysUser
	err := d.db.WithContext(ctx).Model(&SysUser{}).
		Where("mobile = ?", mobile).
		First(&su).Error
	return su, err
}

func NewUserDAO(db *gorm.DB) *DefaultUserDao {
	return &DefaultUserDao{db: db}
}
