package controllers

import (
	"e-clinic/models"
	"net/http"
	"strings"
	"time"

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
	role, _ := c.Get("userRole")

	var dokterID uint
	switch v := rawUserID.(type) {
	case uint:
		dokterID = v
	case float64:
		dokterID = uint(v)
	}

	var user models.User
	dc.DB.First(&user, dokterID)

	var appointments []models.Appointment
	dc.DB.Preload("User").Where("dokter_id = ? AND status = ?", dokterID, "menunggu").Find(&appointments)

	c.HTML(http.StatusOK, "dokter.html", gin.H{
		"title":        "Ruang Dokter - Antrean Pasien",
		"nama":         user.Nama,
		"role":         role,
		"appointments": appointments,
	})
}

func (dc *DoctorController) ShowMedicalRecordForm(c *gin.Context) {
	appointmentID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")

	var user models.User
	dc.DB.First(&user, userID)

	var appointment models.Appointment
	if err := dc.DB.Preload("User").First(&appointment, appointmentID).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/dashboard/dokter")
		return
	}

	var obats []models.Obat
	dc.DB.Where("stok > ?", 0).Find(&obats)

	c.HTML(http.StatusOK, "medical_record_form.html", gin.H{
		"title":       "Pemeriksaan Pasien",
		"nama":        user.Nama,
		"role":        role,
		"appointment": appointment,
		"obats":       obats,
	})
}

func (dc *DoctorController) ProcessMedicalRecord(c *gin.Context) {
	appointmentID := c.Param("id")
	diagnosis := c.PostForm("diagnosis")
	tindakan := c.PostForm("tindakan")

	obatIDs := c.PostFormArray("obat_ids")

	var appointment models.Appointment
	dc.DB.First(&appointment, appointmentID)

	totalHargaObat := 0
	var namaObatList []string

	for _, idStr := range obatIDs {
		var obat models.Obat
		if err := dc.DB.First(&obat, idStr).Error; err == nil {
			if obat.Stok > 0 {
				obat.Stok -= 1
				dc.DB.Save(&obat)
				totalHargaObat += obat.Harga
				namaObatList = append(namaObatList, obat.Nama)
			}
		}
	}

	resepObat := strings.Join(namaObatList, ", ")
	if resepObat == "" {
		resepObat = "Tidak ada resep obat (Hanya Konsultasi)"
	}

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

	transaksi := models.Transaction{
		AppointmentID: appointment.ID,
		TotalBiaya:    totalHargaObat,
		Status:        "menunggu",
	}
	dc.DB.Create(&transaksi)

	dc.DB.Model(&appointment).Update("status", "menunggu_pembayaran")

	c.Redirect(http.StatusSeeOther, "/dashboard/dokter")
}
