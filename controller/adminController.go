package controller

import (
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register Admin godoc
// @Summary Register a new account as admin role.
// @Description registering a new account with admin role, only account with role admin can access this route
// @Tags Admins
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body RegisterInput true "the body to register a admin"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admins/register [post]
func RegisterAdmin(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		c.JSON(http.StatusBadRequest, utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Email = input.Email
	u.Password = input.Password

	// cek role user yg sedang meregister
	role_id, err := GetUserRoleId(c)
	if err != nil || role_id != 2 { // role 2 = admin
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
	}

	u.RoleID = 2 // set default role (user)

	_, err = u.SaveUser(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	user := map[string]string{
		"username": input.Username,
		"email":    input.Email,
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("Register admin berhasil", http.StatusOK, map[string]any{
		"user": user,
	}))

}

// admin role

// Get all admins
// @Summary Get all account with 'admin' role.
// @Description Get a list of account with 'admin' role, only role admin can acces this route
// @Tags Admins
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string][]string
// @Router /admins [get]
func GetAllAdmins(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var admins_data []models.User

	err := db.Select("id", "username", "email", "created_at", "updated_at").Where("role_id = 2").Find(&admins_data).Error
	if err != nil {
		emptydata := make([]string, 0)
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, emptydata))
		return
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, admins_data))
}

// Get Admin by ID godoc
// @Summary Get single admin by ID.
// @Description Get admin data by ID. only role admin can acces to this route
// @Tags Admins
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "admin id"
// @Success 200 {object} map[string][]string
// @Router /admins/{id} [get]
func GetAdminByID(c *gin.Context) {
	var user_data []models.User
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Select("id", "username", "email", "created_at", "updated_at").Where("role_id = 2").First(&user_data, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("admin"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, user_data))
}

// Get profiles data by Admin data ID godoc
// @Summary Get profiles data by Admin id.
// @Description Get all Admins profile data by admin id. only admin can access this route
// @Tags Admins
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "admin id"
// @Success 200 {object} []models.User
// @Router /admins/{id}/profile [get]
func GetAdminProfileByID(c *gin.Context) {
	var user models.User

	db := c.MustGet("db").(*gorm.DB)

	userID := c.Param("id")
	if err := db.Preload("Profiles").Where("role_id = 2").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, user))
}
