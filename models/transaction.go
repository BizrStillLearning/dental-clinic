package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID               uint   `gorm:"primaryKey"`
	AppointmentID    uint   `gorm:"not null;unique"`
	TotalBiaya       int    `gorm:"not null"`
	MetodePembayaran string `gorm:"type:varchar(50)"`
	Status           string `gorm:"type:varchar(50);default:'menunggu'"`
	TanggalBayar     *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	Appointment Appointment `gorm:"foreignKey:AppointmentID"`
}
