package models

import (
	"time"
)

type Phone struct {
	ID             uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	ImageURL       string          `gorm:"not null" json:"image_url"`
	Model          string          `gorm:"not null" json:"model"`
	Price          uint            `gorm:"not null" json:"price"`
	ReleaseDate    time.Time       `gorm:"not null" json:"release_date"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	BrandID        uint            `gorm:"not null" json:"-"`
	Reviews        []Review        `gorm:"foreignKey:PhoneID" json:"reviews,omitempty"`
	Specifications []Specification `gorm:"foreignKey:PhoneID;constraint:onDelete:CASCADE" json:"specification,omitempty"`
}

type PhoneWithBrand struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ImageURL    string    `gorm:"not null" json:"image_url"`
	Model       string    `gorm:"not null" json:"model"`
	Price       uint      `gorm:"not null" json:"price"`
	ReleaseDate time.Time `gorm:"not null" json:"release_date"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Brands      Brand     `gorm:"many2many:PhoneID;" json:"brand"`
}
