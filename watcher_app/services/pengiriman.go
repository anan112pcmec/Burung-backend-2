package services

import (
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
)

func WaitingConfirmation(IdTransaksi, IdKurir int64, status string, db *gorm.DB) {
	if status != "Packaging" {
		return
	}

	if IdTransaksi == 0 {
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err_update_transaksi := tx.Model(&models.Transaksi{}).Where(&models.Transaksi{
			ID: IdTransaksi,
		}).Update("status", "Waiting").Error; err_update_transaksi != nil {
			return err_update_transaksi
		}

		return nil
	}); err != nil {
		return
	}
}

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

// // if err_update_kurir := tx.Model(&models.Kurir{}).Where(&models.Kurir{
// 			ID: IdKurir,
// 		}).UpdateColumn("viewed", gorm.Expr("viewed + 1")).Error; err_update_kurir != nil {
// 			return err_update_kurir
// 		}
