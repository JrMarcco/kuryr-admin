package domain

// UserType 用户类型
type UserType string

const (
	UserTypeAdmin    UserType = "administrator"
	UserTypeOperator UserType = "operator"
)

type SysUser struct {
	Id        uint64
	Username  string
	Password  string
	Email     string
	UserType  UserType
	BizInfo   BizInfo
	CreatedAt int64
	UpdatedAt int64
}
