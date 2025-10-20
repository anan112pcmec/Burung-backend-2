package maintain

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
			MaintainBarangInduk(ctx, db, rds, SE)
			CachingBarangInduk(ctx, db, rds, SE)
			KategoriBarangMaintain(ctx, db)
			VarianBarangMaintain(ctx, db)
			time.Sleep(10 * time.Minute)
		}
	}
}

type UpdateViewLikesBarangInduk struct {
	IdBarangInduk     string
	ViewedBarangInduk string
	LikesBarangInduk  string
}

func (u *UpdateViewLikesBarangInduk) Parse() (Id int, View int, Likes int, status bool) {
	status = true
	Id, err := strconv.Atoi(u.IdBarangInduk)
	if err != nil {
		status = false
	}

	View, err_v := strconv.Atoi(u.ViewedBarangInduk)
	if err_v != nil {
		status = false
	}

	Likes, err_l := strconv.Atoi(u.LikesBarangInduk)
	if err_l != nil {
		status = false
	}

	return
}

type UpdateKomentarBarangInduk struct {
	Id            int32
	TotalKomentar int32
}

type UpdateHargaBarangInduk struct {
	Id            int32
	HargaKategori int32
}

func MaintainBarangInduk(ctx context.Context, db *gorm.DB, rds *redis.Client, SE meilisearch.ServiceManager) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Evaluasi Likes dan Viewed Barang Dari Cache Ke Internal DB
	var KI_rds []UpdateViewLikesBarangInduk
	key := "barang_keys"

	// Ambil semua key dari Redis
	result, err := rds.SMembers(ctx, key).Result()
	if err != nil {
		fmt.Println("Gagal Mendapatkan Key Dari Redis:", err)
		goto MaintainDB
	}

	if len(result) == 0 {
		fmt.Println("Tidak ada key barang di Redis")
		goto MaintainDB
	}

	// Loop untuk ambil data tiap key barang
	for _, kb := range result {
		wg.Add(1)
		go func(key_barang string) {
			defer wg.Done()
			data, err_kb := rds.HGetAll(ctx, key_barang).Result()
			if err_kb != nil || len(data) == 0 {
				return
			}

			if data["id"] == "" {
				return
			}

			mu.Lock()
			KI_rds = append(KI_rds, UpdateViewLikesBarangInduk{
				IdBarangInduk:     data["id"],
				ViewedBarangInduk: data["viewed_barang_induk"],
				LikesBarangInduk:  data["likes_barang_induk"],
			})
			mu.Unlock()
		}(kb)
	}

	wg.Wait()

	// Jika setelah wait ternyata tidak ada data, langsung loncat ke maintain DB
	if len(KI_rds) == 0 {
		goto MaintainDB
	}

	// Update data barang induk di database
	for _, update_b := range KI_rds {
		Id, View, Likes, status := update_b.Parse()
		if !status {
			continue
		}

		if err := db.Model(&models.BarangInduk{}).Where("id = ?", Id).Updates(&models.BarangInduk{
			Viewed: int32(View),
			Likes:  int32(Likes),
		}).Error; err != nil {
			fmt.Println("Gagal Update Barang Induk Id:", Id, "-", err)
			continue
		}
	}

	// Lanjut ke tahap maintain komentar
	goto MaintainDB

MaintainDB:
	// Berfokus Memaintain Komentar dan Harga Kategori
	var KI_db []int64
	var UpdateTotalKomen []UpdateKomentarBarangInduk
	var UpdateHarga []UpdateHargaBarangInduk

	if err := db.Model(&models.BarangInduk{}).Pluck("id", &KI_db).Error; err != nil {
		goto MaintainSE
	}

	if len(KI_db) == 0 {
		goto MaintainSE
	}

	for _, Id := range KI_db {
		var total int64 = 0
		if err := db.Model(&models.Komentar{}).Where(&models.Komentar{IdBarangInduk: int32(Id)}).Count(&total).Error; err != nil {
			continue
		}

		if total == 0 {
			continue
		}

		UpdateTotalKomen = append(UpdateTotalKomen, UpdateKomentarBarangInduk{
			Id:            int32(Id),
			TotalKomentar: int32(total),
		})
	}

	for _, updateKomen := range UpdateTotalKomen {
		if err := db.Model(&models.BarangInduk{}).Where(&models.BarangInduk{ID: updateKomen.Id}).
			Update("total_komentar", updateKomen.TotalKomentar).Error; err != nil {
			fmt.Println("Gagal Update Komentar Barang Id:", updateKomen.Id, "-", err)
			continue
		}
	}

	// Maintain Harga

	for _, Id := range KI_db {
		var origin string = ""
		var harga int64 = 0
		if err := db.Model(&models.Komentar{}).Select("original_kategori").Where(&models.Komentar{IdBarangInduk: int32(Id)}).Take(&origin).Error; err != nil {
			continue
		}

		if harga == 0 || origin == "" {
			continue
		}

		if err := db.Model(&models.KategoriBarang{}).Select("harga").Where(&models.KategoriBarang{
			IdBarangInduk: int32(Id), Nama: origin,
		}).Take(&harga).Error; err != nil {
			continue
		}

		UpdateHarga = append(UpdateHarga, UpdateHargaBarangInduk{
			Id:            int32(Id),
			HargaKategori: int32(harga),
		})
	}

	for _, updateHarga := range UpdateHarga {
		if err := db.Model(&models.BarangInduk{}).Where(&models.BarangInduk{ID: updateHarga.Id}).
			Update("harga_kategoris", updateHarga.HargaKategori).Error; err != nil {
			continue
		}
	}

	goto MaintainSE

