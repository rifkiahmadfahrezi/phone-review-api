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

type roleInput struct {
	Name string `json:"name" bind:"requred"`
}

// Get all roles
// @Summary Get all  roles.
// @Description Get a list of user's roles. only admin can access this route
// @Tags Roles
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} []models.Role
// @Router /roles [get]
func GetAllRoleData(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var roles_data []models.Role

	// apply filtering
	searchKeyword := c.Query("search")
	sort := c.Query("sort")

	query := db.Model(&models.Role{})

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

	err := query.Select("id", "name", "created_at", "updated_at").Find(&roles_data).Error
	if err != nil {
		emptydata := make([]string, 0)
		c.JSON(http.StatusInternalServerError,
			utils.ResponseJSON(err.Error(), http.StatusInternalServerError, emptydata))
		return
	}

	// validsi jika data tidak ditemukan
	if searchKeyword != "" || sort != "" {
		if len(roles_data) == 0 {
			emptydata := make([]string, 0)
			c.JSON(http.StatusNotFound,
				utils.ResponseJSON(lib.ErrMsgNotFound("role"), http.StatusNotFound, emptydata))
			return
		}
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, roles_data))
}

// Get role by ID
// @Summary Get role by ID.
// @Description Get a role data by id. only admin can access this route
// @Tags Roles
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Role id"
// @Success 200 {object} []models.Role
// @Router /roles/{id} [get]
func GetRoleDataByID(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var roles_data []models.Role

	id := c.Param("id")

	err := db.Select("id", "name", "created_at", "updated_at").Where("id = ?", id).Find(&roles_data).Error
	if err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("role"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK,
		utils.ResponseJSON("", http.StatusOK, roles_data))
}

// Create New Role godoc
// @Summary Create New Role
// @Description Creating a new Role data, only admin can access this route
// @Tags Roles
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body roleInput true "example JSON body to create a new Role"
// @Produce json
// @Success 200 {object} models.Role
// @Router /roles [post]
func CreateRole(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input roleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.CustomBindError(err)
		if errorMessage != "" {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errorMessage, http.StatusBadRequest, nil))
		}
		return
	}

	// cek nama role sudah ada atau belum
	var role []models.Role
	if err := db.Where("name = ?", input.Name).Find(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(role) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgAlreadyExist("role"), http.StatusBadRequest, nil))
		return
	}

	role_data := models.Role{
		Name: input.Name,
	}

	db.Create(&role_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgAdded("role"), http.StatusOK, role_data))
}

// Delete Role by id  godoc
// @Summary Delete Role by id .
// @Description Delete a Role by id, only admin can access this route
// @Tags Roles
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Role id"
// @Success 200 {object} map[string]boolean
// @Router /roles/{id} [delete]
func DeleteRoleByID(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	var role_data models.Role

	id := c.Param("id")

	if id == "1" || id == "2" {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON("role 'admin' dan 'user' tidak bisa dihapus", http.StatusBadRequest, nil))
		return
	}

	if err := db.Where("id = ?", id).First(&role_data).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.ErrMsgNotFound("role"), http.StatusBadRequest, nil))
		return
	}

	if err := db.Delete(&role_data).Error; err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			errmsg := fmt.Sprintf("role %s tidak bisa dihapus karena sudah terkait dengan data user", role_data.Name)
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(errmsg, http.StatusBadRequest, nil))
			return
		}
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgDeleted("role"), http.StatusOK, nil))
}

// Update Role data godoc
// @Summary Update Role data.
// @Description Update Role data by id, only account with role admin can access this route
// @Tags Roles
// @Produce json
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "Role id"
// @Param Body body roleInput true "tExample JSON body to update Role data"
// @Success 200 {object} models.Role
// @Router /roles/{id} [put]
func UpdateRole(c *gin.Context) {

	var input roleInput
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	if id == "1" || id == "2" {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON("role 'admin' dan 'user' tidak bisa diupdate", http.StatusBadRequest, nil))
		return
	}

	// cek data role dengan id tsb
	var role models.Role
	if err := db.Where("id = ?", id).First(&role).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("role"), http.StatusNotFound, nil))
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

	var updated_data models.Role
	updated_data.Name = input.Name
	updated_data.UpdatedAt = time.Now()

	// cek nama role sudah ada atau belum
	var role_exist []models.Role
	if err := db.Where("name = ?", updated_data.Name).Find(&role_exist).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
		return
	}

	if len(role_exist) > 0 {
		c.JSON(http.StatusBadRequest,
			utils.ResponseJSON(lib.MsgAlreadyExist("role"), http.StatusBadRequest, nil))
		return
	}

	// update ke tabel
	db.Model(&role).Updates(updated_data)

	c.JSON(http.StatusOK, utils.ResponseJSON(lib.MsgUpdated("role"), http.StatusOK, role))
}

// Get users data by Role data ID godoc
// @Summary Get users data by Role id.
// @Description Get all Users data by role id. only admin can access this route
// @Tags Roles
// @Produce json
// @Param Authorization header string true "Authorization : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "Role id"
// @Success 200 {object} []models.Role
// @Router /roles/{id}/users [get]
func GetUsersDataByRoleId(c *gin.Context) {
	var roles []models.Role

	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	if err := db.Preload("Users").Find(&roles, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			utils.ResponseJSON(lib.ErrMsgNotFound("role"), http.StatusNotFound, nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseJSON("", http.StatusOK, roles))
}
