package provider

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/chunkgar/gin-template/internal/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
)

// AppleProvider Apple 身份提供商
type AppleProvider struct {
	name     string
	jwksURL  string
	issuer   string
	clientID string
}

// NewAppleProvider 创建 Apple 提供商
func NewAppleProvider(clientID string) *AppleProvider {
	return &AppleProvider{
		name:     "apple",
		jwksURL:  "https://appleid.apple.com/auth/keys",
		issuer:   "https://appleid.apple.com",
		clientID: clientID,
	}
}

func (p *AppleProvider) GetName() string {
	return p.name
}

func (p *AppleProvider) GetJWKSURL() string {
	return p.jwksURL
}

func (p *AppleProvider) GetIssuer() string {
	return p.issuer
}

func (p *AppleProvider) GetClientID() string {
	return p.clientID
}

// AppleClaims Apple ID Token claims
type AppleClaims struct {
	Issuer         string `json:"iss"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	Nonce          string `json:"nonce"`
	IssuedAt       int64  `json:"iat"`
	ExpirationTime int64  `json:"exp"`
	Email          string `json:"email,omitempty"`
	EmailVerified  bool   `json:"email_verified,omitempty"`
	IsPrivateEmail bool   `json:"is_private_email,omitempty"`
	RealUserStatus int    `json:"real_user_status,omitempty"`
}

func (c *AppleClaims) GetSubject() string {
	return c.Subject
}

func (c *AppleClaims) GetIssuer() string {
	return c.Issuer
}

func (c *AppleClaims) GetAudience() []string {
	return []string{c.Audience}
}

func (c *AppleClaims) GetExpirationTime() time.Time {
	return time.Unix(c.ExpirationTime, 0)
}

func (c *AppleClaims) GetIssuedAt() time.Time {
	return time.Unix(c.IssuedAt, 0)
}

func (c *AppleClaims) GetEmail() string {
	return c.Email
}

func (c *AppleClaims) GetEmailVerified() bool {
	return c.EmailVerified
}

func (c *AppleClaims) GetNonce() string {
	return c.Nonce
}

func (c *AppleClaims) ToMap() map[string]interface{} {
	data, _ := json.Marshal(c)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

func (p *AppleProvider) ParseClaims(tokenString string) (auth.Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	claimsJSON, err := json.Marshal(claimsMap)
	if err != nil {
		return nil, err
	}

	var claims AppleClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (p *AppleProvider) ValidateClaims(claims auth.Claims, clientID string, nonce string) error {
	// 验证 nonce
	sha256Nonce := fmt.Sprintf("%x", sha256.Sum256([]byte(nonce)))
	if claims.GetNonce() != sha256Nonce {
		return fmt.Errorf("invalid nonce: expected %s, got %s", sha256Nonce, claims.GetNonce())
	}

	// 验证 issuer
	if claims.GetIssuer() != p.issuer {
		return fmt.Errorf("invalid issuer: expected %s, got %s", p.issuer, claims.GetIssuer())
	}

	// 验证 audience
	audiences := claims.GetAudience()
	validAudience := false
	for _, aud := range audiences {
		if aud == clientID {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return fmt.Errorf("invalid audience: expected %s", clientID)
	}

	// 验证过期时间
	if time.Now().After(claims.GetExpirationTime()) {
		return errors.New("token has expired")
	}

	// 验证签发时间（不能是未来时间）
	if time.Now().Before(claims.GetIssuedAt()) {
		return errors.New("token issued in the future")
	}

	return nil
}
