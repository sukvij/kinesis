package database

import (
	"fmt"
	"log"
	"os"
	"time"

	userLogModel "vijju/user-logs/model"
	userModel "vijju/user/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes the MySQL database connection and migrates schemas
func InitDB() *gorm.DB {
	// Load .env file
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println(dsn)
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to MySQL (attempt %d/5): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err := db.AutoMigrate(&userModel.User{}, &userLogModel.Log{}); err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}
	return db
}
