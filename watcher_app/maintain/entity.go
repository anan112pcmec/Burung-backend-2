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
			MaintainSeller(ctx, db, SE)
			CachingSeller(ctx, db, rds)
			time.Sleep(10 * time.Minute)
		}
	}
}

func MaintainSeller(ctx context.Context, db *gorm.DB, SE meilisearch.ServiceManager) {
	// Mengambil Data Seller
	// Mengambil Seluruh data seller

	var sellersData []models.Seller

	_ = db.Model(&models.Seller{}).Find(&sellersData)

	if len(sellersData) == 0 {
		return
	}

	// Maintain Follower
	var FollowerSellerFinal []models.Seller

	for _, data := range sellersData {
		var totalFol int64
		if err := db.Model(&models.Follower{}).
			Where(models.Follower{IdFollowed: int64(data.ID)}).
			Count(&totalFol).Error; err != nil {
			log.Printf("gagal hitung follower untuk seller %d: %v", data.ID, err)
			continue
		}

		FollowerSellerFinal = append(FollowerSellerFinal, models.Seller{
			ID:            data.ID,
			FollowerTotal: int32(totalFol),
		})
	}

	for _, seller := range FollowerSellerFinal {
		if err := db.Model(&models.Seller{}).
			Where(&models.Seller{ID: seller.ID}).
			Update("follower_total", seller.FollowerTotal).Error; err != nil {
			log.Printf("gagal update follower_total untuk seller %d: %v", seller.ID, err)
		}
	}

	// Mengindex data seller yang di dapat ke search engine

	var sellersIndex []models.Seller

	// Mengambil Data Seller
	// Mengambil Seluruh data seller yang telah di maintain

	_ = db.Model(&models.Seller{}).Find(&sellersIndex)

	if len(sellersData) == 0 {
		return
	}

	SellerAll := SE.Index("seller_all")
	var dataSellerIndex []map[string]interface{}

	for _, d := range sellersIndex {
		dataSellerIndex = append(dataSellerIndex, map[string]interface{}{
			"id":                       d.ID,
			"id_seller":                d.ID,
			"nama_seller":              d.Nama,
			"jenis_seller":             d.Jenis,
			"seller_dedication_seller": d.SellerDedication,
			"follower_total_seller":    d.FollowerTotal,
		})
	}

	task, err := SellerAll.AddDocuments(dataSellerIndex, nil)
	if err != nil {
		log.Fatalf("Gagal Menambahkan data seller ke meilisearch")
	} else {
		log.Println("Berhasil Menambahkan data Seller ke meili search")
		log.Println(task)
	}

}

func CachingSeller(ctx context.Context, db *gorm.DB, rds *redis.Client) error {

	if _, err := rds.FlushDB(ctx).Result(); err != nil {
		fmt.Println("❌ Gagal melakukan flush Redis:", err)
	} else {
		fmt.Println("✅ Redis berhasil dibersihkan (FlushDB).")
	}

	var sellersData []models.Seller
	// Mengambil Data Seller
	// Hanya 100 seller dengan follower terbanyak yang akan di cache
	if err := db.Model(&models.Seller{}).Order("follower_total DESC").Limit(100).Find(&sellersData).Error; err != nil {
		return fmt.Errorf("gagal mengambil data seller dari DB: %w", err)
	}

	// return jika sellernya tidak ditemukan
	if len(sellersData) == 0 {
		fmt.Println("⚠️ Tidak ada data seller ditemukan")
		return nil
	}

	// Melakukan caching seller
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

	if err := rds.Del(ctx, "all_seller_keys").Err(); err != nil {
		log.Printf("⚠️ Gagal menghapus all_seller_keys: %v", err)
	}

	keys := make([]interface{}, 0, len(sellersData)+1)
	keys = append(keys, "_init_")
	for _, s := range sellersData {
		keys = append(keys, fmt.Sprintf("seller_data:%v", s.ID))
	}

	if err := rds.SAdd(ctx, "all_seller_keys", keys...).Err(); err != nil {
		log.Fatalf("❌ Gagal membuat all_seller_keys: %v", err)
	} else {
		log.Printf("✅ Berhasil membuat all_seller_keys (%d item)", len(keys))
	}

	fmt.Printf("✅ Sinkronisasi %d seller ke Redis selesai\n", len(sellersData))

	return nil
}
