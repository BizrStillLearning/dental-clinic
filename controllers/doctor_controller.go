package controllers

import (
	"e-clinic/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DoctorController struct {
	DB *gorm.DB
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
