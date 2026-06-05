package controllers

import (
	"e-clinic/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{DB: db}
}

func (tc *TransactionController) ShowKasirDashboard(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")

	var user models.User
	tc.DB.First(&user, userID)

	var transactions []models.Transaction
	tc.DB.Preload("Appointment.User").Preload("Appointment.Dokter").Where("status = ?", "menunggu").Find(&transactions)

	c.HTML(http.StatusOK, "kasir.html", gin.H{
		"title":        "Meja Kasir - Menunggu Pembayaran",
		"nama":         user.Nama,
		"role":         role,
		"transactions": transactions,
	})
}

func (tc *TransactionController) ShowPaymentForm(c *gin.Context) {
	transactionID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("userRole")

	var user models.User
	tc.DB.First(&user, userID)

	var transaksi models.Transaction
	if err := tc.DB.Preload("Appointment.User").Preload("Appointment.Dokter").First(&transaksi, transactionID).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/dashboard/kasir")
		return
	}

	c.HTML(http.StatusOK, "payment_form.html", gin.H{
		"title":     "Proses Pembayaran",
		"nama":      user.Nama,
		"role":      role,
		"transaksi": transaksi,
	})
}

func (tc *TransactionController) ProcessPayment(c *gin.Context) {
	transactionID := c.Param("id")

	totalBiayaStr := c.PostForm("total_biaya")
	totalBiaya, _ := strconv.Atoi(totalBiayaStr)
	metode := c.PostForm("metode_pembayaran")

	var transaksi models.Transaction
	tc.DB.Preload("Appointment").First(&transaksi, transactionID)

	now := time.Now()

	if metode == "Menunggak (Piutang)" {
		transaksi.Status = "menunggak"
	} else {
		transaksi.Status = "lunas"
		transaksi.TanggalBayar = &now
	}

	transaksi.TotalBiaya = totalBiaya
	transaksi.MetodePembayaran = metode
	tc.DB.Save(&transaksi)

	tc.DB.Model(&transaksi.Appointment).Update("status", "lunas")

	c.Redirect(http.StatusSeeOther, "/dashboard/kasir")
}

func (tc *TransactionController) DownloadInvoicePDF(c *gin.Context) {
	transactionID := c.Param("id")

	var trx models.Transaction
	if err := tc.DB.Preload("Appointment.User").Preload("Appointment.Dokter").First(&trx, transactionID).Error; err != nil {
		c.String(http.StatusNotFound, "Data transaksi tidak ditemukan")
		return
	}

	pdf := gofpdf.New("P", "mm", "A5", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(130, 10, "KLINIK UTAMA E-CLINIC", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(130, 5, "Jl. Raya Teknologi No. 037, Surabaya", "", 1, "C", false, 0, "")
	pdf.CellFormat(130, 5, "Telp: (031) 1234567 | Status: "+strings.ToUpper(trx.Status), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.Line(10, pdf.GetY(), 138, pdf.GetY())
	pdf.Ln(4)

	tglBayar := "-"
	if trx.TanggalBayar != nil {
		tglBayar = trx.TanggalBayar.Format("2006-01-02 15:04")
	}
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(30, 6, "No. Transaksi")
	pdf.Cell(5, 6, ":")
	pdf.Cell(90, 6, fmt.Sprintf("TRX-%d", trx.ID))
	pdf.Ln(6)
	pdf.Cell(30, 6, "Nama Pasien")
	pdf.Cell(5, 6, ":")
	pdf.Cell(90, 6, trx.Appointment.User.Nama)
	pdf.Ln(6)
	pdf.Cell(30, 6, "Dokter Jaga")
	pdf.Cell(5, 6, ":")
	pdf.Cell(90, 6, "dr. "+trx.Appointment.Dokter.Nama)
	pdf.Ln(6)
	pdf.Cell(30, 6, "Waktu Bayar")
	pdf.Cell(5, 6, ":")
	pdf.Cell(90, 6, tglBayar)
	pdf.Ln(8)

	pdf.Line(10, pdf.GetY(), 138, pdf.GetY())
	pdf.Ln(2)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(85, 6, "Deskripsi Layanan")
	pdf.Cell(43, 6, "Total")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(85, 6, "Biaya Konsultasi & Jasa Medis Dokter")
	pdf.Cell(43, 6, fmt.Sprintf("Rp %d", trx.TotalBiaya))
	pdf.Ln(6)

	pdf.Line(10, pdf.GetY(), 138, pdf.GetY())
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(85, 6, "TOTAL BAYAR :")
	pdf.Cell(43, 6, fmt.Sprintf("Rp %d", trx.TotalBiaya))
	pdf.Ln(6)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(85, 6, "Metode Pembayaran: "+trx.MetodePembayaran)
	pdf.Ln(15)

	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(128, 5, "Terima kasih atas kepercayaan Anda.", "", 1, "C", false, 0, "")
	pdf.CellFormat(128, 5, "Semoga lekas sembuh.", "", 1, "C", false, 0, "")

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Invoice-TRX-%d.pdf", trx.ID))
	c.Header("Content-Type", "application/pdf")

	_ = pdf.Output(c.Writer)
}

func (tc *TransactionController) ExportFinancialExcel(c *gin.Context) {
	var transactions []models.Transaction
	tc.DB.Preload("Appointment.User").Preload("Appointment.Dokter").Where("status = ?", "lunas").Order("id desc").Find(&transactions)

	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheetName := "Laporan Keuangan"
	index, _ := f.NewSheet(sheetName)
	_ = f.DeleteSheet("Sheet1")

	_ = f.SetCellValue(sheetName, "A1", "LAPORAN PENDAPATAN REAL-TIME E-CLINIC")
	_ = f.SetCellValue(sheetName, "A2", "Diekspor pada: "+time.Now().Format("2006-01-02 15:04:05"))

	headers := map[string]string{
		"A4": "ID Transaksi",
		"B4": "Tanggal Bayar",
		"C4": "Nama Pasien",
		"D4": "Dokter Pemeriksa",
		"E4": "Metode Pembayaran",
		"F4": "Total Pendapatan (Rp)",
	}
	for cell, val := range headers {
		_ = f.SetCellValue(sheetName, cell, val)
	}

	rowNum := 5
	var grandTotal int
	for _, trx := range transactions {
		tgl := "-"
		if trx.TanggalBayar != nil {
			tgl = trx.TanggalBayar.Format("2006-01-02 15:04")
		}

		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("TRX-%d", trx.ID))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), tgl)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), trx.Appointment.User.Nama)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), "dr. "+trx.Appointment.Dokter.Nama)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), trx.MetodePembayaran)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), trx.TotalBiaya)

		grandTotal += trx.TotalBiaya
		rowNum++
	}

	_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), "GRAND TOTAL:")
	_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), grandTotal)

	f.SetActiveSheet(index)

	c.Header("Content-Disposition", "attachment; filename=Laporan-Keuangan-Klinik.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	_ = f.Write(c.Writer)
}
