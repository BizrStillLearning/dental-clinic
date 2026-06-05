package controllers

import (
	"e-clinic/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ObatController struct {
	DB *gorm.DB
}

func NewObatController(db *gorm.DB) *ObatController {
	return &ObatController{DB: db}
}

func (oc *ObatController) ShowObatManagement(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")

	var user models.User
	oc.DB.First(&user, userID)

	var obats []models.Obat
	oc.DB.Order("id desc").Find(&obats)

	c.HTML(http.StatusOK, "obat_management.html", gin.H{
		"title": "Inventaris Obat & Apotek",
		"nama":  user.Nama,
		"role":  role,
		"obats": obats,
	})
}

func (oc *ObatController) ProcessAddObat(c *gin.Context) {
	nama := c.PostForm("nama")
	harga, _ := strconv.Atoi(c.PostForm("harga"))
	stok, _ := strconv.Atoi(c.PostForm("stok"))

	obat := models.Obat{
		Nama:  nama,
		Harga: harga,
		Stok:  stok,
	}

	oc.DB.Create(&obat)
	c.Redirect(http.StatusSeeOther, "/dashboard/superadmin/obat")
}

func (oc *ObatController) DeleteObat(c *gin.Context) {
	id := c.Param("id")
	oc.DB.Delete(&models.Obat{}, id)
	c.Redirect(http.StatusSeeOther, "/dashboard/superadmin/obat")
}
