package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SysUser struct {
	Id        uint64 `gorm:"column:id"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	RealName  string `gorm:"column:real_name"`
	UserType  string `gorm:"column:user_type"`
	BizId     uint64 `gorm:"column:biz_id"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

type UserDao interface {
	Save(ctx context.Context, u SysUser) (SysUser, error)
	DeleteByBizId(ctx context.Context, id uint64) error

	FindByEmail(ctx context.Context, email string) (SysUser, error)
	FindByMobile(ctx context.Context, mobile string) (SysUser, error)
}

var _ UserDao = (*DefaultUserDao)(nil)

type DefaultUserDao struct {
	db *gorm.DB
}

func (d *DefaultUserDao) Save(ctx context.Context, u SysUser) (SysUser, error) {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now

	err := d.db.WithContext(ctx).Model(&SysUser{}).Create(&u).Error
	return u, err
}

func (d *DefaultUserDao) DeleteByBizId(ctx context.Context, id uint64) error {
	return d.db.WithContext(ctx).Model(&SysUser{}).
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
