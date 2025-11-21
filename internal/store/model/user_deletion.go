package model

import "gorm.io/gorm"

type UserDeletion struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	Status  string `gorm:"type:varchar(20);index;comment:删除状态(active:正常,deleted:已删除)" json:"status"`
	AdminID *uint  `gorm:"index" json:"admin_id"`

	// 关联关系
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AdminUser AdminUser `gorm:"foreignKey:AdminID" json:"admin_user,omitempty"`
}
