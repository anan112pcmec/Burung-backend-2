package maintain

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
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

		// kalau tidak ada komentar, skip
		if len(komentarList) == 0 {
			continue
		}

		for _, k := range komentarList {
			// Simpan detail komentar (pakai map agar bisa HSET sekali jalan)
			fields := models.Komentar{
				ID:            k.ID,
				IdBarangInduk: k.IdBarangInduk,
				IdEntity:      k.IdEntity,
				IsiKomentar:   k.IsiKomentar,
				JenisEntity:   k.JenisEntity,
				ParentID:      k.ParentID,
			}

			if ada_err := rds.Del(ctx, fmt.Sprintf("komentar:%d", k.ID)).Err(); ada_err != nil {
				continue
			}

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
				if err := rds.HSet(ctx, fmt.Sprintf("komentar:%d", k.ID), fieldName, value).Err(); err != nil {
					fmt.Println("❌ Gagal Set Redis:", err)
				}
			}

			if err := rds.SAdd(ctx, fmt.Sprintf("komentar_barang:%d", idBarang), fmt.Sprintf("komentar:%d", k.ID)).Err(); err != nil {
				log.Printf("❌ gagal SADD komentar %d ke barang %d: %v", k.ID, idBarang, err)
			}

		}

		log.Printf("✅ komentar barang %d berhasil dimuat ke Redis (%d komentar)", idBarang, len(komentarList))
	}
}
