package services

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

func OnlinePengguna(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadPengguna, rds *redis.Client) {

	data.Validate()
	fmt.Println("Caching Online User")

	key := fmt.Sprintf("pengguna_online:%v", data.Id)

	fields := map[string]interface{}{
		"nama":     data.Nama,
		"username": data.Username,
		"email":    data.Email,
	}

	for field, value := range fields {
		if err := rds.HSet(ctx, key, field, value).Err(); err != nil {
			fmt.Println("Gagal Set Redis:", err)
		}
	}
}

func OfflinePengguna(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadPengguna, rds *redis.Client) {
	data.Validate()
	key := fmt.Sprintf("pengguna_online:%v", data.Id)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Println("Gagal Hapus Redis Key:", err)
	} else {
		fmt.Println("✅ User offline, key dihapus:", key)
	}
}
