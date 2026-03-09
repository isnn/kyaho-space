package user

import (
	"kyaho-space/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(username string) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	CreateUser(user *entities.User) error
	UpdateToken(userID string, token string) error
	FindByToken(token string) (*entities.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUsername(username string) (*entities.User, error) {
	var u entities.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByEmail(email string) (*entities.User, error) {
	var u entities.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *repository) UpdateToken(userID string, token string) error {
	return r.db.Model(&entities.User{}).Where("id = ?", userID).Update("token", token).Error
}

func (r *repository) FindByToken(token string) (*entities.User, error) {
	var u entities.User
	err := r.db.Where("token = ?", token).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

