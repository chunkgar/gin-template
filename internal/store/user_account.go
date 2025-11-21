package store

import "github.com/chunkgar/gin-template/internal/store/model"

type UserAccount interface {
	// FindOrCreate 根据设备ID查找或创建用户账号
	FindOrCreate(deviceID string, accountType model.AccountType) (*model.UserAccount, bool, error)

	// UpdateMeta 更新用户账号元数据
	UpdateMeta(id uint, meta map[string]any) error

	// Bind 绑定第三方或设备账号到指定用户
	Bind(userID uint, accountID string, accountType model.AccountType) (*model.UserAccount, int, error)

	// Unbind 解绑用户账号
	Unbind(userID uint, accountType model.AccountType) error
}
