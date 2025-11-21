package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AccountType 账号类型枚举
type AccountType string

const (
	AccountTypeEmail     AccountType = "email"
	AccountTypePhone     AccountType = "phone"
	AccountTypeWechat    AccountType = "wechat"
	AccountTypeQQ        AccountType = "qq"
	AccountTypeApple     AccountType = "apple"
	AccountTypeGoogle    AccountType = "google"
	AccountTypeIDFV      AccountType = "idfv"
	AccountTypeIDFA      AccountType = "idfa"
	AccountTypeAndroidID AccountType = "android_id"
)

// UserAccount 用户账号表（登录方式映射）
type UserAccount struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	UserID            uint              `gorm:"not null;index" json:"user_id"`
	AccountType       AccountType       `gorm:"type:varchar(20);index:,unique,composite:uk_account,priority:1;not null" json:"account_type"`
	AccountIdentifier string            `gorm:"size:255;not null;index:,unique,composite:uk_account,priority:2;comment:账号标识（邮箱/手机号/第三方ID）" json:"account_identifier"`
	PasswordHash      *string           `gorm:"size:255;comment:密码哈希（第三方登录为NULL）" json:"-"` // json:"-" 不返回给前端
	IsVerified        bool              `gorm:"default:false;comment:是否已验证" json:"is_verified"`
	IsPrimary         bool              `gorm:"default:false;comment:是否为主账号" json:"is_primary"`
	Meta              datatypes.JSONMap `gorm:"column:meta;type:json;comment:'账号元数据'" json:"meta"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`

	// 关联关系
	User       User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OAuthToken *UserOAuthToken `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"oauth_token,omitempty"`
}

// TableName 指定表名
func (UserAccount) TableName() string {
	return "user_accounts"
}

// BeforeCreate GORM Hook - 创建前的验证
func (ua *UserAccount) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加业务逻辑验证
	return nil
}
