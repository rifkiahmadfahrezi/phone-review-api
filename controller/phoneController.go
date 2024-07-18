package controller

import (
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type phoneInput struct {
	Model       string    `json:"model" bind:"required"`
	Price       uint      `json:"price" bind:"required"`
	ImageURL    string    `json:"image_url" bind:"required"`
	ReleaseDate time.Time `json:"release_date"`
	BrandID     uint      `json:"brand_id" bind:"required"`
}

type phoneUpdate struct {
	Model       string    `json:"model"`
	Price       uint      `json:"price"`
	ImageURL    string    `json:"image_url" bind:"url"`
	ReleaseDate time.Time `json:"release_date"`
	BrandID     uint      `json:"brand_id"`
}

type PhonesCompleteResponse struct {
	PhoneID     int       `json:"phone_id"`
	BrandID     int       `json:"brand_id"`
	BrandName   string    `json:"brand_name"`
	PhoneImage  string    `json:"phone_image"`
	PhoneModel  string    `json:"phone_model"`
	FullName    string    `json:"full_name"`
	AVGRating   float64   `json:"avg_rating"`
	Price       float64   `json:"price"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Get all phone data
// @Summary Get all Phones data. (PUBLIC)
// @Description Get a list of Phone.
// @Tags Phones
// @Produce json
// @Success 200 {object} []models.Phone
// @Router /phones [get]
func GetAllPhoneData(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	// apply filtering
	searchKeyword := c.Query("search")
	sort := c.Query("sort")

	query := db.Model(&models.Phone{})

	if searchKeyword != "" {
		q := fmt.Sprintf("%%%s%%", searchKeyword)
		query.Where("name LIKE ?", q)
	}

	switch strings.ToLower(sort) {
	case "desc":
		query.Order("phone_id DESC")
	default:
		query.Order("phone_id ASC")
	}

	var phones_data []PhonesCompleteResponse
	if err := query.Table("phones").
		Select(`brands.name as brand_name, 
            phones.id as phone_id,
            brands.id as brand_id,
            phones.image_url as phone_image, 
            phones.model as phone_model, 
            brands.name || ' ' || phones.model as full_name, 
            ROUND(AVG(reviews.rating), 2) as avg_rating,
            phones.price, phones.release_date, phones.created_at, phones.updated_at`).
		Joins("join reviews on phones.id = reviews.phone_id").
		Joins("join brands on brands.id = phones.brand_id").
		Group("brands.name, phones.id, brands.id, phones.image_url, phones.model, phones.price, phones.release_date, phones.created_at, phones.updated_at").
		Scan(&phones_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	// validsi jika data tidak ditemukan
	if searchKeyword != "" || sort != "" {
		if len(phones_data) == 0 {
			emptydata := make([]string, 0)
			c.JSON(http.StatusNotFound,
				utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusNotFound, emptydata))
			return
		}
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, phones_data))
}

// Create New Phone godoc
// @Summary Create New Phone (ADMIN ONLY)
// @Description Creating a new Phone data, only account with role admin can accsess this route
// @Tags Phones
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body phoneInput true "example JSON body to create a new Phone, sample release_data format = 2023-09-22T00:00:00Z"
// @Produce json
// @Success 200 {object} models.Phone
// @Router /phones [post]
func CreatePhoneData(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input phoneInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	// cek apakah brand id dari input ada di tabel brand
	var brand models.Brand
	if err := db.Where("id = ?", input.BrandID).First(&brand).Error; err != nil {
		idStr := strconv.Itoa(int(input.BrandID))
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("brand dengan id "+idStr), http.StatusNotFound, nil))
		return
	}

	if !isPhoneInputDataValid(c, input) {
		return
	}

	phone_data := models.Phone{
		Model:       input.Model,
		ReleaseDate: input.ReleaseDate,
		Price:       input.Price,
		ImageURL:    input.ImageURL,
		BrandID:     input.BrandID,
	}

	// Add error handling for db.Create
	if err := db.Create(&phone_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON(lib.MsgAdded("phone"), http.StatusOK, phone_data))
}

// GetPhoneById godoc
// @Summary Get Phone. (PUBLIC)
// @Description Get a Phone by id.
// @Tags Phones
// @Produce json
// @Param id path string true "phone id"
// @Success 200 {object} models.Phone
// @Router /phones/{id} [get]
func GetPhoneById(c *gin.Context) {
	var phone []PhonesCompleteResponse
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Table("phones").
		Select(`brands.name as brand_name, 
				phones.id as phone_id,
				brands.id as brand_id,
				phones.image_url as phone_image, 
				phones.model as phone_model, 
				brands.name || ' ' || phones.model as full_name, 
				ROUND(AVG(reviews.rating), 2) as avg_rating,
				phones.price, phones.release_date, phones.created_at, phones.updated_at`).
		Joins("join reviews on phones.id = reviews.phone_id").
		Joins("join brands on brands.id = phones.brand_id").
		Group("brands.name, phones.id, brands.id, phones.image_url, phones.model, phones.price, phones.release_date, phones.created_at, phones.updated_at").
		Where("phones.id = ? ", c.Param("id")).
		Scan(&phone).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	if len(phone) == 0 {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, phone))
}

// Update Phone data godoc
// @Summary Update Phone data. (ADMIN ONLY)
// @Description Update Phone data by id, only account with role admin can access this route
// @Tags Phones
// @Produce json
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "Phone id"
// @Param Body body phoneInput true "Example JSON body to update Phone data, sample release_data format = 2023-09-22T00:00:00Z"
// @Success 200 {object} models.Phone
// @Router /phones/{id} [put]
func UpdatePhoneData(c *gin.Context) {

	var input phoneUpdate
	db := c.MustGet("db").(*gorm.DB)

	// cek data phone dengan id tsb
	var phone models.Phone
	if err := db.Where("id = ?", c.Param("id")).First(&phone).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusNotFound, nil))
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

	// jika brand id diganti maka cek apakah brand id dari input ada di tabel brand
	var brand models.Brand
	if input.BrandID != 0 {
		if err := db.Where("id = ?", input.BrandID).First(&brand).Error; err != nil {
			idStr := strconv.Itoa(int(input.BrandID))
			c.JSON(http.StatusNotFound,
				utils.ResponseJSON(lib.ErrMsgNotFound("brand dengan id "+idStr), http.StatusNotFound, nil))
			return
		}
	}

	if input.ImageURL != "" && !utils.IsValidUrl(input.ImageURL) {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgValidUrl("logo_url"), http.StatusBadRequest, nil))
		return
	}

	var updated_data models.Phone
	updated_data.Model = input.Model
	updated_data.ImageURL = input.ImageURL
	updated_data.ReleaseDate = input.ReleaseDate
	updated_data.Price = input.Price
	if input.BrandID != 0 {
		updated_data.BrandID = input.BrandID
	}
	updated_data.BrandID = brand.ID
	updated_data.UpdatedAt = time.Now()

	// update ke tabel
	db.Model(&phone).Updates(updated_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("phone"), http.StatusOK, phone))
}

// Delete Phone by id  godoc
// @Summary Delete Phone by id . (ADMIN ONLY)
// @Description Delete a Phone by id, only account with role admin can access this route
// @Tags Phones
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Phone id"
// @Success 200 {object} map[string]boolean
// @Router /phones/{id} [delete]
func DeletePhoneData(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	var phone_data models.Phone
	if err := db.Where("id = ?", c.Param("id")).First(&phone_data).Error; err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusBadRequest, nil))
		return
	}

	db.Delete(&phone_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgDeleted("phone"), http.StatusOK, nil))
}

func isPhoneInputDataValid(c *gin.Context, data phoneInput) bool {
	dataErr := []string{}

	if data.ImageURL == "" {
		dataErr = append(dataErr, "image_url")
	}
	if data.Model == "" {
		dataErr = append(dataErr, "model")
	}
	if data.Price == 0 {
		dataErr = append(dataErr, "price")
	}

	if data.ReleaseDate.IsZero() {
		dataErr = append(dataErr, "release_date")
	}

	if len(dataErr) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgRequired(dataErr...), http.StatusBadRequest, nil))
		return false
	}

	if !utils.IsValidUrl(data.ImageURL) {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgValidUrl("image_url"), http.StatusBadRequest, nil))
		return false
	}

	return true
}

// Get phones specification by phone ID godoc
// @Summary Get specification data by Phone id. (PUBLIC)
// @Description Get Phone specifiction data by phone id. if phone's specification data empty, the spec data will not be displayed
// @Tags Phones
// @Produce json
// @Param id path string true "Phone id"
// @Success 200 {object} []models.Phone
// @Router /phones/{id}/specification [get]
func GetPhonesSpecByPhoneId(c *gin.Context) {
	var phones []models.Phone

	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	if err := db.Preload("Specifications").Find(&phones, id).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, phones))
}
