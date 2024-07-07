package models

import (
	"errors"
	"final-project/lib"
	"final-project/utils/token"
	"fmt"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	RoleID    uint      `gorm:"not null" json:"-"`
	Role      Role      `gorm:"foreignKey:RoleID;default:1;" json:"-"`
	Profiles  []Profile `gorm:"foreignKey:UserID;constraint:onDelete:CASCADE" json:"profiles,omitempty"`
	Reviews   []Review  `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
	Comments  []Comment `gorm:"foreignKey:UserID" json:"-"`
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username, email string, password string, db *gorm.DB) (string, error) {

	var err error

	u := User{}

	err = db.Model(User{}).Where("username = ? OR email = ?", username, email).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	fmt.Println("Token:", token)

	if err != nil {
		return "", err
	}

	return token, nil
}

func HashPassword(password_text string) (string, error) {
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(password_text), bcrypt.DefaultCost)
	if errPassword != nil {
		return "", errPassword
	}

	return string(hashedPassword), nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	// cek username dat email duplikat
	if checkDuplicateUsername(db, u.Username) {
		return &User{}, errors.New(lib.MsgAlreadyExist("username"))
	}
	if checkDuplicateEmail(db, u.Email) {
		return &User{}, errors.New(lib.MsgAlreadyExist("email"))
	}

	//turn password into hash
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return &User{}, errPassword
	}
	u.Password = string(hashedPassword)
	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	var err error = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func checkDuplicateUsername(db *gorm.DB, username string) bool {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	return err == nil
}

func checkDuplicateEmail(db *gorm.DB, email string) bool {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	return err == nil
}
