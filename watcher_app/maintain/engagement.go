package maintain

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"
)

func EngagementMaintainLoop(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("❌ BarangMaintainLoop dihentikan")
			return
		default:
			EngagementMaintain(ctx, db, rds)
			PendingTransaksiMaintain(ctx, db, rds)
			time.Sleep(10 * time.Minute)
		}
	}
}

func EngagementMaintain(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	var ids []int32

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

		if len(komentarList) == 0 {
			if err := rds.SAdd(ctx, keyBarang, "__init__").Err(); err != nil {
				log.Printf("❌ gagal SADD __init__ untuk barang %d: %v", idBarang, err)
			}
			continue
		}

		pipe := rds.Pipeline()
		for _, k := range komentarList {
			keyKomentar := fmt.Sprintf("komentar:%d", k.ID)

			pipe.Del(ctx, keyKomentar)

			pipe.HSet(ctx, keyKomentar, map[string]interface{}{
				"id":              k.ID,
				"id_barang_induk": k.IdBarangInduk,
				"id_entity":       k.IdEntity,
				"komentar":        k.Komentar,
				"jenis_entity":    k.JenisEntity,
				"parent_id":       k.ParentID,
			})

			pipe.SAdd(ctx, keyBarang, keyKomentar)
		}

		if _, err := pipe.Exec(ctx); err != nil {
			log.Printf("❌ gagal execute pipeline Redis untuk barang %d: %v", idBarang, err)
		} else {
			log.Printf("✅ %d komentar barang %d berhasil dimuat ke Redis", len(komentarList), idBarang)
		}
	}
}

func PendingTransaksiMaintain(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var hapusTransaksiPending []string

	allKeyspengguna, _ := rds.Keys(ctx, "transaction_pengguna_pending_code:*").Result()

	if len(allKeyspengguna) > 0 {
		for _, key := range allKeyspengguna {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()

				transactionTime, err := rds.HGet(ctx, k, "transaction_time").Result()
				if err != nil {
					return
				}

				if helper.ShouldDelete(transactionTime) {
					mu.Lock()
					hapusTransaksiPending = append(hapusTransaksiPending, k)
					mu.Unlock()
				}
			}(key)
		}

		wg.Wait()
	}

	allKeysSeller, _ := rds.Keys(ctx, "transaction_seller_pending_code:*").Result()

	if len(allKeysSeller) > 0 {
		for _, key := range allKeysSeller {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()

				transactionTime, err := rds.HGet(ctx, k, "transaction_time").Result()
				if err != nil {
					return
				}

				if helper.ShouldDelete(transactionTime) {
					mu.Lock()
					hapusTransaksiPending = append(hapusTransaksiPending, k)
					mu.Unlock()
				}
			}(key)
		}

		wg.Wait()

	}

	if len(hapusTransaksiPending) == 0 {
		return
	}

	pipe := rds.Pipeline()
	for _, key := range hapusTransaksiPending {
		pipe.Del(ctx, key)
	}
	_, _ = pipe.Exec(ctx)

	var transaksi []models.Transaksi
	_ = db.Model(&models.Transaksi{}).Find(&transaksi)
	for _, t := range transaksi {
		const targetLayout = "2006-01-02 15:04:05"
		createdAtStr := t.CreatedAt.Format(targetLayout)

		if !helper.ShouldDelete(createdAtStr) {
			_ = rds.Del(ctx, fmt.Sprintf("transaction_pengguna_pending_id:%v:transaction_code:%s",
				t.IdPengguna, t.KodeOrder))
		}
	}
}
