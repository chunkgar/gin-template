package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

// DefaultTokenVerifier 默认的 token 验证器实现
type DefaultTokenVerifier struct {
	fetcher   JWKSFetcher
	providers map[string]TokenProvider
	mu        sync.RWMutex
}

// NewTokenVerifier 创建新的验证器
func NewTokenVerifier(fetcher JWKSFetcher) *DefaultTokenVerifier {
	return &DefaultTokenVerifier{
		fetcher:   fetcher,
		providers: make(map[string]TokenProvider),
	}
}

// RegisterProvider 注册身份提供商
func (v *DefaultTokenVerifier) RegisterProvider(provider TokenProvider) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.providers[provider.GetName()] = provider
}

// GetProvider 获取已注册的提供商
func (v *DefaultTokenVerifier) GetProvider(name string) (TokenProvider, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	provider, ok := v.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not registered", name)
	}
	return provider, nil
}

// Verify 验证 token（自动检测提供商）
func (v *DefaultTokenVerifier) Verify(ctx context.Context, token string, nonce string) (Claims, TokenProvider, error) {
	// 解析 token 获取 issuer
	unverified, err := v.parseUnverified(token)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse token: %w", err)
	}

	issuer := unverified["iss"].(string)

	// 查找匹配的提供商
	v.mu.RLock()
	var provider TokenProvider
	for _, p := range v.providers {
		if p.GetIssuer() == issuer {
			provider = p
			break
		}
	}
	v.mu.RUnlock()

	if provider == nil {
		return nil, nil, fmt.Errorf("no provider found for issuer: %s", issuer)
	}

	claims, err := v.VerifyWithProvider(ctx, token, provider.GetClientID(), nonce, provider)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify with provider: %w", err)
	}
	return claims, provider, nil
}

// VerifyWithProvider 使用指定的提供商验证 token
func (v *DefaultTokenVerifier) VerifyWithProvider(ctx context.Context, token string, clientID string, nonce string, provider TokenProvider) (Claims, error) {
	// 1. 解析 token header 获取 kid
	unverifiedToken, err := jwt.Parse(token, nil)
	if err != nil && !errors.Is(err, jwt.ErrTokenUnverifiable) {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	kid, ok := unverifiedToken.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid not found in token header")
	}

	// 2. 获取 JWKS
	jwksURL := provider.GetJWKSURL()
	jwks, err := v.fetcher.FetchKeys(ctx, jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	// 3. 找到匹配的公钥
	var matchedKey *JWK
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			matchedKey = &key
			break
		}
	}

	if matchedKey == nil {
		return nil, errors.New("no matching key found")
	}

	// 4. 转换为 RSA 公钥
	publicKey, err := jwkToRSAPublicKey(matchedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JWK to RSA key: %w", err)
	}

	// 5. 验证签名
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is invalid")
	}

	// 6. 解析 claims
	claims, err := provider.ParseClaims(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// 7. 验证 claims
	if err := provider.ValidateClaims(claims, clientID, nonce); err != nil {
		return nil, fmt.Errorf("claims validation failed: %w", err)
	}

	return claims, nil
}

// parseUnverified 解析未验证的 token
func (v *DefaultTokenVerifier) parseUnverified(token string) (map[string]interface{}, error) {
	parser := jwt.NewParser()
	unverifiedToken, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := unverifiedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	return claims, nil
}

// jwkToRSAPublicKey 将 JWK 转换为 RSA 公钥
func jwkToRSAPublicKey(jwk *JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)

	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	return &rsa.PublicKey{N: n, E: eInt}, nil
}
