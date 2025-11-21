package model

import "time"

// MembershipLevel 会员等级枚举
type MembershipLevel string

const (
	MembershipFree    MembershipLevel = "free"
	MembershipBasic   MembershipLevel = "basic"
	MembershipPremium MembershipLevel = "premium"
	MembershipVIP     MembershipLevel = "vip"
)

// Membership 会员信息表
type Membership struct {
	ID                    uint            `gorm:"primaryKey" json:"id"`
	UserID                uint            `gorm:"uniqueIndex;not null" json:"user_id"`
	MembershipLevel       MembershipLevel `gorm:"type:varchar(20);default:free" json:"membership_level"`
	SubscriptionStartDate *time.Time      `json:"subscription_start_date"`
	SubscriptionEndDate   *time.Time      `json:"subscription_end_date"`
	AutoRenew             bool            `gorm:"default:false" json:"auto_renew"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Membership) TableName() string {
	return "memberships"
}

// IsActive 判断会员是否有效
func (m *Membership) IsActive() bool {
	if m.SubscriptionEndDate == nil {
		return false
	}
	return m.SubscriptionEndDate.After(time.Now())
}
