package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// CachedJWKSFetcher 带缓存的 JWKS 获取器
type CachedJWKSFetcher struct {
	cache      map[string]*cachedJWKS
	mu         sync.RWMutex
	httpClient *http.Client
	defaultTTL time.Duration
}

type cachedJWKS struct {
	keys      *JWKSet
	expiresAt time.Time
}

// NewCachedJWKSFetcher 创建新的 JWKS 获取器
func NewCachedJWKSFetcher(defaultTTL time.Duration) *CachedJWKSFetcher {
	return &CachedJWKSFetcher{
		cache: make(map[string]*cachedJWKS),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		defaultTTL: defaultTTL,
	}
}

// FetchKeys 获取公钥（带缓存）
func (f *CachedJWKSFetcher) FetchKeys(ctx context.Context, url string) (*JWKSet, error) {
	// 检查缓存
	if keys, ok := f.GetCachedKeys(url); ok {
		return keys, nil
	}

	// 从远程获取
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var jwks JWKSet
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	// 缓存结果
	f.SetCache(url, &jwks, f.defaultTTL)

	return &jwks, nil
}

// GetCachedKeys 从缓存获取公钥
func (f *CachedJWKSFetcher) GetCachedKeys(url string) (*JWKSet, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	cached, ok := f.cache[url]
	if !ok || time.Now().After(cached.expiresAt) {
		return nil, false
	}

	return cached.keys, true
}

// SetCache 设置缓存
func (f *CachedJWKSFetcher) SetCache(url string, keys *JWKSet, ttl time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.cache[url] = &cachedJWKS{
		keys:      keys,
		expiresAt: time.Now().Add(ttl),
	}
}

// ClearCache 清除缓存
func (f *CachedJWKSFetcher) ClearCache() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cache = make(map[string]*cachedJWKS)
}
