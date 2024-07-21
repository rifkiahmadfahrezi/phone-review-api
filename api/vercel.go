package api

import (
	"final-project/config"
	"final-project/docs"
	"final-project/routes"
	"final-project/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	App *gin.Engine
)

func init() {
	App = gin.New()

	environment := utils.GetEnv("ENVIRONMENT", "development")

	if environment == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	docs.SwaggerInfo.Title = "Phone review REST API"
	docs.SwaggerInfo.Description = "This is REST API Phone review."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = utils.GetEnv("API_HOST", "localhost:8080")
	if environment == "development" {
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}
	db := config.ConnectDatabase()
	// sqlDB, _ := db.DB()
	// defer sqlDB.Close()

	// load initial role & user/admin auth information
	// seed.Load()

	routes.SetupRouter(db, App)
}

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	App.ServeHTTP(w, r)
}
