package controllers

import (
	"e-clinic/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppointmentController struct {
	DB *gorm.DB
}

func NewAppointmentController(db *gorm.DB) *AppointmentController {
	return &AppointmentController{DB: db}
}

func (ac *AppointmentController) ShowAppointmentForm(c *gin.Context) {
	var dokters []models.User
	ac.DB.Where("role = ?", "dokter").Find(&dokters)

	c.HTML(http.StatusOK, "appointment.html", gin.H{
		"title":   "Buat Janji Temu",
		"dokters": dokters,
	})
}

func (ac *AppointmentController) ProcessAppointment(c *gin.Context) {
	userID, _ := c.Get("userID")

	dokterID := c.PostForm("dokter_id")
	tanggal := c.PostForm("tanggal")
	keluhan := c.PostForm("keluhan")

	appointment := models.Appointment{
		UserID:   userID.(uint),
		DokterID: stringToUint(dokterID),
		Tanggal:  tanggal,
		Keluhan:  keluhan,
		Status:   "menunggu",
	}

	if err := ac.DB.Create(&appointment).Error; err != nil {
		fmt.Println("=====================================")
		fmt.Println("DB ERROR:", err)
		fmt.Println("=====================================")

		c.HTML(http.StatusInternalServerError, "appointment.html", gin.H{
			"error": "Gagal membuat janji temu. Silakan coba lagi.",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/dashboard/")
}

func stringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}
