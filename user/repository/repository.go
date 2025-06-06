package repository

import (
	"vijju/user/model"

	"gorm.io/gorm"
)

type Repository struct {
	Db   *gorm.DB
	User *model.User
}

func NewRepository(db *gorm.DB, user *model.User) *Repository {
	return &Repository{Db: db, User: user}
}

func (repository *Repository) GetAllUsers() (*[]model.User, error) {
	users := &[]model.User{}
	if err := repository.Db.Find(users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repository *Repository) CreateUser() (*model.User, error) {
	// users := &[]model.User{}
	if err := repository.Db.Create(repository.User).Error; err != nil {
		return nil, err
	}
	return repository.User, nil
}
