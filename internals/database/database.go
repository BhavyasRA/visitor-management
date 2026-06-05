package database

import (
	"fmt"
	"log"

	"entry-system/internals/config"
	"entry-system/internals/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	user := config.GetEnv("DB_USER", "postgres")
	password := "Google$$321"
	dbName := config.GetEnv("DB_NAME", "entry_system")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbName,
		port,
	)
	log.Println("DB_HOST:", host)
	log.Println("DB_USER:", user)
	log.Println("DB_NAME:", dbName)
	log.Println("DB_PASSWORD_LENGTH:", len(password))

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal("Database connection failed")
	}

	DB = db

	err = DB.AutoMigrate(
		&models.User{},
		&models.Authentication{},
		&models.Role{},
		&models.Permission{},
		&models.Visitor{},
		&models.VisitorDocument{},
		&models.EntryLog{},
	)

	if err != nil {
		log.Fatal("Migration failed")
	}

	log.Println("Database connected successfully")
}
