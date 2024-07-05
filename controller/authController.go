package controller

import (
	"errors"
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"final-project/utils/token"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username"`
	Email    string `json:"Email"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginUser godoc
// @Summary Login.
// @Description Logging in to get jwt token to access admin or user api by roles.
// @Tags Auth
// @Param Body body LoginInput true "the body to login a user choose using email or username"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		c.JSON(http.StatusBadRequest, utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Email = input.Email
	u.Password = input.Password

	token, err := models.LoginCheck(u.Username, u.Email, u.Password, db)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, utils.ResponseJSON("Username atau password salah", http.StatusBadRequest, nil))
		return
	}

	user := map[string]string{
		"username": u.Username,
		"email":    u.Email,
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("Login berhasil", http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	}))

}

// Register godoc
// @Summary Register a user.
// @Description registering a user from public access.
// @Tags Auth
// @Param Body body RegisterInput true "the body to register a user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func RegisterUser(c *gin.Context) {
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
	u.RoleID = 1 // set default role (user)

	_, err := u.SaveUser(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	user := map[string]string{
		"username": input.Username,
		"email":    input.Email,
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("Register berhasil", http.StatusOK, map[string]any{
		"user": user,
	}))
}

type changePasswordInput struct {
	CurrentPassword string `json:"current_password" bind:"required"`
	NewPassword     string `json:"new_password" bind:"required"`
}

// ChangePassword godoc
// @Summary Change password
// @Description changging current logged in user's password
// @Tags Auth
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body changePasswordInput true "body for changing user's password, user id is taken from the authorization token"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/change-password [put]
func ChangePassword(c *gin.Context) {
	var input changePasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	// currentUserId adalah id user yg sedang login
	currentUserId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	// cek apakah user dengan id tsb ada
	if err := db.Where("id = ?", currentUserId).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	currentPassword := user.Password

	// verifikasi password pw lama di db dan pw lama di input
	err = models.VerifyPassword(input.CurrentPassword, currentPassword)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON("Password salah, gagal memperbarui password", http.StatusBadRequest, nil))
		return
	}

	newHashedPassword, err := models.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Password gagal diperbarui, ada masalah diserver", http.StatusInternalServerError, nil))
		return
	}

	var updated_data models.User
	updated_data.Password = newHashedPassword
	updated_data.UpdatedAt = time.Now()

	// Update the user record with the new password
	result := db.Model(&user).Updates(updated_data)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Error updating user: "+result.Error.Error(), http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("Password berhasil diperbarui", http.StatusOK, nil))
}

func GetUserRoleId(c *gin.Context) (uint, error) {
	db := c.MustGet("db").(*gorm.DB)
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return 0, err
	}

	var user []models.User
	if err := db.Where("id = ?", userID).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return 0, err
	}

	// jika data tidak ditemukan
	if len(user) == 0 {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusNotFound, nil))
		return 0, errors.New(lib.ErrMsgNotFound("user"))
	}

	return user[0].RoleID, nil
}
