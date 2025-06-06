package main

import (
	"vijju/database"
	kinesisSer "vijju/kinesis/service"
	"vijju/user/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.InitDB()
	kinesisServie := kinesisSer.NewKinesisService(db)
	go func() {
		kinesisServie.ReadDataFromKinesis()
	}()
	app := gin.Default()
	// userController := controller.NewUserController(db)
	(controller.NewUserController(db)).UserApis(app)
	app.Run(":8080")

}
