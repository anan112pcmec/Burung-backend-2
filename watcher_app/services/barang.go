package services

import (
	"context"
	"fmt"
	"log"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

// //////////////////////////////////////////////////////////////////////////////////////////////
// BARANG INDUK
// //////////////////////////////////////////////////////////////////////////////////////////////

// 1. Lebih Bertujuan Untuk Melakukan Caching

func BarangMasuk(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadBarang, rds *redis.Client, SE meilisearch.ServiceManager) {
	fmt.Println("🔔 Mulai proses caching Barang")

	if data.OriginalKategori == "" {
		fmt.Println("⚠️ OriginalKategori kosong, skip caching")
		return
	}

	go func() {
		var harga int32
		err := db.Model(models.KategoriBarang{}).
			Where(models.KategoriBarang{Nama: data.OriginalKategori, IdBarangInduk: data.ID}).
			Select("harga").
			Take(&harga).Error

		if err != nil {
			fmt.Printf("❌ Gagal ambil harga dari DB untuk kategori %s (barang ID %d): %v\n",
				data.OriginalKategori, data.ID, err)
			harga = 0
		} else {
			fmt.Printf("✅ Berhasil ambil harga: %d untuk kategori %s\n", harga, data.OriginalKategori)
		}

		key := fmt.Sprintf("barang:%v", data.ID)

		fields := map[string]interface{}{
			"id_barang_induk":             data.ID,
			"id_seller_barang_induk":      data.SellerID,
			"nama_barang_induk":           data.NamaBarang,
			"jenis_barang_induk":          data.JenisBarang,
			"original_kategori":           data.OriginalKategori,
			"deskripsi_barang_induk":      data.Deskripsi,
			"tanggal_rilis_barang_induk":  data.TanggalRilis,
			"viewed_barang_induk":         data.Viewed,
			"likes_barang_induk":          data.Likes,
			"total_komentar_barang_induk": data.TotalKomentar,
			"created_at":                  data.CreatedAt,
			"updated_at":                  data.UpdatedAt,
			"deleted_at":                  data.DeletedAt,
			"harga":                       harga,
		}

		for field, value := range fields {
			if err := rds.HSet(ctx, key, field, value).Err(); err != nil {
				fmt.Printf("❌ Gagal set Redis field=%s value=%v error=%v\n", field, value, err)
			} else {
				fmt.Printf("✅ Redis set OK field=%s value=%v\n", field, value)
			}
		}
	}()

	go func() {
		barangIndukIndex := SE.Index("barang_induk_all")
		var documents []map[string]interface{}

		fmt.Println("baraang", data)
		doc := map[string]interface{}{
			"id":                         data.BarangInduk.ID,
			"id_barang_induk":            data.BarangInduk.ID,
			"nama_barang_induk":          data.NamaBarang,
			"id_seller_barang_induk":     data.SellerID,
			"jenis_barang_induk":         data.JenisBarang,
			"tanggal_rilis_barang_induk": data.TanggalRilis,
			"viewed_barang_induk":        data.Viewed,
			"likes_barang_induk":         data.Likes,
		}
		documents = append(documents, doc)

		task, err := barangIndukIndex.AddDocuments(documents, nil)
		if err != nil {
			log.Fatal("❌ Gagal menambahkan dokumen ke Meilisearch:", err)
		}

		log.Println("✅ Task UID:", task.TaskUID)
	}()

	go func() {
		key := fmt.Sprintf("barang:%v", data.ID)
		rds.SAdd(ctx, fmt.Sprintf("jenis_%s_barang", helper.ConvertJenisBarangReverse(data.JenisBarang)), key)
	}()

	fmt.Println("🎉 Proses caching selesai untuk barang:", data.NamaBarang)
}

func HapusBarang(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadBarang, rds *redis.Client) {
	fmt.Println("Mulai Hapus Barang")

	key := fmt.Sprintf("barang:%v", data.ID)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Println("Gagal Hapus Redis Key:", err)
	} else {
		fmt.Println("✅ User offline, key dihapus:", key)
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////////
// VARIAN BARANG
// //////////////////////////////////////////////////////////////////////////////////////////////

// 1. Lebih Bertujuan Untuk Mengendalikan Barang Apalagi Jika Berkenaan dengan Transaksi
