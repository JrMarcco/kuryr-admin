package domain

// BizInfo 业务信息表
type BizInfo struct {
	Id        uint64
	BizKey    string
	BizSecret string
	BizName   string
	CreatedAt int64
	UpdatedAt int64
}
