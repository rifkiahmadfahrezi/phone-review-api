package controller

import (
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"final-project/utils/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type profileInput struct {
	Biodata  string    `json:"biodata"`
	ImageURL string    `json:"image_url" binding:"url"`
	FullName string    `json:"full_name"`
	Birthday time.Time `json:"birthday"`
	Email    string    `json:"email" binding:"email"`
}

// Create Profile for user godoc
// @Summary Create Profile for user
// @Description Creating a profile data for user, only registered user can access this route
// @Tags Profiles
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body profileInput true "example JSON body to create a profile for user, user_id is taken from the authorization token"
// @Produce json
// @Success 200 {object} models.Profile
// @Router /profiles [post]
func CreateProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input profileInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}
	// ambil user id
	userID, err := token.ExtractTokenID(c)
	// cek data Profile sudah ada atau blm
	var profile []models.Profile
	if err := db.Where("user_id = ?", userID).Find(&profile).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(profile) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgAlreadyExist("profile"), http.StatusBadRequest, nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Gagal mengambil id user", http.StatusInternalServerError, nil))
		return
	}

	profile_data := models.Profile{
		Biodata:  input.Biodata,
		ImageURL: input.ImageURL,
		FullName: input.FullName,
		Birthday: &input.Birthday,
		Email:    input.Email,
		UserID:   userID,
	}

	db.Create(&profile_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("profile"), http.StatusOK, profile_data))
}

// Update Profile for user godoc
// @Summary Update Profile for user
// @Description Updating a profile data for user, only registered user can access this route
// @Tags Profiles
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body profileInput true "example JSON body to update a profile for user, user_id is taken from the authorization token"
// @Produce json
// @Success 200 {object} models.Profile
// @Router /profiles [put]
func UpdateProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// ambil user id
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Gagal mengambil id user", http.StatusInternalServerError, nil))
		return
	}

	// cek apakah user dengan id (userID) ada
	var user []models.User
	if err := db.Where("id = ?", userID).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(user) == 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusBadRequest, nil))
		return
	}

	// Validate input
	var input profileInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	profile_data := models.Profile{
		Biodata:  input.Biodata,
		ImageURL: input.ImageURL,
		FullName: input.FullName,
		Birthday: &input.Birthday,
		Email:    input.Email,
		UserID:   userID,
	}

	db.Create(&profile_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("profile"), http.StatusOK, profile_data))
}
