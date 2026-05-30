package controllers

import (
	"e-clinic/models"
	"fmt"
	"net/http"
	"time" // Ditambahkan untuk mengambil tanggal hari ini otomatis

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DoctorController struct {
	DB *gorm.DB
}

func NewDoctorController(db *gorm.DB) *DoctorController {
	return &DoctorController{DB: db}
}

func (dc *DoctorController) ShowDoctorDashboard(c *gin.Context) {
	rawUserID, _ := c.Get("userID")

	var dokterID uint
	switch v := rawUserID.(type) {
	case uint:
		dokterID = v
	case float64:
		dokterID = uint(v)
	}

	fmt.Println("\n=====================================")
	fmt.Println(">>> [DEBUG] Dokter yg Login ID:", dokterID)

	var appointments []models.Appointment

	err := dc.DB.Preload("User").Where("dokter_id = ? AND status = ?", dokterID, "menunggu").Find(&appointments).Error

	if err != nil {
		fmt.Println(">>> [DEBUG] Error GORM:", err)
	}
	fmt.Println(">>> [DEBUG] Jumlah Antrean Ditemukan:", len(appointments))
	fmt.Println("=====================================\n")

	c.HTML(http.StatusOK, "dokter.html", gin.H{
		"title":        "Ruang Dokter - Antrean Pasien",
		"appointments": appointments,
	})
}

func (dc *DoctorController) ShowMedicalRecordForm(c *gin.Context) {
	appointmentID := c.Param("id")

	var appointment models.Appointment
	if err := dc.DB.Preload("User").First(&appointment, appointmentID).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/dashboard/dokter")
		return
	}

	c.HTML(http.StatusOK, "medical_record_form.html", gin.H{
		"title":       "Pemeriksaan Pasien",
		"appointment": appointment,
	})
}

func (dc *DoctorController) ProcessMedicalRecord(c *gin.Context) {
	appointmentID := c.Param("id")

	diagnosis := c.PostForm("diagnosis")
	tindakan := c.PostForm("tindakan")
	resepObat := c.PostForm("resep_obat")

	var appointment models.Appointment
	dc.DB.First(&appointment, appointmentID)

	record := models.MedicalRecord{
		AppointmentID: appointment.ID,
		UserID:        appointment.UserID,
		DokterID:      appointment.DokterID,
		Diagnosis:     diagnosis,
		Tindakan:      tindakan,
		ResepObat:     resepObat,
		Tanggal:       time.Now().Format("2006-01-02"),
	}

	dc.DB.Create(&record)

	dc.DB.Model(&appointment).Update("status", "selesai")

	c.Redirect(http.StatusSeeOther, "/dashboard/dokter")
}
