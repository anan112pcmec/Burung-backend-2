package maintain

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"

)

// convertJenisBarang akan mengubah nama jenis internal jadi format DB

func BarangMaintainLoop(ctx context.Context, db *gorm.DB, rds *redis.Client, SE meilisearch.ServiceManager) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("❌ BarangMaintainLoop dihentikan")
			return
		default:
			CachingBarangMaintain(ctx, db, rds, SE)
			InternalBarangMaintain(ctx, db)
			time.Sleep(10 * time.Minute)
		}
	}
}

type UpdateBarangInduk struct {
	IdBarangInduk     int32
	ViewedBarangInduk int32
	LikesBarangInduk  int32
}

func CachingBarangMaintain(ctx context.Context, db *gorm.DB, rds *redis.Client, SE meilisearch.ServiceManager) {

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Evaluasi Internal db Dari Cache

	keys, _ := rds.Keys(ctx, "barang:*").Result()

	if len(keys) != 0 {
		updateBarangInduk := make([]UpdateBarangInduk, 0, len(keys))

		for _, k := range keys {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()

				id := strings.TrimPrefix(k, "barang:")

				jumlah_viewed, err_viewed := rds.HGet(ctx, k, "viewed_barang_induk").Result()
				if err_viewed != nil {
					return
				}

				jumlah_likes, err_likes := rds.HGet(ctx, k, "likes_barang_induk").Result()
				if err_likes != nil {
					return
				}

				id_barang_induk, err_id_barang_induk := strconv.Atoi(id)
				if err_id_barang_induk != nil {
					return
				}

				jumlah_viewed_barang, err_jumlah_viewed := strconv.Atoi(jumlah_viewed)
				if err_jumlah_viewed != nil {
					return
				}

				jumlah_likes_barang, err_jumlah_likes := strconv.Atoi(jumlah_likes)
				if err_jumlah_likes != nil {
					return
				}

				data := UpdateBarangInduk{
					IdBarangInduk:     int32(id_barang_induk),
					ViewedBarangInduk: int32(jumlah_viewed_barang),
					LikesBarangInduk:  int32(jumlah_likes_barang),
				}

				mu.Lock()
				updateBarangInduk = append(updateBarangInduk, data)
				mu.Unlock()
			}(k)
		}

		wg.Wait()

		_ = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for _, update_data := range updateBarangInduk {
				if err_update := tx.Model(&models.BarangInduk{}).Where(&models.BarangInduk{
					ID: update_data.IdBarangInduk,
				}).Updates(&models.BarangInduk{
					Viewed: update_data.ViewedBarangInduk,
					Likes:  update_data.LikesBarangInduk,
				}).Error; err_update != nil {
					return err_update
				}
			}
			return nil
		})
	}

	// Evaluasi Cache dari internal db

	idbar := []int32{}
	if err := db.Model(&models.BarangInduk{}).Pluck("id", &idbar).Error; err != nil {
		log.Println("❌ Gagal Mendapatkan Id Barang:", err)
		return
	}

	if len(idbar) == 0 {
		log.Println("❌ Tidak ada Id Barang ditemukan")
		return
	}

	dataBarangInduk := []models.BarangInduk{}
	if err := db.Where("id IN ?", idbar).Find(&dataBarangInduk).Error; err != nil {
		log.Println("❌ Gagal mengambil data barang:", err)
		return
	}

	barangIndukIndex := SE.Index("barang_induk_all")
	var documents []map[string]interface{}

	for _, b := range dataBarangInduk {
		fmt.Println("barang", b.NamaBarang)
		doc := map[string]interface{}{
			"id":                         b.ID,
			"id_barang_induk":            b.ID,
			"nama_barang_induk":          b.NamaBarang,
			"id_seller_barang_induk":     b.SellerID,
			"original_kategori":          b.OriginalKategori,
			"deskripsi":                  b.Deskripsi,
			"jenis_barang_induk":         b.JenisBarang,
			"tanggal_rilis_barang_induk": b.TanggalRilis,
			"viewed_barang_induk":        b.Viewed,
			"likes_barang_induk":         b.Likes,
		}
		documents = append(documents, doc)
	}

	task, err := barangIndukIndex.AddDocuments(documents, nil)
	if err != nil {
		log.Fatal("❌ Gagal menambahkan dokumen ke Meilisearch:", err)
	}

	for i := range documents {
		key := fmt.Sprintf("barang:%v", documents[i]["id"])

		if err := rds.Del(ctx, key).Err(); err != nil {
			log.Printf("⚠️ gagal hapus key lama %s: %v", key, err)
		}

		if err := rds.HSet(ctx, key, documents[i]).Err(); err != nil {
			log.Printf("❌ gagal HSET key %s: %v", key, err)
		}
	}

	log.Println("✅ Task UID:", task.TaskUID)

	fmt.Println("Barang Maintain Jalan")
	if err_buat_key := rds.SAdd(ctx, "barang_keys", "_init_").Err(); err_buat_key != nil {
		fmt.Println("Gagal Membuat keys semua barang")
	} else {
		var barang_induk []int32
		if err := db.Model(&models.BarangInduk{}).Pluck("id", &barang_induk).Error; err != nil {
			fmt.Println("Gagal Ambil id Semua Barang")
		} else {
			for _, data_id := range barang_induk {

				redisKey := fmt.Sprintf("barang:%v", data_id)

				if err := rds.SAdd(ctx, "barang_keys", redisKey).Err(); err != nil {
					fmt.Printf("❌ Gagal masukin %s ke barang_keys: %v\n", redisKey, err)
				} else {
					fmt.Printf("✅ Berhasil masukin %s ke barang_keys\n", redisKey)
				}
			}
		}
	}

	var idSeller []int32
	if err := db.Model(&models.Seller{}).Pluck("id", &idSeller).Error; err != nil {
		fmt.Println("❌ Gagal mendapatkan ID seluruh seller:", err)
		return
	}

	for _, id := range idSeller {
		key := fmt.Sprintf("barang_seller:%v", id)

		if err := rds.Del(ctx, key).Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus Redis key %s: %v\n", key, err)
		}

		if err := rds.SAdd(ctx, key, "_init").Err(); err != nil {
			fmt.Printf("❌ Gagal buat set Redis untuk seller %v: %v\n", id, err)
		} else {
			fmt.Printf("✅ Redis set siap untuk seller %v\n", id)
		}

		var idBarangInduk []int32
		if err := db.Model(&models.BarangInduk{}).
			Where(&models.BarangInduk{SellerID: id}).
			Pluck("id", &idBarangInduk).Error; err != nil {
			fmt.Println("❌ Gagal mendapatkan barang induk:", err)
		}

		for _, barangID := range idBarangInduk {
			if err := rds.SAdd(ctx, key, fmt.Sprintf("barang:%v", barangID)).Err(); err != nil {
				fmt.Printf("⚠️ Gagal tambah barang %v ke Redis untuk seller %v: %v\n", barangID, id, err)
			}
		}
	}

	jenisBarang := [...]string{
		"Pakaian&Fashion", "Kosmetik&Kecantikan", "Elektronik&Gadget",
		"Buku&Media", "Makanan&Minuman", "Ibu&Bayi", "Mainan",
		"Olahraga&Outdoor", "Otomotif&Sparepart", "RumahTangga",
		"AlatTulis", "Perhiasan&Aksesoris", "ProdukDigital",
		"Bangunan&Perkakas", "Musik&Instrumen", "Film&Broadcasting",
		"SemuaBarang",
	}

	for _, jenis := range jenisBarang {
		go func(j string) {
			key := fmt.Sprintf("jenis_%s_barang", j)

			if err := rds.Del(ctx, key).Err(); err != nil {
				fmt.Printf("⚠️ Gagal hapus Redis key %s: %v\n", key, err)
			}

			if err := rds.SAdd(ctx, key, "_init", 1).Err(); err != nil {
				fmt.Printf("❌ Gagal buat hash Redis untuk jenis %s: %v\n", j, err)
			} else {
				fmt.Printf("✅ Hash Redis siap untuk jenis %s\n", j)
			}

			var idBarangInduk []int32
			if err := db.Model(&models.BarangInduk{}).
				Where(&models.BarangInduk{JenisBarang: helper.ConvertJenisBarang(j)}).
				Pluck("id", &idBarangInduk).Error; err != nil {
				fmt.Println("❌ Gagal mendapatkan barang induk:", err)
			}

			for _, barangID := range idBarangInduk {
				if err := rds.SAdd(ctx, key, fmt.Sprintf("barang:%v", barangID)).Err(); err != nil {
					fmt.Printf("⚠️ Gagal tambah barang %v ke Redis untuk seller %v: %v\n", barangID, barangID, err)
				}
			}
		}(jenis)
	}
}

