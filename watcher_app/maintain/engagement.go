package maintain

import (
	"context"
	"fmt"
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
			PendingTransaksiMaintain(ctx, db, rds)
			time.Sleep(10 * time.Minute)
		}
	}
}

func PendingTransaksiMaintain(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	if _, err := rds.FlushDB(ctx).Result(); err != nil {
		fmt.Println("❌ Gagal melakukan flush Redis:", err)
	} else {
		fmt.Println("✅ Redis berhasil dibersihkan (FlushDB).")
	}

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
