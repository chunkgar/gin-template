package model

import "time"

// UserOAuthToken 第三方OAuth Token表（可选）
type UserOAuthToken struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AccountID    uint       `gorm:"not null" json:"account_id"`
	AccessToken  string     `gorm:"type:text" json:"-"`
	RefreshToken string     `gorm:"type:text" json:"-"`
	ExpiresAt    *time.Time `json:"expires_at"`
	Scope        string     `gorm:"size:255" json:"scope"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// 关联关系
	Account UserAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

// TableName 指定表名
func (UserOAuthToken) TableName() string {
	return "user_oauth_tokens"
}
