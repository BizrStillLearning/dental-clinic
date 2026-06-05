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

func (pc *PatientController) ShowPatientRegister(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")

	var user models.User
	pc.DB.First(&user, userID)

	c.HTML(http.StatusOK, "patient_register.html", gin.H{
		"title": "Pendaftaran Pasien Baru",
		"nama":  user.Nama,
		"role":  role,
	})
}

func (pc *PatientController) ProcessRegisterPatient(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")

	var user models.User
	pc.DB.First(&user, userID)

	namaLengkap := c.PostForm("nama")
	nik := c.PostForm("nik")
	tanggalLahir := c.PostForm("tanggal_lahir")
	alamat := c.PostForm("alamat")
	noHp := c.PostForm("nomor_hp")
	riwayatAlergi := c.PostForm("alergi")

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
			"title": "Pendaftaran Pasien Baru",
			"nama":  user.Nama,
			"role":  role,
			"error": "Gagal mendaftarkan pasien. NIK mungkin sudah digunakan.",
		})
		return
	}

	c.HTML(http.StatusOK, "patient_register.html", gin.H{
		"title":   "Pendaftaran Pasien Baru",
		"nama":    user.Nama,
		"role":    role,
		"success": "Pasien berhasil didaftarkan dengan Nomor RM: " + nomorRM,
	})
}
