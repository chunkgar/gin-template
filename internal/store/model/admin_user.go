package model

import "gorm.io/gorm"

type AdminUser struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;column:username;type:varchar(255);not null"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
	Role     string `gorm:"column:role;type:varchar(100);default:'admin'"`
	IsActive bool   `gorm:"column:is_active;default:true"`
}
