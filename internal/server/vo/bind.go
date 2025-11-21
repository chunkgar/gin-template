package vo

type BindIDTokenRequest struct {
	IDToken  string `json:"id_token"`
	Type     string `json:"type"`     // IDToken类型: apple
	Platform string `json:"platform"` // 设备平台：ios, android
	Nonce    string `json:"nonce"`    // 随机数
}

type UnbindRequest struct {
	Type string `json:"type"` // IDToken类型: apple
}
