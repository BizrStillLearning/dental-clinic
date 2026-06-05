package models

import (
	"time"

	"gorm.io/gorm"
)

type Obat struct {
	ID        uint   `gorm:"primaryKey"`
	Nama      string `gorm:"type:varchar(100);not null"`
	Harga     int    `gorm:"not null"`
	Stok      int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
