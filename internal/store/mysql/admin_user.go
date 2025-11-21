package mysql

import (
	"github.com/chunkgar/gin-template/internal/store/model"
	"gorm.io/gorm"
)

type adminUserStore struct {
	db *gorm.DB
}

func (s *adminUserStore) GetByName(name string) (*model.AdminUser, error) {
	var adminUser model.AdminUser
	if err := s.db.Where("username = ? AND is_active = ?", name, true).First(&adminUser).Error; err != nil {
		return nil, err
	}
	return &adminUser, nil
}

func (s *adminUserStore) GetByID(id uint) (*model.AdminUser, error) {
	var adminUser model.AdminUser
	if err := s.db.Where("id = ? AND is_active = ?", id, true).First(&adminUser).Error; err != nil {
		return nil, err
	}
	return &adminUser, nil
}
