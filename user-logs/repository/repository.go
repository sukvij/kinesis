package repository

import (
	"fmt"
	"vijju/user-logs/model"

	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	Log *model.Log
}

func NewRepository(db *gorm.DB, log *model.Log) *Repository {
	return &Repository{DB: db, Log: log}
}

func (repo *Repository) CreateLog() {
	if err := repo.DB.Create(repo.Log).Error; err != nil {
		fmt.Println("logs creation error occured. err - ", err)
	}
}
