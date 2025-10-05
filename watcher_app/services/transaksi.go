package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"

)

func ApprovedTransaksiChange(data notify_payload.NotifyResponseTransaksi, db *gorm.DB) {
	start := time.Now()
	fmt.Printf("\n🔹 [START] ApprovedTransaksiChange | TransaksiID=%d | Status=%s | User=%d | Kuantitas=%d | Time=%s\n",
		data.ID, data.Status, data.IdPengguna, data.Kuantitas, start.Format(time.RFC3339))

	if err := db.Transaction(func(tx *gorm.DB) error {
		fmt.Printf("🚀 Transaction BEGIN | TransaksiID=%d\n", data.ID)

		if data.Status == "Diproses" {
			fmt.Printf("📝 Preparing UPDATE VarianBarang | WHERE: {IdTransaksi:%d, Status:'Diproses', HoldBy:%d} | UPDATE: {Status:'Terjual'} | Limit=%d\n",
				data.ID, data.IdPengguna, data.Kuantitas)

			var biayaongkir int16
			err_bk := tx.Model(&models.Ongkir{}).
				Where(&models.Ongkir{Nama: data.JenisPengiriman}).
				Select("value").Take(&biayaongkir).Error
			if err_bk != nil {
				fmt.Printf("❌ Gagal ambil biaya ongkir | TransaksiID=%d | Err=%v\n", data.ID, err_bk)
			}

			var id_kategori int64
			_ = tx.Model(&models.VarianBarang{}).
				Where(&models.VarianBarang{
					IdTransaksi: data.ID,
					Status:      "Diproses",
					HoldBy:      data.IdPengguna,
				}).
				Limit(1).
				Select("id_kategori").Take(&id_kategori).Error

			var kategorinya models.KategoriBarang
			_ = tx.Model(&models.KategoriBarang{}).
				Where(&models.KategoriBarang{ID: id_kategori}).
				Select("berat_gram", "dimensi_lebar_cm", "dimensi_panjang_cm", "id_alamat_gudang").
				Take(&kategorinya).Error

			beratTotalBarangPengirian := kategorinya.BeratGram * data.Kuantitas / 1000

			var biayalayanan int32
			var layanan string
			switch {
			case beratTotalBarangPengirian <= 10:
				layanan = "Motor"
				_ = tx.Model(&models.LayananPengirimanKurir{}).
					Where(&models.LayananPengirimanKurir{NamaLayanan: layanan}).
					Select("harga_layanan").Take(&biayalayanan).Error
			case beratTotalBarangPengirian <= 20:
				layanan = "Mobil"
				_ = tx.Model(&models.LayananPengirimanKurir{}).
					Where(&models.LayananPengirimanKurir{NamaLayanan: layanan}).
					Select("harga_layanan").Take(&biayalayanan).Error
			case beratTotalBarangPengirian <= 30:
				layanan = "Pickup"
				_ = tx.Model(&models.LayananPengirimanKurir{}).
					Where(&models.LayananPengirimanKurir{NamaLayanan: layanan}).
					Select("harga_layanan").Take(&biayalayanan).Error
			default:
				layanan = "Truk"
				_ = tx.Model(&models.LayananPengirimanKurir{}).
					Where(&models.LayananPengirimanKurir{NamaLayanan: layanan}).
					Select("harga_layanan").Take(&biayalayanan).Error
			}

			biayaKirim := biayaongkir - 5000
			kurirPaid := int32(biayaKirim) + biayalayanan

			pengiriman := models.Pengiriman{
				IdTransaksi:         data.ID,
				IdKurir:             0,
				NomorResi:           data.KodeOrder,
				Layanan:             layanan,
				JenisPengiriman:     data.JenisPengiriman,
				IdAlamatPengambilan: kategorinya.IDAlamat,
				IdAlamatPengiriman:  data.IdAlamat,
				Status:              "Packaging",
				BiayaKirim:          biayaKirim,
				KurirPaid:           kurirPaid,
				BeratTotalKG:        beratTotalBarangPengirian,
			}

			_ = tx.Create(&pengiriman)

			q := tx.Model(&models.VarianBarang{}).
				Where(&models.VarianBarang{
					IdTransaksi: data.ID,
					Status:      "Diproses",
					HoldBy:      data.IdPengguna,
				}).
				Limit(int(data.Kuantitas)).
				Updates(&models.VarianBarang{Status: "Terjual"})

			if q.Error != nil {
				fmt.Printf("❌ ERROR executing UPDATE | TransaksiID=%d | User=%d | Kuantitas=%d | Err=%v\n",
					data.ID, data.IdPengguna, data.Kuantitas, q.Error)
				return q.Error
			}

			if q.RowsAffected == 0 {
				fmt.Printf("⚠️ UPDATE executed but no rows affected | TransaksiID=%d | User=%d | Kuantitas=%d\n",
					data.ID, data.IdPengguna, data.Kuantitas)
			} else {
				fmt.Printf("✅ UPDATE success | TransaksiID=%d | RowsAffected=%d | User=%d | Kuantitas=%d\n",
					data.ID, q.RowsAffected, data.IdPengguna, data.Kuantitas)
			}

		} else {
			fmt.Printf("ℹ️ Status transaksi bukan 'Diproses' (Status=%s), tidak ada aksi update | TransaksiID=%d\n",
				data.Status, data.ID)
		}

		fmt.Printf("📌 Transaction about to COMMIT | TransaksiID=%d\n", data.ID)
		return nil
	}); err != nil {
		fmt.Printf("❌ Transaction ROLLBACK | TransaksiID=%d | Err=%v\n", data.ID, err)
	} else {
		fmt.Printf("✅ Transaction COMMIT | TransaksiID=%d\n", data.ID)
	}

	end := time.Now()
	fmt.Printf("🔹 [END] ApprovedTransaksiChange | TransaksiID=%d | Duration=%v ms\n\n",
		data.ID, end.Sub(start).Milliseconds())
}

func UnapproveTransaksiChange(data notify_payload.NotifyResponseTransaksi, db *gorm.DB) {
	var id_varian_barangs []int64

	if err_ambil_id := db.Model(&models.VarianBarang{}).
		Where(&models.VarianBarang{
			IdTransaksi:   data.ID,
			IdBarangInduk: data.IdBarangInduk,
		}).
		Limit(int(data.Kuantitas)). // batasi sesuai jumlah kuantitas
		Pluck("id", &id_varian_barangs).Error; err_ambil_id != nil {
		fmt.Println("Gagal Ambil Id:", err_ambil_id)
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, id_varian := range id_varian_barangs {
			if err_update := tx.Model(&models.VarianBarang{}).
				Where(models.VarianBarang{
					ID: id_varian,
				}).
				Updates(&models.VarianBarang{
					Status:       "Down",
					HoldBy:       0,
					HolderEntity: " ",
				}).Error; err_update != nil {
				return err_update
			}
		}
		return nil
	}); err != nil {
		fmt.Println("Gagal menjalankan Unapprove Transaksi Change:", err)
	}
}
