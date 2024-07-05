package config

import (
	"final-project/models"
	"final-project/utils"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {

	dbProvider := utils.GetEnv("DB_PROVIDER", "mysql")
	environment := utils.GetEnv("ENVIRONMENT", "development")

	var db *gorm.DB

	if environment == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	host := utils.GetEnv("DB_HOST", "127.0.0.1")
	port := utils.GetEnv("DB_PORT", "3306")
	user := utils.GetEnv("DB_USER", "root")
	password := utils.GetEnv("DB_PASSWORD", "password")
	dbname := utils.GetEnv("DB_NAME", "db_phone_review")

	switch dbProvider {
	case "postgre":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta", host, user, password, dbname, port)
		dbPsg, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		db = dbPsg
	default:
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

		dbGorm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			panic(err.Error())
		}

		db = dbGorm
	}

	// Migrate the schema
	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Profile{},
		&models.Brand{},
		&models.Phone{},
		&models.Specification{},
		&models.Review{},
		&models.Comment{},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

// func seedData(db *gorm.DB) {
// 	// Cek apakah tabel Role kosong
// 	var count int64
// 	db.Model(&models.Role{}).Count(&count)
// 	if count == 0 {
// 		 roles := []models.Role{
// 			  {Name: "user"},
// 			  {Name: "admin"},
// 		 }

// 		 for _, role := range roles {
// 			  db.Create(&role)
// 		 }
// 	}

// 	db.Model(&models.User{}).Count(&count)
// 	if count == 0 {
// 		 users := []models.User{
// 			  {Username: "admin", Email: "admin@test.com", Password: "hashed_password", RoleID: 2},
// 			  {Username: "user", Email: "user@test.com", Password: "hashed_password", RoleID: 1},
// 		 }

// 		 for _, user := range users {
// 			  db.Create(&user)
// 		 }
// 	}
// }
