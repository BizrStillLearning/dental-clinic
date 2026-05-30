package main

import (
	"e-clinic/config"
	"e-clinic/controllers"
	"e-clinic/middlewares"
	"e-clinic/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedUser(db *gorm.DB, nama, email, password, role string) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		db.Create(&models.User{
			Nama:     nama,
			Email:    email,
			Password: string(hashedPassword),
			Role:     role,
		})
	}
}

func main() {
	db := config.ConnectDB()

	db.AutoMigrate(&models.User{}, &models.Patient{}, &models.Appointment{}, &models.MedicalRecord{})

	seedUser(db, "Administrator Utama", "superadmin@gmail.com", "admin123", "superadmin")
	seedUser(db, "dr. Budi Santoso", "dokter@klinik.com", "dokter123", "dokter")
	seedUser(db, "Siti Resepsionis", "resepsionis@klinik.com", "resepsionis123", "resepsionis")
	seedUser(db, "Pasien Dummy", "pasien@klinik.com", "pasien123", "pasien")

	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	authController := controllers.NewAuthController(db)
	patientController := controllers.NewPatientController(db)
	appointmentController := controllers.NewAppointmentController(db)
	doctorController := controllers.NewDoctorController(db)
	staffController := controllers.NewStaffController(db)
	historyController := controllers.NewHistoryController(db)

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login E-Clinic"})
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{"title": "Daftar Akun"})
	})

	r.POST("/register", authController.ProcessRegister)
	r.POST("/login", authController.ProcessLogin)
	r.GET("/logout", authController.Logout)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/login")
	})

	protected := r.Group("/dashboard")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			role, _ := c.Get("userRole")

			var user models.User
			db.First(&user, userID)

			c.HTML(http.StatusOK, "dashboard.html", gin.H{
				"title": "Dashboard E-Clinic",
				"nama":  user.Nama,
				"role":  role,
			})
		})

		protected.GET("/superadmin/staff", middlewares.RoleBlockMiddleware("superadmin"), staffController.ShowStaffManagement)
		protected.POST("/superadmin/staff", middlewares.RoleBlockMiddleware("superadmin"), staffController.ProcessAddStaff)

		// Rute Dokter
		protected.GET("/dokter", middlewares.RoleBlockMiddleware("dokter"), doctorController.ShowDoctorDashboard)

		// Rute Pemeriksaan Medis
		protected.GET("/dokter/periksa/:id", middlewares.RoleBlockMiddleware("dokter"), doctorController.ShowMedicalRecordForm)
		protected.POST("/dokter/periksa/:id", middlewares.RoleBlockMiddleware("dokter"), doctorController.ProcessMedicalRecord)

		// Rute Resepsionis & Pasien
		protected.GET("/register-pasien", middlewares.RoleBlockMiddleware("superadmin", "resepsionis"), patientController.ShowRegisterPatient)
		protected.POST("/register-pasien", middlewares.RoleBlockMiddleware("superadmin", "resepsionis"), patientController.ProcessRegisterPatient)

		protected.GET("/buat-janji", middlewares.RoleBlockMiddleware("pasien"), appointmentController.ShowAppointmentForm)
		protected.POST("/buat-janji", middlewares.RoleBlockMiddleware("pasien"), appointmentController.ProcessAppointment)

		protected.GET("/riwayat", middlewares.RoleBlockMiddleware("pasien"), historyController.ShowPatientHistory)
	}

	r.Run(":8080")
}
