package services

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

func UpCacheKomentar(ctx context.Context, data notify_payload.NotifyResponsePayloadKomentar, rds *redis.Client) {
	key := fmt.Sprintf("komentar:%v", data.ID)

	fields := models.Komentar{
		ID:            data.ID,
		IdBarangInduk: data.IdBarangInduk,
		IdEntity:      data.IdEntity,
		IsiKomentar:   data.IsiKomentar,
		JenisEntity:   data.JenisEntity,
		ParentID:      data.ParentID,
	}

	fmt.Println("Caching NEW Komentar")

	v := reflect.ValueOf(fields)
	t := reflect.TypeOf(fields)

	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		parts := strings.Split(tag, ",")
		fieldName := parts[0]
		if fieldName == "" {
			fieldName = t.Field(i).Name
		}

		value := fmt.Sprintf("%v", v.Field(i).Interface()) // konversi ke string
		if err := rds.HSet(ctx, key, fieldName, value).Err(); err != nil {
			fmt.Println("❌ Gagal Set Redis:", err)
		}
	}
}

func EditCacheKomentar(ctx context.Context, data notify_payload.NotifyResponsePayloadKomentar, rds *redis.Client) {
	key := fmt.Sprintf("komentar:%v", data.ID)

	if err := rds.HSet(ctx, key, "isi_komentar", data.IsiKomentar).Err(); err != nil {
		fmt.Println("Gagal update Redis:", err)
	}
}
