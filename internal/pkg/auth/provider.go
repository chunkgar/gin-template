package auth

import (
	"context"
	"time"
)

// TokenProvider 定义了身份提供商的接口
type TokenProvider interface {
	// GetName 返回提供商名称
	GetName() string

	// GetJWKSURL 返回 JWKS 端点 URL
	GetJWKSURL() string

	// GetIssuer 返回预期的 issuer
	GetIssuer() string

	// GetClientID 返回客户端 ID
	GetClientID() string

	// ValidateClaims 验证自定义 claims
	ValidateClaims(claims Claims, clientID string, nonce string) error

	// ParseClaims 解析 token 为特定的 claims 结构
	ParseClaims(tokenString string) (Claims, error)
}

// Claims 定义了通用的 token claims 接口
type Claims interface {
	GetSubject() string
	GetIssuer() string
	GetAudience() []string
	GetExpirationTime() time.Time
	GetIssuedAt() time.Time
	GetEmail() string
	GetEmailVerified() bool
	GetNonce() string
	ToMap() map[string]interface{}
}

// JWKSFetcher 定义了获取公钥的接口
type JWKSFetcher interface {
	FetchKeys(ctx context.Context, url string) (*JWKSet, error)
	GetCachedKeys(url string) (*JWKSet, bool)
	SetCache(url string, keys *JWKSet, ttl time.Duration)
}

// JWKSet 表示 JSON Web Key Set
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

// JWK 表示单个 JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// TokenVerifier 通用 token 验证器
type TokenVerifier interface {
	Verify(ctx context.Context, token string, nonce string) (Claims, TokenProvider, error)
	VerifyWithProvider(ctx context.Context, token string, clientID string, nonce string, provider TokenProvider) (Claims, error)
}
