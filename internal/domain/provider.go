package domain

// Provider 供应商领域对象。
type Provider struct {
	Id           uint64 `json:"id"`
	ProviderName string `json:"provider_name"` // 供应商名称
	Channel      int32  `json:"channel"`       // 渠道

	Endpoint string `json:"endpoint"`  // 接口地址
	RegionId string `json:"region_id"` // 区域 ID

	AppId     string `json:"app_id"`     // 应用 ID
	ApiKey    string `json:"api_key"`    // 接口密钥
	ApiSecret string `json:"api_secret"` // 接口密钥

	Weight     int `json:"weight"`      // 权重
	QpsLimit   int `json:"qps_limit"`   // 每秒请求限制
	DailyLimit int `json:"daily_limit"` // 每日请求限制

	AuditCallbackUrl string `json:"audit_callback_url"` // 审核回调地址

	ActiveStatus string `json:"active_status"` // 状态
}
