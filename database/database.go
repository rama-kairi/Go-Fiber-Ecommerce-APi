package database

import (
	"log"
	"os"

	"github.com/rama-kairi/fiber-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	Db *gorm.DB
}

var Database DBInstance

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err.Error())
		os.Exit(1)
	}

	log.Println("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{})

	Database = DBInstance{Db: db}
}
