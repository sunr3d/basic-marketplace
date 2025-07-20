package models

import "time"

type Adv struct {
	ID          uint `gorm:"primaryKey"`
	Title       string `gorm:"size:35;not null"`
	Description string `gorm:"size:100;not null"`
	ImageURL    string `gorm:"size:255"`
	Price       float64 `gorm:"not null"`
	OwnerID     uint `gorm:"not null"`
	CreatedAt   time.Time
}
