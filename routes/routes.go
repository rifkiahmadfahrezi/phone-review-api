package routes

import (
	"final-project/controller"
	"final-project/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"}

	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true
	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")

	r.Use(cors.New(corsConfig))

	// set db to gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// auth routes
	authMiddlewareRoutes := r.Group("/auth")
	// ⬇ PUBLIC ROUTES
	r.POST("/auth/register", controller.RegisterUser)
	r.POST("/auth/login", controller.Login)
	authMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	// ⬇ REGISTERED ACCOUNT ONLY (user/admin)
	// ID untuk change password diambil dari token
	authMiddlewareRoutes.PUT("/change-password", controller.ChangePassword)
	// ⬇ ADMIN ONLY
	authMiddlewareRoutes.Use(middleware.RoleMiddleware(db))
	authMiddlewareRoutes.POST("/register-admin", controller.RegisterAdmin)

	// users route
	userMiddlewareRoutes := r.Group("/users")
	// ⬇ PUBLIC ROUTES
	r.GET("/users", controller.GetAllUser)      // get user role accounts only
	r.GET("/users/:id", controller.GetUserByID) // get user role accounts only
	r.GET("/users/:id/profile", controller.GetUserProfileByID)
	userMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	// ⬇ For registered account only (user/admin)
	userMiddlewareRoutes.PUT("", controller.UpdateUser)
	userMiddlewareRoutes.DELETE("", controller.DeleteMyAccount)
	// ⬇ For ADMIN ONLY
	userMiddlewareRoutes.Use(middleware.RoleMiddleware(db))
	userMiddlewareRoutes.DELETE("/:id", controller.DeleteUserById)

	// profile routes
	profileMiddlewareRoutes := r.Group("/profiles")
	// ⬇ PUBLIC ROUTES

	profileMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	// ⬇ REGISTERED ACCOUNT ONLY (user/admin)
	// ID user diambil dari token
	profileMiddlewareRoutes.POST("", controller.CreateProfile)

	// brands route
	brandsMiddlewareRoutes := r.Group("/brands")
	// ⬇ BRANDS PUBLIC ROUTES
	r.GET("/brands", controller.GetAllBrandData)
	r.GET("/brands/:id", controller.GetBrandById)
	r.GET("/brands/:id/phones", controller.GetPhonesDataByBrandId)
	brandsMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	// ⬇ ADMIN ONLY
	brandsMiddlewareRoutes.Use(middleware.RoleMiddleware(db))
	brandsMiddlewareRoutes.POST("", controller.CreateBrand)
	brandsMiddlewareRoutes.PUT("/:id", controller.UpdateBrand)
	brandsMiddlewareRoutes.DELETE("/:id", controller.DeleteBrandByID)

	// phones route
	phonesMiddlewareRoutes := r.Group("/phones")
	// ⬇ PUBLIC ROUTES
	r.GET("/phones/:id", controller.GetPhoneById)
	r.GET("/phones", controller.GetAllPhoneData)
	phonesMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	// ⬇ ADMIN ONLY
	phonesMiddlewareRoutes.Use(middleware.RoleMiddleware(db))
	phonesMiddlewareRoutes.POST("", controller.CreatePhoneData)
	phonesMiddlewareRoutes.PUT("/:id", controller.UpdatePhoneData)
	phonesMiddlewareRoutes.DELETE("/:id", controller.DeletePhoneData)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
