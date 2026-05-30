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

	c.HTML(http.StatusOK, "medical_history.html", gin.H{
		"title":   "Riwayat Rekam Medis",
		"records": records,
	})
}
