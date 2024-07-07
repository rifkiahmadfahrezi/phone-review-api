package middleware

import (
	"net/http"
	"strings"

	"final-project/lib"
	"final-project/models"
	"final-project/utils"
	"final-project/utils/token"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userRole struct {
	RoleName string `json:"role_name"`
	RoleID   int    `json:"role_id"`
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, "token tidak valid")
			return
		}
		c.Next()
	}
}

func RoleMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// validasi token
		if err := token.TokenValid(c); err != nil {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON("Token tidak valid", http.StatusBadRequest, nil))
			return
		}

		userID, err := token.ExtractTokenID(c)

		if err != nil {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
			return
		}

		// cek data user
		var user []models.User
		if err := db.Where("id = ?", userID).Find(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest,
				utils.ResponseJSON(err.Error(), http.StatusBadRequest, nil))
			return
		}

		// jika data tidak ditemukan
		if len(user) == 0 {
			c.JSON(http.StatusNotFound,
				utils.ResponseJSON(lib.ErrMsgNotFound("user"), http.StatusNotFound, nil))
			return
		}

		// ambil nama role berdasarkan user id
		var userRole userRole
		if err := db.Table("roles").
			Select(`roles.name as role_name,
						roles.id as role_id`).
			Joins("join users on users.role_id = roles.id").
			Where("users.id = ?", user[0].ID).Scan(&userRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError,
				utils.ResponseJSON(err.Error(), http.StatusInternalServerError, nil))
			return
		}

		if strings.ToLower(userRole.RoleName) != "admin" {
			c.JSON(http.StatusForbidden,
				utils.ResponseJSON("anda tidak bisa mengakses route ini", http.StatusForbidden, nil))
			c.Abort()
			return
		}

		// if user[0].RoleID != 2 { // 2 == admin
		// 	c.JSON(http.StatusForbidden,
		// 		utils.ResponseJSON("anda tidak bisa mengakses route ini", http.StatusForbidden, nil))
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}
