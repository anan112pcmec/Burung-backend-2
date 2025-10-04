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

	// Ambil semua id barang induk
	if errDB := db.Model(&models.BarangInduk{}).Pluck("id", &ids).Error; errDB != nil {
		log.Printf("❌ gagal mendapatkan id barang untuk maintain komentar: %v", errDB)
		return
	}

	for _, idBarang := range ids {
		var komentarList []models.Komentar
		if err := db.Model(&models.Komentar{}).
			Where(&models.Komentar{IdBarangInduk: idBarang}).
			Find(&komentarList).Error; err != nil {
			log.Printf("❌ gagal mengambil komentar untuk barang %d: %v", idBarang, err)
			continue
		}

		keyBarang := fmt.Sprintf("komentar_barang:%d", idBarang)
		if err := rds.Del(ctx, keyBarang).Err(); err != nil {
			log.Printf("⚠️ gagal hapus komentar_barang lama untuk barang %d: %v", idBarang, err)
		}

		// Jika tidak ada komentar, kasih "__init__"
		if len(komentarList) == 0 {
			if err := rds.SAdd(ctx, keyBarang, "__init__").Err(); err != nil {
				log.Printf("❌ gagal SADD __init__ untuk barang %d: %v", idBarang, err)
			}
			continue
		}

		// Pakai pipeline biar lebih cepat
		pipe := rds.Pipeline()
		for _, k := range komentarList {
			keyKomentar := fmt.Sprintf("komentar:%d", k.ID)

			// Hapus komentar lama
			pipe.Del(ctx, keyKomentar)

			// Set ulang detail komentar
			pipe.HSet(ctx, keyKomentar, map[string]interface{}{
				"id":              k.ID,
				"id_barang_induk": k.IdBarangInduk,
				"id_entity":       k.IdEntity,
				"komentar":        k.Komentar,
				"jenis_entity":    k.JenisEntity,
				"parent_id":       k.ParentID,
			})

			// Tambahkan ke set barang
			pipe.SAdd(ctx, keyBarang, keyKomentar)
		}

		if _, err := pipe.Exec(ctx); err != nil {
			log.Printf("❌ gagal execute pipeline Redis untuk barang %d: %v", idBarang, err)
		} else {
			log.Printf("✅ %d komentar barang %d berhasil dimuat ke Redis", len(komentarList), idBarang)
		}
	}
}
