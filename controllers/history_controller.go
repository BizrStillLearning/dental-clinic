package controllers

import (
	"e-clinic/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HistoryController struct {
	DB *gorm.DB
}

type HistoryView struct {
	Record models.MedicalRecord
	Trx    models.Transaction
}

func NewHistoryController(db *gorm.DB) *HistoryController {
	return &HistoryController{DB: db}
}

func (hc *HistoryController) ShowPatientHistory(c *gin.Context) {
	rawUserID, _ := c.Get("userID")

	var userID uint
	switch v := rawUserID.(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	}

	var records []models.MedicalRecord
	hc.DB.Preload("Dokter").Preload("Appointment").Where("user_id = ?", userID).Find(&records)

	var views []HistoryView
	for _, rec := range records {
		var trx models.Transaction
		hc.DB.Where("appointment_id = ?", rec.AppointmentID).First(&trx)

		views = append(views, HistoryView{
			Record: rec,
			Trx:    trx,
		})
	}

	c.HTML(http.StatusOK, "medical_history.html", gin.H{
		"title": "Riwayat Medis & Tagihan",
		"views": views,
	})
}

func (hc *HistoryController) ShowAllHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")

	var user models.User
	hc.DB.First(&user, userID)

	var records []models.MedicalRecord

	hc.DB.Preload("User").Preload("Dokter").Preload("Appointment").Find(&records)

	c.HTML(http.StatusOK, "all_medical_history.html", gin.H{
		"title":   "Rekam Medis",
		"nama":    user.Nama,
		"role":    role,
		"records": records,
	})
}
