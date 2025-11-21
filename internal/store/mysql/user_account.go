package mysql

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chunkgar/gin-template/internal/code"
	"github.com/chunkgar/gin-template/internal/store/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type userAccountStore struct {
	db *gorm.DB
}

// GenerateUsername 生成指定长度的随机用户名
func generateUsername(length int) string {
	// 定义字符集：大小写字母 + 数字
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 初始化随机数种子
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成随机用户名
	username := make([]byte, length)
	for i := range username {
		username[i] = charset[rand.Intn(len(charset))]
	}

	return string(username)
}

func (u *userAccountStore) FindOrCreate(deviceID string, accountType model.AccountType) (*model.UserAccount, bool, error) {
	var userAccount model.UserAccount
	isNew := false

	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("User").Where("account_identifier = ? AND account_type = ?", deviceID, accountType).First(&userAccount).Error; err != nil {
			// 没有找到，创建用户和用户账号
			if err == gorm.ErrRecordNotFound {
				// 创建用户
				user := model.User{
					Nickname:  "用户" + generateUsername(8), // TODO: 随机昵称？
					AvatarURL: "",                         // TODO: 随机头像？
				}

				if err = tx.Create(&user).Error; err != nil {
					return err
				}

				userAccount = model.UserAccount{
					UserID:            user.ID,
					AccountIdentifier: deviceID,
					AccountType:       accountType,
					IsVerified:        true,
				}
				if err = tx.Create(&userAccount).Error; err != nil {
					return err
				}

				userAccount.User = user
				isNew = true

			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, isNew, err
	}

	return &userAccount, isNew, nil
}

func (u *userAccountStore) UpdateMeta(id uint, meta map[string]any) error {
	return u.db.Model(&model.UserAccount{}).Where("id = ?", id).Update("meta", datatypes.JSONMap(meta)).Error
}

func (u *userAccountStore) Bind(userID uint, accountID string, accountType model.AccountType) (*model.UserAccount, int, error) {
	var found model.UserAccount

	if err := u.db.Where("account_identifier = ? AND account_type = ?", accountID, accountType).First(&found).Error; err == nil {
		return nil, code.ErrAccountAlreadyBound, fmt.Errorf("account already bound")
	} else if err != gorm.ErrRecordNotFound {
		return nil, code.ErrDatabase, err
	}

	var ua model.UserAccount
	if err := u.db.Transaction(func(tx *gorm.DB) error {
		ua = model.UserAccount{
			UserID:            userID,
			AccountIdentifier: accountID,
			AccountType:       accountType,
		}
		if err := tx.Create(&ua).Error; err != nil {
			return err
		}
		return tx.Preload("User").First(&ua, ua.ID).Error
	}); err != nil {
		return nil, code.ErrDatabase, err
	}

	return &ua, code.ErrSuccess, nil
}

func (u *userAccountStore) Unbind(userID uint, accountType model.AccountType) error {
	return u.db.Delete(&model.UserAccount{}, "user_id = ? AND account_type = ?", userID, accountType).Error
}
