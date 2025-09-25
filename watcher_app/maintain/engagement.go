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
		log.Printf("❌ gagal mendapatkan id barang untuk maintain komentar: %v", errDB)
		return
	}

	for _, id := range ids {
		key := fmt.Sprintf("komentar_barang:%d", id)
		if err_hapus := rds.Del(ctx, key).Err(); err_hapus != nil {
			log.Printf("Gagal Hapus Dulu")
		}
		if err := rds.SAdd(ctx, key, "__init__").Err(); err != nil {
			log.Printf("gagal membuat SADD untuk barang id=%d: %v", id, err)
			continue
		}
		// kalau ingin log success, bisa dibatasi
		log.Printf("SADD set komentar untuk barang id=%d berhasil", id)
	}

	for _, idBarang := range ids {
		var komentarList []models.Komentar
		if err := db.Where("id_barang_induk = ?", idBarang).Find(&komentarList).Error; err != nil {
			log.Printf("❌ gagal mengambil komentar untuk barang %d: %v", idBarang, err)
			continue
		}

		if len(komentarList) == 0 {
			continue
		}

		// 🧹 Bersihkan set lama untuk komentar barang ini
		if err := rds.Del(ctx, fmt.Sprintf("komentar_barang:%d", idBarang)).Err(); err != nil {
			log.Printf("⚠️ gagal hapus komentar_barang lama untuk barang %d: %v", idBarang, err)
		}

		for _, k := range komentarList {
			keyKomentar := fmt.Sprintf("komentar:%d", k.ID)

			// 🧹 Bersihkan komentar lama per ID
			if err := rds.Del(ctx, keyKomentar).Err(); err != nil {
				log.Printf("⚠️ gagal hapus komentar lama %d: %v", k.ID, err)
			}

			// 🚀 Simpan komentar (langsung tanpa reflect biar lebih efisien)
			if err := rds.HSet(ctx, keyKomentar, map[string]interface{}{
				"id":              k.ID,
				"id_barang_induk": k.IdBarangInduk,
				"id_entity":       k.IdEntity,
				"komentar":        k.Komentar,
				"jenis_entity":    k.JenisEntity,
				"parent_id":       k.ParentID,
			}).Err(); err != nil {
				log.Printf("❌ gagal HSET Redis komentar %d: %v", k.ID, err)
			}

			// Tambahkan ke daftar komentar barang
			if err := rds.SAdd(ctx, fmt.Sprintf("komentar_barang:%d", idBarang), keyKomentar).Err(); err != nil {
				log.Printf("❌ gagal SADD komentar %d ke barang %d: %v", k.ID, idBarang, err)
			}
		}

		log.Printf("✅ %d komentar barang %d berhasil dimuat ke Redis", len(komentarList), idBarang)
	}

}
