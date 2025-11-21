package model

import (
	"gorm.io/gorm"
)

// User 用户基础信息表
type User struct {
	gorm.Model
	Nickname  string    `gorm:"size:50" json:"nickname"`
	AvatarURL string    `gorm:"size:255;column:avatar_url" json:"avatar_url"`

	// 关联关系
	Accounts   []UserAccount `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"accounts,omitempty"`
	Membership *Membership   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"membership,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
