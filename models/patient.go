package models

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	ID            uint   `gorm:"primaryKey"`
	NomorRM       string `gorm:"type:varchar(20);unique;not null"`
	NamaLengkap   string `gorm:"type:varchar(100);not null"`
	NIK           string `gorm:"type:varchar(16);unique"`
	TanggalLahir  string `gorm:"type:date"`
	Alamat        string `gorm:"type:text"`
	NoHP          string `gorm:"type:varchar(20)"`
	RiwayatAlergi string `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
