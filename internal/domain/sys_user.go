package domain

// UserType 用户类型
type UserType string

const (
	UserTypeAdmin    UserType = "administrator"
	UserTypeOperator UserType = "operator"
)

type SysUser struct {
	Id        uint64
	Email     string
	Password  string
	RealName  string
	UserType  UserType
	BizId     uint64
	CreatedAt int64
	UpdatedAt int64
}
