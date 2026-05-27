package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	DokterID  uint   `gorm:"not null"`
	Tanggal   string `gorm:"type:date;not null"`
	Keluhan   string `gorm:"type:text;not null"`
	Status    string `gorm:"type:varchar(20);default:'menunggu'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User   User `gorm:"foreignKey:UserID"`
	Dokter User `gorm:"foreignKey:DokterID"`
}
