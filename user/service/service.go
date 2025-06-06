package service

import (
	kinesisModel "vijju/kinesis/model"
	kinesisSer "vijju/kinesis/service"
	"vijju/user/model"
	"vijju/user/repository"

	"gorm.io/gorm"
)

type UserService struct {
	DB   *gorm.DB
	User *model.User
}

func NewService(db *gorm.DB, user *model.User) *UserService {
	return &UserService{DB: db, User: user}
}

func (service *UserService) GetAllUsers() (*[]model.User, error) {
	repo := repository.NewRepository(service.DB, service.User)
	return repo.GetAllUsers()
}

func (service *UserService) CreateUser() (*model.User, error) {
	repo := repository.NewRepository(service.DB, service.User)
	res, err := repo.CreateUser()
	// call kinesis put data
	kinesisService := kinesisSer.NewKinesisService(service.DB)
	kinesisService.PutData(&kinesisModel.UserEvent{UserID: repo.User.ID, Action: "user-created"}, int(service.User.ID))
	return res, err
}
