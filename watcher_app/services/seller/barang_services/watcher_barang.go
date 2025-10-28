package seller_barang_watcher

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

// //////////////////////////////////////////////////////////////////////////////////////////////
// BARANG INDUK
// //////////////////////////////////////////////////////////////////////////////////////////////

// 1. Lebih Bertujuan Untuk Melakukan Caching

func BarangReady(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadKategoriBarang) {
	if data.IDRekening != 0 && data.IDAlamat != 0 {
		if err := db.Model(&models.VarianBarang{}).Where(&models.VarianBarang{
			IdKategori: data.ID,
		}).Update("status", "Ready").Error; err != nil {
			fmt.Println("Gagal Mengubah Status Menjadi Ready")
		}
	}
}
