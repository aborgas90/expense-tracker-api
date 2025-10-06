package db

import (
	"fmt"
	"log"
	"os"

	"github.com/aborgas90/expense-tracker-api/internal/model"

	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var databaseInstance *gorm.DB

func InIt() *gorm.DB {

	var err error
	databaseInstance, err = connectDb()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = performMigration()
	if err != nil {
		log.Fatalf("Could not auto migrate: %v", err)
	}
	return databaseInstance
}

func connectDb() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	if err := godotenv.Load(); err != nil {
		log.Println(" .env file not found, use system env instead")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("CHECK ENV " + dbUsername + " " + dbPassword + " " + dbHost + " " + dbPort + " " + dbName)

	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		return nil, err
	}
	return db, nil
}

func performMigration() error {
	return databaseInstance.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Transaction{},
		&model.Goal{},
		&model.GoalDeposit{},
	)
}
