package models

import (
	"time"
)

type Review struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Rating    uint      `gorm:"not null" json:"rating"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	PhoneID   uint      `gorm:"not null" json:"phone_id"`
	Comments  []Comment `gorm:"foreignKey:ReviewID;constraint:onDelete:CASCADE" json:"comments,omitempty"`
	User      User      `json:"user"`
}