MaintainSE:

	dataBarangInduk := []models.BarangInduk{}

	// Ambil 100 barang paling populer dari DB
	if err := db.Model(&models.BarangInduk{}).
		Order("viewed DESC, likes DESC").
		Limit(100).
		Find(&dataBarangInduk).Error; err != nil {
		fmt.Println("❌ Gagal mengambil data Barang Induk:", err)
		return
	}

	if len(dataBarangInduk) == 0 {
		fmt.Println("⚠️ Tidak ada data barang untuk diindeks ke Meilisearch")
		return
	}

	// Pastikan index Meilisearch valid
	barangIndukIndex := SE.Index("barang_induk_all")
	if barangIndukIndex == nil {
		fmt.Println("❌ Index Meilisearch 'barang_induk_all' tidak ditemukan")
		return
	}

	// Siapkan dokumen untuk Meilisearch
	documents := make([]map[string]interface{}, 0, len(dataBarangInduk))

	for _, b := range dataBarangInduk {
		fmt.Println(" Barang:", b.NamaBarang)

		documents = append(documents, map[string]interface{}{
			"id":                          b.ID,
			"id_barang_induk":             b.ID,
			"nama_barang_induk":           b.NamaBarang,
			"id_seller_barang_induk":      b.SellerID,
			"original_kategori":           b.OriginalKategori,
			"deskripsi":                   b.Deskripsi,
			"jenis_barang_induk":          b.JenisBarang,
			"tanggal_rilis_barang_induk":  b.TanggalRilis,
			"viewed_barang_induk":         b.Viewed,
			"likes_barang_induk":          b.Likes,
			"total_komentar_barang_induk": b.TotalKomentar,
		})
	}

	task, err := barangIndukIndex.AddDocuments(documents, nil)
	if err != nil {
		fmt.Println("❌ Gagal menambahkan dokumen ke Meilisearch:", err)
		return
	}

	log.Println("✅ Task UID terkirim ke Meilisearch:", task.TaskUID)

}

