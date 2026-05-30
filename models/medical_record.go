package models

import (
	"time"

	"gorm.io/gorm"
)

type MedicalRecord struct {
	ID            uint   `gorm:"primaryKey"`
	AppointmentID uint   `gorm:"not null"`
	UserID        uint   `gorm:"not null"`
	DokterID      uint   `gorm:"not null"`
	Diagnosis     string `gorm:"type:text;not null"`
	Tindakan      string `gorm:"type:text"`
	ResepObat     string `gorm:"type:text"`
	Tanggal       string `gorm:"type:date;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	Appointment Appointment `gorm:"foreignKey:AppointmentID"`
	User        User        `gorm:"foreignKey:UserID"`
	Dokter      User        `gorm:"foreignKey:DokterID"`
}
