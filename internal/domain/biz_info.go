package domain

// BizInfo 业务信息表
type BizInfo struct {
	Id           uint64 `json:"id"`
	BizKey       string `json:"biz_key"`
	BizSecret    string `json:"biz_secret"`
	BizName      string `json:"biz_name"`
	Contact      string `json:"contact"`
	ContactEmail string `json:"contact_email"`
	CreatorId    uint64 `json:"creator_id"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`

	Creator SysUser `json:"creator"`
}
