package vo

type LoginAppleRequest struct {
	Code string `json:"code"`
}

type LoginResponse struct {
	UserID    uint   `json:"user_id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

type LoginAnonRequest struct {
	DeviceID string `json:"device_id"`
	Type     string `json:"type"`     // 设备ID类型: idfv, idfa, android id, ...
	Platform string `json:"platform"` // 设备平台：ios, android
}

type LoginIDTokenRequest struct {
	IDToken  string `json:"id_token"`
	Type     string `json:"type"`     // IDToken类型: apple
	Platform string `json:"platform"` // 设备平台：ios, android
	Nonce    string `json:"nonce"`    // 随机数
}
