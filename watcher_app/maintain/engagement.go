package maintain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
)

func EngagementMaintainLoop(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("❌ BarangMaintainLoop dihentikan")
			return
		default:
			EngagementMaintain(ctx, db, rds)
			time.Sleep(10 * time.Minute)
		}
	}
}

func EngagementMaintain(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	var ids []int32

	// ambil semua id barang
	if errDB := db.Model(&models.BarangInduk{}).Pluck("id", &ids).Error; errDB != nil {
		log.Printf("gagal mendapatkan id barang untuk maintain komentar: %v", errDB)
		return
	}

	for _, id := range ids {
		key := fmt.Sprintf("komentar_barang:%d", id)
		if err := rds.SAdd(ctx, key, "__init__").Err(); err != nil {
			log.Printf("gagal membuat SADD untuk barang id=%d: %v", id, err)
			continue
		}
		// kalau ingin log success, bisa dibatasi
		log.Printf("SADD set komentar untuk barang id=%d berhasil", id)
	}
}
