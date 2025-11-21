package mysql

import (
	"errors"

	"github.com/chunkgar/gin-template/internal/store/model"
	"gorm.io/gorm"
)

type userStore struct {
	db *gorm.DB
}

func (u *userStore) GetByID(userID uint) (*model.User, error) {
	var user model.User
	if err := u.db.Where("id = ?", userID).Preload("Accounts").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userStore) RequestDeletion(userID uint) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否已删除
		var count int64
		tx.Model(&model.UserDeletion{}).Where("user_id = ?", userID).Count(&count)
		if count > 0 {
			return errors.New("duplicated request")
		}

		// 创建用户删除记录
		userDeletion := &model.UserDeletion{
			UserID: userID,
			Status: "active",
		}
		if err := tx.Create(userDeletion).Error; err != nil {
			return err
		}

		return nil
	})
}
