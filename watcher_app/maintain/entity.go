package maintain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
)

func EntityMaintainLoop(ctx context.Context, db *gorm.DB, rds *redis.Client, SE meilisearch.ServiceManager) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("❌ BarangMaintainLoop dihentikan")
			return
		default:
			MaintainSeller(ctx, db, rds, SE)
			time.Sleep(10 * time.Minute)
		}
	}
}

func MaintainSeller(ctx context.Context, db *gorm.DB, rds *redis.Client, SE meilisearch.ServiceManager) error {
	var sellersData []models.Seller
	if err := db.Find(&sellersData).Error; err != nil {
		return fmt.Errorf("gagal mengambil data seller dari DB: %w", err)
	}

	if len(sellersData) == 0 {
		fmt.Println("⚠️ Tidak ada data seller ditemukan")
		return nil
	}

	pipe := rds.Pipeline()
	for _, s := range sellersData {
		key := fmt.Sprintf("seller_data:%v", s.ID)
		cache := map[string]interface{}{
			"id_seller":                s.ID,
			"username_seller":          s.Username,
			"nama_seller":              s.Nama,
			"email_seller":             s.Email,
			"jam_operasional_seller":   s.JamOperasional,
			"seller_dedication_seller": s.SellerDedication,
			"jenis_seller":             s.Jenis,
			"punchline_seller":         s.Punchline,
			"deskripsi_seller":         s.Deskripsi,
			"follower_total_seller":    s.FollowerTotal,
		}
		pipe.HSet(ctx, key, cache)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("gagal menyimpan data seller ke Redis: %w", err)
	}

	dedicationList := []string{
		"Pakaian & Fashion", "Kosmetik & Kecantikan", "Elektronik & Gadget", "Buku & Media",
		"Makanan & Minuman", "Ibu & Bayi", "Mainan", "Olahraga & Outdoor",
		"Otomotif & Sparepart", "Rumah Tangga", "Alat Tulis", "Perhiasan & Aksesoris",
		"Produk Digital", "Bangunan & Perkakas", "Musik & Instrumen",
		"Film & Broadcasting", "Semua Barang",
	}

	jenisList := []string{"Brands", "Distributors", "Personal"}

	pipe = rds.Pipeline()
	for _, d := range dedicationList {
		setKey := fmt.Sprintf("seller_dedication:%s", d)
		pipe.Del(ctx, setKey)
		members := []interface{}{"_init_"}
		for _, s := range sellersData {
			if s.SellerDedication == d {
				members = append(members, fmt.Sprintf("seller_data:%v", s.ID))
			}
		}
		if len(members) > 0 {
			pipe.SAdd(ctx, setKey, members...)
		}
	}

	for _, j := range jenisList {
		setKey := fmt.Sprintf("seller_jenis:%s", j)
		pipe.Del(ctx, setKey)
		members := []interface{}{"_init_"}
		for _, s := range sellersData {
			if s.Jenis == j {
				members = append(members, fmt.Sprintf("seller_data:%v", s.ID))
			}
		}
		if len(members) > 0 {
			pipe.SAdd(ctx, setKey, members...)
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("gagal menyimpan kategori seller ke Redis: %w", err)
	}

	fmt.Printf("✅ Sinkronisasi %d seller ke Redis selesai\n", len(sellersData))

	SellerAll := SE.Index("seller_all")
	var dataSellerIndex []map[string]interface{}

	for _, d := range sellersData {
		dataSellerIndex = append(dataSellerIndex, map[string]interface{}{
			"id":                       d.ID,
			"nama_seller":              d.Nama,
			"jenis_seller":             d.Jenis,
			"seller_dedication_seller": d.SellerDedication,
		})
	}

	task, err := SellerAll.AddDocuments(dataSellerIndex, nil)
	if err != nil {
		log.Fatalf("Gagal Menambahkan data seller ke meilisearch")
	} else {
		log.Println("Berhasil Menambahkan data Seller ke meili search")
		log.Println(task)
	}

	return nil
}
