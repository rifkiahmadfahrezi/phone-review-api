package controller

import (
	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type specificationInput struct {
	Network           string `json:"network" binding:"required"`
	OperatingSystem   string `json:"operating_system" binding:"required"`
	Storage           uint   `json:"storage" binding:"required"`
	Memory            uint   `json:"memory" binding:"required"`
	Camera            uint   `json:"camera" binding:"required"`
	Battery           string `json:"battery" binding:"required"`
	AdditionalFeature string `json:"additional_feature"`
}

type specificationUpdate struct {
	Network           string `json:"network"`
	OperatingSystem   string `json:"operating_system"`
	Storage           uint   `json:"storage"`
	Memory            uint   `json:"memory"`
	Camera            uint   `json:"camera"`
	Battery           string `json:"battery"`
	AdditionalFeature string `json:"additional_feature"`
}

// Create Specification for phone godoc
// @Summary Create Specification for phone
// @Description Creating a specification data for phone, only admin can access this route
// @Tags Phones
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "Phone id"
// @Param Body body specificationInput true "example JSON body to create a specification for phone"
// @Produce json
// @Success 200 {object} models.Specification
// @Router /phones/{id}/specification [post]
func CreateSpecification(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input specificationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	phoneID := c.Param("id")

	// cek data specification ada berdasarkan phone_id
	var specification []models.Specification
	if err := db.Where("phone_id = ?", phoneID).Find(&specification).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(specification) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON("specification untuk phone ini sudah tersedia", http.StatusBadRequest, nil))
		return
	}

	// cek data phone ada berdasarkan ID
	var phone []models.Phone
	if err := db.Where("id = ?", phoneID).Find(&phone).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(phone) == 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.ErrMsgNotFound("phone"), http.StatusBadRequest, nil))
		return
	}

	specification_data := models.Specification{
		Network:           input.Network,
		OperatingSystem:   input.OperatingSystem,
		Storage:           input.Storage,
		Memory:            input.Memory,
		Camera:            input.Camera,
		Battery:           input.Battery,
		AdditionalFeature: input.AdditionalFeature,
	}

	if err := db.Create(&specification_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("specification"), http.StatusOK, specification_data))
}

// Update Specification for phone godoc
// @Summary Update Specification for phone
// @Description Creating a specification data for phone, only admin can access this route
// @Tags Phones
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "phone id"
// @Param Body body specificationInput true "example JSON body to update a specification for phone"
// @Produce json
// @Success 200 {object} models.Specification
// @Router /phones/{id}/specification [put]
func UpdateSpecification(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input specificationUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	phoneID := c.Param("id")

	// cek data phone ada
	var spec models.Specification
	if err := db.Where("phone_id = ?", phoneID).First(&spec).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	// Update yg diinput saja
	if input.Network != "" {
		spec.Network = input.Network
	}
	if input.OperatingSystem != "" {
		spec.OperatingSystem = input.OperatingSystem
	}
	if input.Storage != 0 {
		spec.Storage = input.Storage
	}
	if input.Memory != 0 {
		spec.Memory = input.Memory
	}
	if input.Camera != 0 {
		spec.Camera = input.Camera
	}
	if input.Battery != "" {
		spec.Battery = input.Battery
	}
	if input.AdditionalFeature != "" {
		spec.AdditionalFeature = input.AdditionalFeature
	}
	spec.UpdatedAt = time.Now()

	if err := db.Save(&spec).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON("Failed to update specification", http.StatusInternalServerError, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("specification"), http.StatusOK, spec))
}
