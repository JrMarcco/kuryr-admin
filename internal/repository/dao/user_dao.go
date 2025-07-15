package dao

import (
	"context"

	"gorm.io/gorm"
)

type SysUser struct {
	Id        uint64
	Username  string
	Password  string
	Email     string
	UserType  string
	BizId     uint64
	CreatedAt int64
	UpdatedAt int64
}

func (su SysUser) TableName() string {
	return "sys_user"
}

type UserDAO interface {
	FindByUsername(ctx context.Context, username string) (SysUser, error)
}

var _ UserDAO = (*DefaultUserDAO)(nil)

type DefaultUserDAO struct {
	db *gorm.DB
}

func (d *DefaultUserDAO) FindByUsername(ctx context.Context, username string) (SysUser, error) {
	var su SysUser
	err := d.db.WithContext(ctx).Model(&SysUser{}).
		Where("username = ?", username).
		First(&su).Error
	return su, err
}

func NewUserDAO(db *gorm.DB) *DefaultUserDAO {
	return &DefaultUserDAO{db: db}
}
