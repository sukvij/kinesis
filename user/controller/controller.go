package controller

import (
	"net/http"
	"vijju/user/model"
	"vijju/user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (userController *UserController) UserApis(app *gin.Engine) {
	app.GET("/users", userController.getAllUsers)
	app.POST("/users", userController.createUser)
}

func (controller *UserController) getAllUsers(ctx *gin.Context) {
	service := service.NewService(controller.DB, nil)
	res, err := service.GetAllUsers()
	if err != nil {
		ctx.JSON(100, err)
		return
	}
	ctx.JSON(200, res)
}

func (controller *UserController) createUser(ctx *gin.Context) {
	user := &model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service := service.NewService(controller.DB, user)
	res, err := service.CreateUser()
	if err != nil {
		ctx.JSON(100, err)
		return
	}
	ctx.JSON(200, res)
}
