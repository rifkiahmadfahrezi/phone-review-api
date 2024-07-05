package controller

import (
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"final-project/utils/token"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type deleteUserInput struct {
	Password string `json:"password" bind:"required"`
}

type userUpdate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Get all users
// @Summary Get all account with user role.
// @Description Get a list of account with 'user' role.
// @Tags Users
// @Produce json
// @Success 200 {object} map[string][]string
// @Router /users [get]
func GetAllUser(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var users_data []models.User

	// apply filtering
	searchKeyword := c.Query("search")
	sort := c.Query("sort")

	query := db.Model(&models.User{})

	if searchKeyword != "" {
		q := fmt.Sprintf("%%%s%%", searchKeyword)
		query.Where("username LIKE ? OR email Like ?", q, q)
	}

	switch strings.ToLower(sort) {
	case "desc":
		query.Order("id DESC")
	default:
		query.Order("id ASC")
	}

	err := query.Select("id", "username", "email", "created_at", "updated_at").Where("role_id != 2").Find(&users_data).Error
	if err != nil {
		emptydata := make([]string, 0)
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, emptydata))
		return
	}

	// validsi jika data tidak ditemukan
	if searchKeyword != "" || sort != "" {
		if len(users_data) == 0 {
			emptydata := make([]string, 0)
			c.JSON(http.StatusNotFound,
				utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusNotFound, emptydata))
			return
		}
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, users_data))
}

// Get User by ID godoc
// @Summary Get single user by ID.
// @Description Get user data by ID.
// @Tags Users
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} map[string][]string
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	var user_data []models.User
	id := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Select("id", "username", "email", "created_at", "updated_at").Where("role_id != 2").First(&user_data, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, user_data))
}

// Delete account
// @Summary Delete user's own account
// @Description Will delete the user account itself
// @Tags Users
// @Produce json
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Param Body body deleteUserInput true "the body to delete user's own account"
// @Security BearerToken
// @Success 200 {object} []models.User
// @Router /users [delete]
func DeleteMyAccount(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// ambil user id
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}
	var user models.User
	// cek apakah user dengan id tsb ada
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	var input deleteUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	// verifikasi password
	err = models.VerifyPassword(input.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON("password salah, gagal menghapus akun", http.StatusBadRequest, nil))
		return
	}
	db.Delete(&user)
	c.JSON(http.StatusOK,
		utils.ResponseJSON(lib.MsgDeleted("user"), http.StatusOK, nil))
}

// Delete User by id  godoc
// @Summary Delete User by id .
// @Description Delete a User by id, only admin can access this route
// @Tags Users
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "User id"
// @Success 200 {object} map[string][]string
// @Router /users/{id} [delete]
func DeleteUserById(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// ambil user id
	userID := c.Param("id")
	var user models.User
	// cek apakah user dengan id tsb ada
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	db.Delete(&user)
	c.JSON(http.StatusOK,
		utils.ResponseJSON(lib.MsgDeleted("user"), http.StatusOK, user))
}

// Update User data godoc
// @Summary Update User data.
// @Description update its own user data
// @Tags Users
// @Produce json
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body userUpdate true "tExample JSON body to update User data"
// @Success 200 {object} models.User
// @Router /users/ [put]
func UpdateUser(c *gin.Context) {

	var input userUpdate
	db := c.MustGet("db").(*gorm.DB)
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Terjadi kesalahan pada server", http.StatusInternalServerError, nil))
		return
	}
	// cek data user dengan id tsb
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusNotFound, nil))
		return
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	var updated_data models.User
	updated_data.Username = input.Username
	updated_data.Email = input.Email

	// cek username sudah ada atau belum
	var user_exist []models.User
	if err := db.Where("username = ?", updated_data.Username).Find(&user_exist).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(user_exist) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgAlreadyExist("user"), http.StatusBadRequest, nil))
		return
	}

	// update ke tabel
	db.Model(&user).Updates(updated_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("user"), http.StatusOK, user))
}

// Get profiles data by User data ID godoc
// @Summary Get profiles data by User id.
// @Description Get all Users profile data by user id.
// @Tags Users
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} []models.User
// @Router /users/{id}/profile [get]
func GetUserProfileByID(c *gin.Context) {
	var user models.User

	db := c.MustGet("db").(*gorm.DB)

	userID := c.Param("id")
	if err := db.Preload("Profiles").Where("role_id != 2").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, user))
}
