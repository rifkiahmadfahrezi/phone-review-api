package main

import (
	"final-project/config"
	"final-project/docs"
	"final-project/routes"
	"final-project/seed"
	"final-project/utils"
	"log"

	"github.com/joho/godotenv"
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/

func main() {

	environment := utils.GetEnv("ENVIRONMENT", "development")
	if environment == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	docs.SwaggerInfo.Title = "Phone reviews API"
	docs.SwaggerInfo.Description = "This is a phone reviews API"
	docs.SwaggerInfo.Version = "1.0"
	if environment == "development" {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	} else {
		docs.SwaggerInfo.Host = "xenophobic-blanch-superbootcamp-45cdd211.koyeb.app/"
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	db := config.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// load initial data role dan user
	seed.Load(db)

	r := routes.SetupRouter(db)
	r.Run()

}
