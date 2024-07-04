package models

import (
	"time"
)

type Specification struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Network           string    `gorm:"not null" json:"network"`
	OperatingSystem   string    `gorm:"not null" json:"operating_system"`
	Storage           uint      `gorm:"not null" json:"storage"`
	Memory            uint      `gorm:"not null" json:"memory"`
	Camera            uint      `gorm:"not null" json:"camera"`
	Battery           string    `gorm:"not null" json:"battery"`
	AdditionalFeature string    `json:"additional_feature"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	PhoneID           uint      `gorm:"not null" json:"phone_id"`
}