func CachingBarangInduk(ctx context.Context, db *gorm.DB, rds *redis.Client, SE meilisearch.ServiceManager) {
	fmt.Println("🚀 Memulai proses maintain cache barang...")

	// ⚠️ Hati-hati, ini akan menghapus semua data Redis
	if _, err := rds.FlushDB(ctx).Result(); err != nil {
		fmt.Println("❌ Gagal melakukan flush Redis:", err)
	} else {
		fmt.Println("✅ Redis berhasil dibersihkan (FlushDB).")
	}

	key := "barang_keys"

	// Bersihkan key lama
	if result, err := rds.SMembers(ctx, key).Result(); err == nil {
		for _, keys := range result {
			if err := rds.Del(ctx, keys).Err(); err != nil {
				log.Printf("⚠️ Gagal hapus key lama %s: %v", keys, err)
			}
		}
		_ = rds.Del(ctx, key).Err()
	}

	dataBarangInduk := []models.BarangInduk{}

	// Ambil 100 barang paling populer dari DB
	if err := db.Model(&models.BarangInduk{}).
		Order("viewed DESC, likes DESC").
		Limit(100).
		Find(&dataBarangInduk).Error; err != nil {
		fmt.Println("❌ Gagal mengambil data Barang Induk:", err)
		return
	}

	if len(dataBarangInduk) == 0 {
		fmt.Println("⚠️ Tidak ada data barang untuk di-cache.")
		return
	}

	documents := make([]map[string]interface{}, 0, len(dataBarangInduk))

	for _, b := range dataBarangInduk {
		doc := map[string]interface{}{
			"id":                          b.ID,
			"id_barang_induk":             b.ID,
			"nama_barang_induk":           b.NamaBarang,
			"id_seller_barang_induk":      b.SellerID,
			"original_kategori":           b.OriginalKategori,
			"deskripsi":                   b.Deskripsi,
			"jenis_barang_induk":          b.JenisBarang,
			"tanggal_rilis_barang_induk":  b.TanggalRilis,
			"viewed_barang_induk":         b.Viewed,
			"likes_barang_induk":          b.Likes,
			"total_komentar_barang_induk": b.TotalKomentar,
			"harga":                       b.HargaKategoris,
		}
		documents = append(documents, doc)
	}

	// Simpan data barang ke Redis (hash per barang)
	for _, data := range documents {
		keyBarang := fmt.Sprintf("barang:%v", data["id"])
		if err := rds.HSet(ctx, keyBarang, data).Err(); err != nil {
			log.Printf("⚠️ Gagal menyimpan data ke Redis untuk %s: %v", keyBarang, err)
		}
	}

	fmt.Println("✅ Barang berhasil dimasukkan ke Redis")

	// Bangun ulang key utama barang_keys
	for _, data := range documents {
		keyBarang := fmt.Sprintf("barang:%v", data["id"])
		if err := rds.SAdd(ctx, "barang_keys", keyBarang).Err(); err != nil {
			fmt.Printf("❌ Gagal menambah %s ke barang_keys: %v\n", keyBarang, err)
		}
	}
	fmt.Println("✅ Redis key 'barang_keys' berhasil diperbarui")

	// Maintain per seller
	var idSeller []int32
	if err := db.Model(&models.Seller{}).Pluck("id", &idSeller).Error; err != nil {
		fmt.Println("❌ Gagal mendapatkan ID seluruh seller:", err)
		return
	}

	for _, id := range idSeller {
		keySeller := fmt.Sprintf("barang_seller:%v", id)

		_ = rds.Del(ctx, keySeller).Err()
		_ = rds.SAdd(ctx, keySeller, "_init").Err()

		var idBarangInduk []int32
		if err := db.Model(&models.BarangInduk{}).
			Where(&models.BarangInduk{SellerID: id}).
			Pluck("id", &idBarangInduk).Error; err != nil {
			fmt.Println("❌ Gagal mendapatkan barang induk untuk seller:", id)
			continue
		}

		for _, barangID := range idBarangInduk {
			if err := rds.SAdd(ctx, keySeller, fmt.Sprintf("barang:%v", barangID)).Err(); err != nil {
				fmt.Printf("⚠️ Gagal menambahkan barang %v ke Redis seller %v: %v\n", barangID, id, err)
			}
		}
	}

	fmt.Println("✅ Redis seller-barang mapping selesai")

	// Maintain jenis barang secara paralel
	jenisBarang := [...]string{
		"Pakaian&Fashion", "Kosmetik&Kecantikan", "Elektronik&Gadget",
		"Buku&Media", "Makanan&Minuman", "Ibu&Bayi", "Mainan",
		"Olahraga&Outdoor", "Otomotif&Sparepart", "RumahTangga",
		"AlatTulis", "Perhiasan&Aksesoris", "ProdukDigital",
		"Bangunan&Perkakas", "Musik&Instrumen", "Film&Broadcasting",
		"SemuaBarang",
	}

	var wg sync.WaitGroup
	for _, jenis := range jenisBarang {
		wg.Add(1)
		go func(j string) {
			defer wg.Done()

			keyJenis := fmt.Sprintf("jenis_%s_barang", j)
			_ = rds.Del(ctx, keyJenis).Err()
			_ = rds.SAdd(ctx, keyJenis, "_init").Err()

			var idBarangInduk []int32
			if err := db.Model(&models.BarangInduk{}).
				Where(&models.BarangInduk{JenisBarang: helper.ConvertJenisBarang(j)}).
				Pluck("id", &idBarangInduk).Error; err != nil {
				fmt.Println("❌ Gagal mendapatkan barang induk untuk jenis:", j)
				return
			}

			for _, barangID := range idBarangInduk {
				if err := rds.SAdd(ctx, keyJenis, fmt.Sprintf("barang:%v", barangID)).Err(); err != nil {
					fmt.Printf("⚠️ Gagal menambahkan barang %v ke Redis untuk jenis %s: %v\n", barangID, j, err)
				}
			}
			fmt.Printf("✅ Redis siap untuk jenis %s\n", j)
		}(jenis)
	}

	wg.Wait()

	fmt.Println("🎯 Caching barang selesai sepenuhnya.")
}

func KategoriBarangMaintain(ctx context.Context, db *gorm.DB) {

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
}

func VarianBarangMaintain(ctx context.Context, db *gorm.DB) {
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
