package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gotodo/internal/persistence/record"
	"log"
	"os"
	"time"
)

// NewDatabaseConnection
// Do: Function to open connection with database mysql
// Param: Context
func NewDatabaseConnection(ctx context.Context, path string) (db *gorm.DB, errs error) {
	err := godotenv.Load(path)

	if err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}

	// Do get from environment file
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	hostname := os.Getenv("MYSQL_HOST")
	databaseName := os.Getenv("MYSQL_NAME")

	databaseConnection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, hostname, databaseName)

	//connection, err := gorm.Open(mysql.Open(databaseConnection), &gorm.Config{})

	connection, err := gorm.Open(mysql.Open(databaseConnection), &gorm.Config{})

	if err != nil {
		panic("Failed to create connection to database")
	}

	err = connection.AutoMigrate(&record.TaskRecord{}, &record.AccountRecord{}, &record.UserDetailRecord{})
	if err != nil {
		return nil, err
	}

	// Check if the MyModel table exists in the database
	if connection.Migrator().HasTable(&record.TaskRecord{}) ||
		connection.Migrator().HasTable(&record.AccountRecord{}) ||
		connection.Migrator().HasTable(&record.UserDetailRecord{}) {
		fmt.Println("table record already migrations")
	} else {
		fmt.Println("table record not have migrations")
	}

	sqlDB, err := connection.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = sqlDB.PingContext(ctx)

	if err != nil {
		_ = fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database")

	return connection, nil
}
