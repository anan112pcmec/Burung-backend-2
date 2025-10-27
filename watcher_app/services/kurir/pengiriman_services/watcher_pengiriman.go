package kurir_pengiriman_watcher

import (
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
)

func DiperjalananConfirmation(IdTransaksi int64, status string, db *gorm.DB) {
	if status != "Picked Up" {
		return
	}

	if IdTransaksi == 0 {
		return
	}

	_ = db.Model(&models.Transaksi{}).Where(&models.Transaksi{
		ID: IdTransaksi,
	}).Update("status", "Dikirim")
}

func SampaiConfirmation(IdTransaksi int64, status string, db *gorm.DB) {
	if status != "Sampai" {
		return
	}
	if IdTransaksi == 0 {
		return
	}
	_ = db.Model(&models.Transaksi{}).Where(&models.Transaksi{
		ID: IdTransaksi,
	}).Update("status", "Selesai")
}
