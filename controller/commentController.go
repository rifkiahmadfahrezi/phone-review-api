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

type commentInput struct {
	Content string `json:"content" binding:"required"`
}

// Create Comment godoc
// @Summary Create Comment on a review
// @Description Creating a comment on a phone review, (user ID is taken from the JWT token)
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Tags Reviews
// @Security BearerToken
// @Param id path string true "review id"
// @Param Body body commentInput true "example JSON body to create a comment, user_id is taken from the authorization token"
// @Produce json
// @Success 200 {object} models.Comment
// @Router /reviews/{id}/comments [post]
func CreateComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input commentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	// cek apakah data review ada atau tdk
	// jika tdk ada maka user tidak bisa komen
	reviewIDstr := c.Param("id")

	reviewID, err := strconv.Atoi(reviewIDstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	var rev []models.Review
	if err := db.Where("id = ?", reviewID).First(&rev).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	if len(rev) == 0 {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("review"), http.StatusNotFound, nil))
		return
	}

	// ambil user id
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Gagal mengambil id user", http.StatusInternalServerError, nil))
		return
	}

	comment_data := models.Comment{
		Content:  input.Content,
		UserID:   userID,
		ReviewID: uint(reviewID),
	}

	db.Create(&comment_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("comment"), http.StatusOK, comment_data))
}

// Update Comment godoc
// @Summary Update Comment on a review
// @Description This route will update comment data, will only be able to update data related to the logged in user (user ID is taken from the JWT token)
// @Tags Comments
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "comment id"
// @Param Body body commentInput true "example JSON body to update a comment, user_id is taken from the authorization token"
// @Produce json
// @Success 200 {object} models.Comment
// @Router /comments/{id} [put]
func UpdateComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input commentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	commentID := c.Param("id")

	// cek komen berdasarkan ID
	if err := db.Where("id = ?", commentID).First(&models.Comment{}).Error; err != nil {
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

	// cek jika user sudah memberi comment atau blm
	var rev models.Comment
	if err := db.Where("user_id = ?", userID).Find(&rev).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON("user ini belum memberikan comment", http.StatusNotFound, nil))
		return
	}

	// cek data comment ada atau tdk
	if err := db.Where("id = ? AND user_id = ?", commentID, userID).Find(&rev).Error; err != nil {
		msg := fmt.Sprintf("comment dgn ID (%s) dari user (%d) tidak ditemukan", commentID, userID)
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(msg, http.StatusNotFound, nil))
		return
	}

	// Update yg diinput saja
	if input.Content != "" {
		rev.Content = input.Content
	}

	var updated_data commentInput
	updated_data.Content = rev.Content

	rev.UpdatedAt = time.Now()

	if err := db.Model(&models.Comment{}).Where("user_id = ? AND id = ?", userID, commentID).Updates(&updated_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Failed to update review", http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("review"), http.StatusOK, rev))
}

// Get comments data by Review data ID godoc
// @Summary Get comments data by Review id.
// @Description Get all Comments data by review id.
// @Tags Reviews
// @Produce json
// @Param id path string true "Review id"
// @Success 200 {object} []models.Review
// @Router /reviews/{id}/comments [get]
func GetCommentsDataByReviewId(c *gin.Context) {
	var reviews []models.Review

	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	if err := db.Preload("Comments.User").Preload("User").Find(&reviews, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("review"), http.StatusNotFound, nil))
		return
	}

	var response []map[string]interface{}
	for _, review := range reviews {
		r := map[string]interface{}{
			"id":         review.ID,
			"rating":     review.Rating,
			"content":    review.Content,
			"username":   review.User.Username,
			"created_at": review.CreatedAt,
			"updated_at": review.UpdatedAt,
			"user_id":    review.UserID,
			"phone_id":   review.PhoneID,
		}

		var comments []map[string]interface{}
		for _, comment := range review.Comments {
			c := map[string]interface{}{
				"id":         comment.ID,
				"content":    comment.Content,
				"username":   comment.User.Username,
				"created_at": comment.CreatedAt,
				"updated_at": comment.UpdatedAt,
				"user_id":    comment.UserID,
				"review_id":  comment.ReviewID,
			}
			comments = append(comments, c)
		}

		r["comments"] = comments
		response = append(response, r)
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, response))
}

// Delete Comment by id  godoc
// @Summary Delete Comment by id .
// @Description Delete a Comment by id, (user ID is taken from the JWT token so only loogged in user can delete its own coment)
// @Tags Comments
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Comment id"
// @Success 200 {object} map[string]boolean
// @Router /comments/{id} [delete]
func DeleteCommentByID(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var comment_data models.Comment
	if err := db.Where("id = ?", c.Param("id")).First(&comment_data).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.ErrMsgNotFound("comment"), http.StatusBadRequest, nil))
		return
	}

	if err := db.Delete(&comment_data).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgDeleted("comment"), http.StatusOK, nil))
}

// Get all phone comments
// @Summary Get all Phones comments. (ADMIN ONLY)
// @Description Get all of comments data.
// @Tags Comments
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.Comment
// @Router /comments [get]
func GetAllCommentData(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var comments_data []models.Comment

	err := db.Preload("User").Find(&comments_data).Error
	if err != nil {
		emptydata := make([]string, 0)
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, emptydata))
		return
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, comments_data))
}
