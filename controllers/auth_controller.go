package controllers

import (
	"e-clinic/models"
	"e-clinic/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) ProcessRegister(c *gin.Context) {
	nama := c.PostForm("nama")
	email := c.PostForm("email")
	password := c.PostForm("password")

	role := "pasien"

	if nama == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Semua kolom wajib diisi!"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Nama:     nama,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Email sudah terdaftar atau terjadi kesalahan."})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func (ac *AuthController) ProcessLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	if err := ac.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Email atau password salah!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Email atau password salah!"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Gagal memproses sesi login."})
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/dashboard/")
}

func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/login")
}
