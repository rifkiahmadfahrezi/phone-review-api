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
}

// Create Profile for user godoc
// @Summary Create Profile for user
// @Description Creating a profile data for user, user ID is taken from JWT Token so only acount's owner can create the profile
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

	// lakukan validasi image url hanya jika diisi oleh user
	// krn image url sudah diberi nilai default di struct nya
	if input.ImageURL != "" {
		if err := c.ShouldBindJSON(&input); err != nil {
			errorMessage := utils.CustomBindError(err)
			if errorMessage != "" {
				c.JSON(http.StatusBadRequest,
					utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
			}
			return
		}
	}
	// ambil user id
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Gagal mengambil id user", http.StatusInternalServerError, nil))
		return
	}
	// cek data Profile sudah ada atau blm
	var profile []models.Profile
	if err := db.Where("user_id = ?", userID).Find(&profile).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(profile) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON("anda sudah membuat profile", http.StatusBadRequest, nil))
		return
	}

	profile_data := models.Profile{
		Biodata:  input.Biodata,
		ImageURL: input.ImageURL,
		FullName: input.FullName,
		Birthday: &input.Birthday,
		UserID:   userID,
	}

	if err := db.Create(&profile_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("profile"), http.StatusOK, profile_data))
}

// Update Profile for user godoc
// @Summary Update Profile for user
// @Description Updating a profile data for user, user ID is taken from JWT Token so only acount's owner can update the profile
// @Tags Profiles
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body profileInput true "example JSON body to update a profile for user, user_id is taken from the authorization token"
// @Produce json
// @Success 200 {object} models.Profile
// @Router /profiles [put]
func UpdateProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	// user id dari token
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	// cek jika user sudah mengisi profile atau blm
	var profile models.Profile
	if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON("user ini belum mengisikan profile", http.StatusNotFound, nil))
		return
	}

	// Update yg diinput saja
	if input.Biodata != "" {
		profile.Biodata = input.Biodata
	}
	if input.Birthday != nil {
		profile.Birthday = input.Birthday
	}
	if input.FullName != "" {
		profile.FullName = input.FullName
	}

	profile.UpdatedAt = time.Now()

	if err := db.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Failed to update profile", http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("profile"), http.StatusOK, profile))
}
