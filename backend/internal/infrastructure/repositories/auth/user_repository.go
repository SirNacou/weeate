package auth

import (
	"errors"

	domain_auth "github.com/SirNacou/weeate/backend/internal/domain/auth"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain_auth.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) WithTx(tx *gorm.DB) domain_auth.UserRepository {
	return &GormUserRepository{db: tx}
}

func (r *GormUserRepository) CreateUser(user *domain_auth.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) GetUserByUsername(username string) (*domain_auth.User, error) {
	var user domain_auth.User
	if err := r.db.Where(&domain_auth.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_auth.ErrUserNotFound
		}

		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&domain_auth.User{}, &domain_auth.User{Model: gorm.Model{ID: id}}).Error
}
