package store

import "github.com/chunkgar/gin-template/internal/store/model"

type User interface {
	GetByID(userID uint) (*model.User, error)
	RequestDeletion(userID uint) error
}
