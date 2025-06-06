package service

import (
	"vijju/user-logs/model"
	"vijju/user-logs/repository"

	"gorm.io/gorm"
)

type UserLogService struct {
	Db  *gorm.DB
	Log *model.Log
}

func NewUserLogService(db *gorm.DB, log *model.Log) *UserLogService {
	return &UserLogService{Db: db, Log: log}
}

func (service *UserLogService) CreateLog() {
	userLogRepo := repository.NewRepository(service.Db, service.Log)
	userLogRepo.CreateLog()
}
