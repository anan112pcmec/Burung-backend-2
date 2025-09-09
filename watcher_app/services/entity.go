package services

import (
	"context"
	"fmt"
	"reflect"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

func OnlinePengguna(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadPengguna, rds *redis.Client) {

	fmt.Println("Caching Online User")

	key := fmt.Sprintf("pengguna_online:%v", data.ID)

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
	key := fmt.Sprintf("pengguna_online:%v", data.ID)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Println("Gagal Hapus Redis Key:", err)
	} else {
		fmt.Println("✅ User offline, key dihapus:", key)
	}
}

func UpSeller(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsePayloadSeller, rds *redis.Client) {
	fmt.Println("Caching NEW Seller")

	key := fmt.Sprintf("seller_data:%v", data.ID)

	fields := models.Seller{
		ID:               data.ID,
		Username:         data.Username,
		Nama:             data.Nama,
		SellerDedication: data.SellerDedication,
		FollowerTotal:    data.FollowerTotal,
	}

	v := reflect.ValueOf(fields)
	t := reflect.TypeOf(fields)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		if fieldName == "" {
			fieldName = t.Field(i).Name
		}
		value := v.Field(i).Interface()

		if err := rds.HSet(ctx, key, fieldName, value).Err(); err != nil {
			fmt.Println("Gagal Set Redis:", err)
		}
	}
}

func HapusSeller(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsePayloadSeller, rds *redis.Client) {
	key := fmt.Sprintf("seller_data:%v", data.ID)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Println("Gagal Hapus Redis Key:", err)
	} else {
		fmt.Println("✅ seller dihapus, key dihapus:", key)
	}
}