func InternalBarangMaintain(ctx context.Context, db *gorm.DB) {

	// EVALUATING STOK KATEGORI

	var BarangInduks []models.BarangInduk

	_ = db.Model(&models.BarangInduk{}).Find(&BarangInduks)

	if len(BarangInduks) == 0 {
		return
	}

	for _, bi := range BarangInduks {
		var KategoriBarangs []models.KategoriBarang

		if err := db.Model(&models.KategoriBarang{}).Where(&models.KategoriBarang{
			IdBarangInduk: bi.ID,
		}).Find(&KategoriBarangs).Error; err != nil {
			continue
		}

		for _, kb := range KategoriBarangs {
			var jumlah_real int64 = 0
			_ = db.Model(&models.VarianBarang{}).Where(&models.VarianBarang{
				IdBarangInduk: bi.ID,
				IdKategori:    kb.ID,
				Status:        "Ready",
			}).Count(&jumlah_real)

			if jumlah_real == 0 {
				continue
			}

			if jumlah_real == int64(kb.Stok) {
				continue
			}

			_ = db.Model(&kb).Updates(map[string]interface{}{
				"stok": jumlah_real})
		}
	}

	// EVALUATING STATUS VARIAN BARANG

	var varian_barang []models.VarianBarang

	_ = db.Model(&models.VarianBarang{}).Where(&models.VarianBarang{Status: "Down"}).Or(&models.VarianBarang{Status: "Ready"}).Find(&varian_barang)

	var idErrorDown []int64

	for _, vb := range varian_barang {
		if vb.HoldBy != 0 || vb.HolderEntity != "" {
			idErrorDown = append(idErrorDown, vb.ID)
		}
	}

	if len(idErrorDown) > 0 {
		if err := db.Model(&models.VarianBarang{}).
			Where("id IN ?", idErrorDown).
			Updates(map[string]interface{}{
				"hold_by":       0,
				"holder_entity": "",
			}).Error; err != nil {
			log.Println("Gagal memperbarui data varian:", err)
		} else {
			log.Printf("Berhasil mereset %d varian yang statusnya Down.\n", len(idErrorDown))
		}
	} else {
		log.Println("Tidak ada varian Down yang masih di-hold.")
	}
}
