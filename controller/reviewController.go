package controller

import (
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"final-project/utils/token"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type reviewInput struct {
	Rating  uint   `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content" binding:"required"`
}
type reviewUpdate struct {
	Rating  uint   `json:"rating" binding:"min=1,max=5"`
	Content string `json:"content" `
}

// Create New Review godoc
// @Summary Create New Review
// @Description This route will create review data , user ID is taken from the JWT token, one user only can give one review to one phone
// @Tags Phones
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "Phone id"
// @Param Body body reviewInput true "example JSON body to create a new Review"
// @Produce json
// @Success 200 {object} models.Review
// @Router /phones/{id}/reviews [post]
func CreateReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input reviewInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	phoneIDstr := c.Param("id")

	phoneID, err := strconv.Atoi(phoneIDstr)

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	// user id dari token
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	// cek phone ada atau tidak
	var phone []models.Phone
	if err := db.Where("id = ?", phoneID).Find(&phone).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(phone) == 0 {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusNotFound, nil))
		return
	}

	// cek apakah user sudah memberikan review ke phone dengan id (phoneID) atau belm
	var rev []models.Review
	if err := db.Where("phone_id = ? AND user_id = ?", phoneID, userID).First(&rev).Error; err != nil {
		// jika data tidak ditemukan berarti user belum memberikan review terhadap phone ini
		// dan user boleh menambahkan reviewnya
		if err.Error() != "record not found" {
			c.JSON(http.StatusInternalServerError,
				utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
			return
		}
	}

	// jika data ditemukan berarti user sudah memberikan review terhadap phone ini
	// dan user tidak boleh menambahkan reviewnya
	if len(rev) > 0 {
		msg := "anda sudah memberikan review ke phone ini"
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(msg, http.StatusBadRequest, nil))
		return
	}

	review_data := models.Review{
		Rating:  input.Rating,
		Content: input.Content,
		UserID:  userID,
		PhoneID: uint(phoneID),
	}

	db.Create(&review_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("review"), http.StatusOK, review_data))
}

// Update Review for phone godoc
// @Summary Update Review for phone
// @Description This route will update review data , user ID is taken from the JWT token
// @Tags Reviews
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "review id"
// @Param Body body reviewInput true "example JSON body to update a review for phone"
// @Produce json
// @Success 200 {object} models.Review
// @Router /reviews/{id} [put]
func UpdateReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input reviewUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	reviewID := c.Param("id")

	// user id dari token
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	// cek data review ada atau tdk dan milik user yg login
	var rev models.Review
	if err := db.Where("id = ? AND user_id = ?", reviewID, userID).First(&rev).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON("Review tidak ditemukan atau bukan milik user", http.StatusNotFound, nil))
		return
	}

	// Update yang diinput saja
	if input.Rating != 0 {
		rev.Rating = input.Rating
	}
	if input.Content != "" {
		rev.Content = input.Content
	}

	rev.UpdatedAt = time.Now()

	if err := db.Save(&rev).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Failed to update review", http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("review"), http.StatusOK, rev))
}

type ReviewsResponse struct {
	models.Review
	Username   string `json:"username"`
	PhoneModel string `json:"phone_model"`
}

// Get all reviews data
// @Summary Get all reviews (ADMIN ONLY)
// @Description Get all Reviews
// @Tags Reviews
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []map[string]any
// @Router /reviews [get]
func GetAllReviews(c *gin.Context) {
	var data []ReviewsResponse

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Table("reviews").
		Select("reviews.*, users.username as username, phones.model as phone_model").
		Joins("join users on reviews.user_id = users.id").
		Joins("join phones on reviews.phone_id = phones.id").
		Scan(&data).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, data))
}

type Reviews struct {
	models.Review
	Username string `json:"username"`
}

// Get reviews data by Phone data ID godoc
// @Summary Get reviews data by Phone id. (PUBLIC)
// @Description Get all Reviews data by phone id. if reviews data is empty, the review data will not be displayed
// @Tags Phones
// @Produce json
// @Param id path string true "Phone id"
// @Success 200 {object} []models.Phone
// @Router /phones/{id}/reviews [get]
func GetReviewsDataByPhoneId(c *gin.Context) {
	var reviews []Reviews

	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	if err := db.Table("reviews").
		Select("reviews.*, users.username as username").
		Joins("join users on reviews.user_id = users.id").
		Where("reviews.phone_id = ?", id).
		Find(&reviews).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(err.Error(), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, reviews))
}

// Delete account
// @Summary Delete review by ID
// @Description This route will delete review data, based on the review ID and will only be able to delete data related to the logged in user (user ID is taken from the JWT token)
// @Tags Reviews
// @Produce json
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Param id path string true "review id"
// @Security BearerToken
// @Success 200 {object} []models.Review
// @Router /reviews/{id} [delete]
func DeleteReviewById(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	reviewID := c.Param("id")

	var review []models.Review
	// cek apakah review dengan id tsb ada
	if err := db.Where("id = ?", reviewID).Find(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	if len(review) == 0 {
		msg := "data review tidak ditemukan"
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(msg, http.StatusNotFound, nil))
		return
	}

	if err := db.Model(&models.Review{}).Where("id = ?", reviewID).Delete(&review).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}
	c.JSON(http.StatusOK,
		utils.ResponseJSON(lib.MsgDeleted("review"), http.StatusOK, nil))
}
