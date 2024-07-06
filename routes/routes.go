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

	// untuk data account dgn role 'user'
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

	// untuk data account dgn role 'admins'
	adminMiddlewareRoutes := r.Group("/admins")
	// ⬇ Hanya bisa diakses account dgn role 'admin' yg sudah login
	adminMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	adminMiddlewareRoutes.Use(middleware.RoleMiddleware(db))
	adminMiddlewareRoutes.GET("", controller.GetAllAdmins)
	adminMiddlewareRoutes.GET("/:id", controller.GetAdminByID)
	adminMiddlewareRoutes.GET("/:id/profile", controller.GetAdminProfileByID)
	adminMiddlewareRoutes.POST("/register", controller.RegisterAdmin)

	// profile routes
	profileMiddlewareRoutes := r.Group("/profiles")
	profileMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	// ⬇ REGISTERED ACCOUNT ONLY (user/admin)
	// ID user diambil dari token
	profileMiddlewareRoutes.POST("", controller.CreateProfile)

	// role routes
	roleMiddlewareRoutes := r.Group("/roles")
	// ⬇ Hanya bisa diakses account dgn role 'admin' yg sudah login
	roleMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	roleMiddlewareRoutes.Use(middleware.RoleMiddleware(db))
	roleMiddlewareRoutes.GET("", controller.GetAllRoleData)
	roleMiddlewareRoutes.GET("/:id", controller.GetRoleDataByID)
	roleMiddlewareRoutes.POST("", controller.CreateRole)
	roleMiddlewareRoutes.PUT("/:id", controller.UpdateRole)
	roleMiddlewareRoutes.DELETE("/:id", controller.DeleteRoleByID)
	roleMiddlewareRoutes.GET("/:id/users", controller.GetUsersDataByRoleId)

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
	r.GET("/phones/:id/specification", controller.GetPhonesSpecByPhoneId)
	r.GET("/phones/:id/reviews", controller.GetReviewsDataByPhoneId)
	phonesMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	// ⬇ logged in account only (user/admin)
	phonesMiddlewareRoutes.POST("/:id/reviews", controller.CreateReview)
	phonesMiddlewareRoutes.PUT("/:id/reviews", controller.UpdateReview)
	// ⬇ ADMIN ONLY
	phonesMiddlewareRoutes.Use(middleware.RoleMiddleware(db))
	phonesMiddlewareRoutes.POST("", controller.CreatePhoneData)
	phonesMiddlewareRoutes.PUT("/:id", controller.UpdatePhoneData)
	phonesMiddlewareRoutes.DELETE("/:id", controller.DeletePhoneData)
	// create & update phone specification
	phonesMiddlewareRoutes.POST("/:id/specification", controller.CreateSpecification)
	phonesMiddlewareRoutes.PUT("/:id/specification", controller.UpdateSpecification)

	reviewsMiddlewareRoutes := r.Group("/reviews")
	reviewsMiddlewareRoutes.Use(middleware.JwtAuthMiddleware())
	reviewsMiddlewareRoutes.DELETE("/:id", controller.DeleteReviewById)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
