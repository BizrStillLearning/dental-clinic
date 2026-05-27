package controllers

import (
	"e-clinic/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PatientController struct {
	DB *gorm.DB
}

func NewPatientController(db *gorm.DB) *PatientController {
	return &PatientController{DB: db}
}

func (pc *PatientController) ShowRegisterPatient(c *gin.Context) {
	role, _ := c.Get("userRole")
	c.HTML(http.StatusOK, "patient_register.html", gin.H{
		"title": "Pendaftaran Pasien Baru",
		"role":  role,
	})
}

// Memproses data form pendaftaran pasien
func (pc *PatientController) ProcessRegisterPatient(c *gin.Context) {
	role, _ := c.Get("userRole")

	namaLengkap := c.PostForm("nama_lengkap")
	nik := c.PostForm("nik")
	tanggalLahir := c.PostForm("tanggal_lahir")
	alamat := c.PostForm("alamat")
	noHp := c.PostForm("no_hp")
	riwayatAlergi := c.PostForm("riwayat_alergi")

	nomorRM := fmt.Sprintf("RM-%d", time.Now().Unix())

	patient := models.Patient{
		NomorRM:       nomorRM,
		NamaLengkap:   namaLengkap,
		NIK:           nik,
		TanggalLahir:  tanggalLahir,
		Alamat:        alamat,
		NoHP:          noHp,
		RiwayatAlergi: riwayatAlergi,
	}

	if err := pc.DB.Create(&patient).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "patient_register.html", gin.H{
			"error": "Gagal mendaftarkan pasien. NIK mungkin sudah digunakan.",
			"role":  role,
		})
		return
	}

	c.HTML(http.StatusOK, "patient_register.html", gin.H{
		"success": "Pasien berhasil didaftarkan dengan Nomor RM: " + nomorRM,
		"role":    role,
	})
}
