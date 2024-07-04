package config

import (
	"final-project/models"
	"final-project/utils"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// host := "aws-0-ap-southeast-1.pooler.supabase.com"
	// port := "6543"
	// user := "postgres.qkhezuksjpypdrluzkpu"
	// password := "dbphonereview123:)"
	// dbname := "postgres"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := utils.GetEnv("DB_HOST", "localhost")
	port := utils.GetEnv("DB_PORT", "3306")
	user := utils.GetEnv("DB_USER", "root")
	password := utils.GetEnv("DB_PASSWORD", "password")
	dbname := utils.GetEnv("DB_NAME", "db_phone_review")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
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
