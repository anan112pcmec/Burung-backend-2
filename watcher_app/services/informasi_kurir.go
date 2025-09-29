package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
)

func VerifiedKurir(ctx context.Context, id_kurir int64, status_perizinan string, db *gorm.DB) {
	if status_perizinan == "Diizinkan" {
		var diizinkan_info_kurir string = ""
		var diizinkan_info_kendaraan string = ""

		_ = db.Model(models.InformasiKurir{}).Select("informasi_status_perizinan").Where(models.InformasiKurir{
			IDkurir: id_kurir,
		}).Take(&diizinkan_info_kurir)

		_ = db.Model(models.InformasiKendaraanKurir{}).Select("informasi_status_perizinan").Where(models.InformasiKendaraanKurir{
			ID: id_kurir,
		}).Take(&diizinkan_info_kendaraan)

		if diizinkan_info_kendaraan == "Diizinkan" && diizinkan_info_kurir == "Diizinkan" {
			_ = db.Model(models.Kurir{}).Where(models.Kurir{
				ID: id_kurir,
			}).Update("verified", true)
		}
	}
}
