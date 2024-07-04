package models

import (
	"time"
)

type Brand struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LogoURL     string    `gorm:"not null" json:"logo_url"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Phones      []Phone   `gorm:"foreignKey:BrandID;" json:"phones,omitempty"`
}
