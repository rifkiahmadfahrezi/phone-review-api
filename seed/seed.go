package seed

import (
	"final-project/models"

	"gorm.io/gorm"
)

// Load initial data (role, user)
func Load(db *gorm.DB) {

	// insert initial role
	role_data := []models.Role{
		{Name: "user"},
		{Name: "admin"},
	}

	for _, role := range role_data {
		db.FirstOrCreate(&role, models.Role{Name: role.Name})
	}

	//insert initial user & admin
	user_data := []models.User{
		{
			Username: "user",
			Email:    "user@gmail.com",
			Password: "$2a$10$ijrkqTmYqvCdEmR/CuIJ4eNH0Br6.CDxGBoJytqE7fxuVWihqeaoO", // user123
			RoleID:   1,                                                              // user
		},
		{
			Username: "admin",
			Email:    "admin@gmail.com",
			Password: "$2a$10$g/i0uCYW4smgD5ccBo0TwOrC5JOxfdFMTHSOnDYyIxmY5Zhd8.zpa", // admin
			RoleID:   2,                                                              // admin
		},
	}

	for _, user := range user_data {
		db.FirstOrCreate(&user, models.User{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
			RoleID:   user.RoleID,
		})
	}

}
