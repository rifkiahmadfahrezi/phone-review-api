package controller

import (
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type brandInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoUrl     string `json:"logo_url" binding:"required"`
}

type brandUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoUrl     string `json:"logo_url"`
}

// Get all phone brands
// @Summary Get all Phones brands.
// @Description Get a list of Phone brands.
// @Tags Brands
// @Produce json
// @Success 200 {object} []models.Brand
// @Router /brands [get]
func GetAllBrandData(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var brands_data []models.Brand

	// apply filtering
	searchKeyword := c.Query("search")
	sort := c.Query("sort")

	query := db.Model(&models.Brand{})

	if searchKeyword != "" {
		q := fmt.Sprintf("%%%s%%", searchKeyword)
		query.Where("name LIKE ?", q)
	}

	switch strings.ToLower(sort) {
	case "desc":
		query.Order("id DESC")
	default:
		query.Order("id ASC")
	}

	err := query.Select("id", "logo_url", "name", "description", "created_at", "updated_at").Find(&brands_data).Error
	if err != nil {
		emptydata := make([]string, 0)
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, emptydata))
		return
	}

	// validsi jika data tidak ada
	if searchKeyword != "" || sort != "" {
		if len(brands_data) == 0 {
			emptydata := make([]string, 0)
			c.JSON(http.StatusNotFound,
				utils.ResponseJSON(lib.ErrMsgNotFound("brands"), http.StatusNotFound, emptydata))
			return
		}
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, brands_data))
}

// Create New Brand godoc
// @Summary Create New Brand
// @Description Creating a new Brand data, only account with role admin can access this route
// @Tags Brands
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body brandInput true "example JSON body to create a new Brand"
// @Produce json
// @Success 200 {object} models.Brand
// @Router /brands [post]
func CreateBrand(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input brandInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	// cek nama brand sudah ada atau belum
	var brand []models.Brand
	if err := db.Where("name = ?", input.Name).Find(&brand).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(brand) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgAlreadyExist("brand"), http.StatusBadRequest, nil))
		return
	}

	// validasi data brand
	if !isBrandDataInputValid(c, input) {
		return
	}

	brand_data := models.Brand{
		Name:        input.Name,
		Description: input.Description,
		LogoURL:     input.LogoUrl,
	}

	db.Create(&brand_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("brand"), http.StatusOK, brand_data))
}

// Get phones data by Brand data ID godoc
// @Summary Get phones data by Brand id.
// @Description Get all Phones data by brand id.
// @Tags Brands
// @Produce json
// @Param id path string true "Brand id"
// @Success 200 {object} []models.Brand
// @Router /brands/{id}/phones [get]
func GetPhonesDataByBrandId(c *gin.Context) {
	var brands []models.Brand

	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	if err := db.Preload("Phones").Find(&brands, id).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("brand"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, brands))
}

// Get  Brand data by ID godoc
// @Summary Get single brand data by ID.
// @Description Get Brand data by ID.
// @Tags Brands
// @Produce json
// @Param id path string true "Brand id"
// @Success 200 {object} []models.Brand
// @Router /brands/{id} [get]
func GetBrandById(c *gin.Context) {
	var brands_data []models.Brand

	db := c.MustGet("db").(*gorm.DB)
	brandID := c.Param("id")

	if err := db.Select("id", "logo_url", "name", "description", "created_at", "updated_at").First(&brands_data, brandID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}

	// jika data tidak ditemukan
	if len(brands_data) == 0 {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("brand"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, brands_data))
}

// Update Brand data godoc
// @Summary Update Brand data.
// @Description Update Brand data by id, only account with role admin can access this route
// @Tags Brands
// @Produce json
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "Brand id"
// @Param Body body brandInput true "tExample JSON body to update Brand data"
// @Success 200 {object} models.Brand
// @Router /brands/{id} [put]
func UpdateBrand(c *gin.Context) {

	var input brandUpdate
	db := c.MustGet("db").(*gorm.DB)

	// cek data brand dengan id tsb
	var brand models.Brand
	if err := db.Where("id = ?", c.Param("id")).First(&brand).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("brand"), http.StatusNotFound, nil))
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

	if input.LogoUrl != "" && !utils.IsValidUrl(input.LogoUrl) {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgValidUrl("logo_url"), http.StatusBadRequest, nil))
		return
	}

	var updated_data models.Brand
	updated_data.Name = input.Name
	updated_data.LogoURL = input.LogoUrl
	updated_data.Description = input.Description
	updated_data.UpdatedAt = time.Now()

	// cek nama brand sudah ada atau belum
	var brand_exist []models.Brand
	if err := db.Where("name = ?", updated_data.Name).Find(&brand_exist).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(brand_exist) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgAlreadyExist("brand"), http.StatusBadRequest, nil))
		return
	}

	// update ke tabel
	db.Model(&brand).Updates(updated_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("brand"), http.StatusOK, brand))
}

// Delete Brand by id  godoc
// @Summary Delete Brand by id .
// @Description Delete a Brand by id
// @Tags Brands
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Brand id"
// @Success 200 {object} map[string]boolean
// @Router /brands/{id} [delete]
func DeleteBrandByID(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	var brand_data models.Brand
	if err := db.Where("id = ?", c.Param("id")).First(&brand_data).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.ErrMsgNotFound("brand"), http.StatusBadRequest, nil))
		return
	}

	if err := db.Delete(&brand_data).Error; err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			errmsg := fmt.Sprintf("brand %s tidak bisa dihapus karena sudah terkait dengan data phone", brand_data.Name)
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errmsg, http.StatusBadRequest, nil))
			return
		}
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgDeleted("brand"), http.StatusOK, nil))
}

func isBrandDataInputValid(c *gin.Context, data brandInput) bool {

	if data.Name == "" && data.LogoUrl == "" {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgRequired("name", "logo_url"), http.StatusBadRequest, nil))
		return false
	}

	if data.Name == "" {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgRequired("name"), http.StatusBadRequest, nil))
		return false
	}

	if data.LogoUrl == "" {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgRequired("logo_url"), http.StatusBadRequest, nil))
		return false
	}

	if !utils.IsValidUrl(data.LogoUrl) {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgValidUrl("logo_url"), http.StatusBadRequest, nil))
		return false
	}

	return true
}
