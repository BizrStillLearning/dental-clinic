package controllers

import (
	"e-clinic/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffController struct {
	DB *gorm.DB
}

func NewStaffController(db *gorm.DB) *StaffController {
	return &StaffController{DB: db}
}

func (sc *StaffController) ShowStaffManagement(c *gin.Context) {
	var staffs []models.User
	sc.DB.Where("role IN ?", []string{"dokter", "resepsionis"}).Find(&staffs)

	c.HTML(http.StatusOK, "staff_management.html", gin.H{
		"title":  "Manajemen Staf Klinik",
		"staffs": staffs,
	})
}

func (sc *StaffController) ProcessAddStaff(c *gin.Context) {
	nama := c.PostForm("nama")
	email := c.PostForm("email")
	password := c.PostForm("password")
	role := c.PostForm("role")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Nama:     nama,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	sc.DB.Create(&user)

	c.Redirect(http.StatusSeeOther, "/dashboard/superadmin/staff")
}
