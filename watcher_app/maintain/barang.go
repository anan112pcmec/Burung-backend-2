package maintain

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"
)

// convertJenisBarang akan mengubah nama jenis internal jadi format DB

func BarangMaintainLoop(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("❌ BarangMaintainLoop dihentikan")
			return
		default:
			BarangMaintain(ctx, db, rds)
			time.Sleep(10 * time.Minute)
		}
	}
}

func BarangMaintain(ctx context.Context, db *gorm.DB, rds *redis.Client) {
	var idSeller []int32

	// Ambil semua ID seller
	if err := db.Model(models.Seller{}).Pluck("id", &idSeller).Error; err != nil {
		fmt.Println("❌ Gagal mendapatkan ID seluruh seller:", err)
		return
	}

	// Maintain key barang_seller
	// Maintain key barang_seller
	for _, id := range idSeller {
		key := fmt.Sprintf("barang_seller:%v", id)

		// Hapus key lama
		if err := rds.Del(ctx, key).Err(); err != nil {
			fmt.Printf("⚠️ Gagal hapus Redis key %s: %v\n", key, err)
		}

		// Buat inisialisasi set kosong dengan marker
		if err := rds.SAdd(ctx, key, "_init").Err(); err != nil {
			fmt.Printf("❌ Gagal buat set Redis untuk seller %v: %v\n", id, err)
		} else {
			fmt.Printf("✅ Redis set siap untuk seller %v\n", id)
		}

		// Ambil barang induk
		var idBarangInduk []int32
		if err := db.Model(models.BarangInduk{}).
			Where(models.BarangInduk{SellerID: id}).
			Pluck("id", &idBarangInduk).Error; err != nil {
			fmt.Println("❌ Gagal mendapatkan barang induk:", err)
		}

		// Simpan ke Redis set
		for _, barangID := range idBarangInduk {
			if err := rds.SAdd(ctx, key, fmt.Sprintf("barang:%v", barangID)).Err(); err != nil {
				fmt.Printf("⚠️ Gagal tambah barang %v ke Redis untuk seller %v: %v\n", barangID, id, err)
			}
		}
	}

	// Selesai update barang berdasarkan seller
	jenisBarang := [...]string{
		"Pakaian&Fashion", "Kosmetik&Kecantikan", "Elektronik&Gadget",
		"Buku&Media", "Makanan&Minuman", "Ibu&Bayi", "Mainan",
		"Olahraga&Outdoor", "Otomotif&Sparepart", "RumahTangga",
		"AlatTulis", "Perhiasan&Aksesoris", "ProdukDigital",
		"Bangunan&Perkakas", "Musik&Instrumen", "Film&Broadcasting",
		"SemuaBarang",
	}

	// Maintain key jenis barang
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
			if err := db.Model(models.BarangInduk{}).
				Where(models.BarangInduk{JenisBarang: helper.ConvertJenisBarang(j)}).
				Pluck("id", &idBarangInduk).Error; err != nil {
				fmt.Println("❌ Gagal mendapatkan barang induk:", err)
			}

			// Simpan ke Redis set
			for _, barangID := range idBarangInduk {
				if err := rds.SAdd(ctx, key, fmt.Sprintf("barang:%v", barangID)).Err(); err != nil {
					fmt.Printf("⚠️ Gagal tambah barang %v ke Redis untuk seller %v: %v\n", barangID, barangID, err)
				}
			}
		}(jenis)
	}
}
