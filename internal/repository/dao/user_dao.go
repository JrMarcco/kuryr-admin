package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	FindById(ctx context.Context, id uint64) (SysUser, error)
	FindByBizId(ctx context.Context, bizId uint64) (SysUser, error)
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

	res := d.db.WithContext(ctx).Model(&SysUser{}).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"email":      u.Email,
			"password":   u.Password,
			"real_name":  u.RealName,
			"updated_at": now,
		}),
	}).Create(&u)
	if res.Error != nil {
		return SysUser{}, res.Error
	}
	return u, nil
}

func (d *DefaultUserDao) FindById(ctx context.Context, id uint64) (SysUser, error) {
	var su SysUser
	err := d.db.WithContext(ctx).Model(&SysUser{}).
		Where("id = ?", id).
		First(&su).Error
	return su, err
}

func (d *DefaultUserDao) FindByBizId(ctx context.Context, bizId uint64) (SysUser, error) {
	var su SysUser
	err := d.db.WithContext(ctx).Model(&SysUser{}).
		Where("biz_id = ?", bizId).
		First(&su).Error
	return su, err
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
