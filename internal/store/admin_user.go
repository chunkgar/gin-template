package store

import "github.com/chunkgar/gin-template/internal/store/model"

type AdminUser interface {
	GetByName(name string) (*model.AdminUser, error)
	GetByID(id uint) (*model.AdminUser, error)
}
