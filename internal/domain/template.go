package domain

type ChannelTemplate struct {
	Id      uint64  `json:"id"`
	BizId   uint64  `json:"biz_id"`
	BizType BizType `json:"biz_type"`

	TplName string `json:"tpl_name"`
	TplDesc string `json:"tpl_desc"`

	Channel            int32  `json:"channel"`
	NotificationType   int32  `json:"notification_type"`
	ActivatedVersionId uint64 `json:"activated_version_id"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`

	Versions []ChannelTemplateVersion `json:"versions"`
}

type ChannelTemplateVersion struct {
	Id    uint64 `json:"id"`
	TplId uint64 `json:"tpl_id"`

	VersionName string `json:"version_name"`
	Signature   string `json:"signature"`
	Content     string `json:"content"`
	ApplyRemark string `json:"apply_remark"`

	AuditId         uint64 `json:"audit_id"`
	AuditorId       uint64 `json:"auditor_id"`
	AuditTime       int64  `json:"audit_time"`
	AuditStatus     string `json:"audit_status"`
	RejectionReason string `json:"rejection_reason"`
	LastReviewAt    int64  `json:"last_review_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
