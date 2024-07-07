package models

import (
	"time"
)

type Profile struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Biodata   string     `json:"biodata"`
	ImageURL  string     `gorm:"not null;default:'https://i.ibb.co.com/XjBCcsL/user-1.png'" json:"image_url"`
	FullName  string     `gorm:"not null" json:"full_name"`
	Birthday  *time.Time `json:"birthday"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	UserID    uint       `gorm:"not null" json:"user_id"`
}
