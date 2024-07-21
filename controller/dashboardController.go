package controller

import (
	"final-project/models"
	"final-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get number of all data
// @Summary number of all data (ADMIN ONLY)
// @Description Get a number of all data
// @Tags Dashboard
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []map[string]any
// @Router /dashboard/all-count-data [get]
func GetAllDataCount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var numUsers, numAdmins, numPhones, numBrands, numReviews int64

	if err := db.Model(&models.User{}).Where("role_id = ?", 1).Count(&numUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}
	if err := db.Model(&models.User{}).Where("role_id = ?", 2).Count(&numAdmins).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}
	if err := db.Model(&models.Phone{}).Count(&numPhones).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}
	if err := db.Model(&models.Brand{}).Count(&numBrands).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}
	if err := db.Model(&models.Review{}).Count(&numReviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	response := []map[string]any{
		{
			"name":        "users",
			"num_of_data": int(numUsers),
		},
		{
			"name":        "admins",
			"num_of_data": int(numAdmins),
		},
		{
			"name":        "phones",
			"num_of_data": int(numPhones),
		},
		{
			"name":        "brands",
			"num_of_data": int(numBrands),
		},
		{
			"name":        "reviews",
			"num_of_data": int(numReviews),
		},
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, response))
}
