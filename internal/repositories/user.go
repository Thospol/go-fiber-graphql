package repositories

import "gorm.io/gorm"

// Repository repository interface
type UserRepository interface {
	Create(db *gorm.DB, i interface{}) error
	FindOneObjectByID(db *gorm.DB, id uint64, i interface{}) error
}

type userRepository struct {
	Repository
}

// NewUserRepository new user repository
func NewUserRepository() UserRepository {
	return &userRepository{
		NewRepository(),
	}
}
