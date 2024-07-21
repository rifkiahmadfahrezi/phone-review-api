package main

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/
import "final-project/api"

func main() {

	api.App.Run()

	// environment := utils.GetEnv("ENVIRONMENT", "development")
	// if environment == "development" {
	// 	err := godotenv.Load()
	// 	if err != nil {
	// 		log.Fatal("Error loading .env file")
	// 	}
	// }

	// api_host := utils.GetEnv("API_HOST", "localhost")

	// docs.SwaggerInfo.Title = "Phone reviews API"
	// docs.SwaggerInfo.Description = "Phone reviews API, Rifki ahmad fahrezi"
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = api_host
	// if environment == "development" {
	// 	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// } else {
	// 	docs.SwaggerInfo.Schemes = []string{"https"}
	// }

	// db := config.ConnectDatabase()
	// sqlDB, _ := db.DB()
	// defer sqlDB.Close()

	// // load initial data role dan user
	// seed.Load(db)

	// r := routes.SetupRouter(db)
	// r.Run()

}
