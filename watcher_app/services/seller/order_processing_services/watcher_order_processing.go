package seller_order_processing_watcher

import (
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/message_broker/notification"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

func ApprovedTransaksiChange(data notify_payload.NotifyResponsePayloadTransaksi, db *gorm.DB, conn *amqp091.Connection) {
	var notif notification.Notification
	start := time.Now()
	fmt.Printf("\n[INFO] [START] ApprovedTransaksiChange | TransaksiID=%d | Status=%s | User=%d | Kuantitas=%d | Time=%s\n",
		data.ID, data.Status, data.IdPengguna, data.Kuantitas, start.Format(time.RFC3339))

	if err := db.Transaction(func(tx *gorm.DB) error {
		fmt.Printf("[INFO] Transaction BEGIN | TransaksiID=%d\n", data.ID)

		if data.Status == "Diproses" {
			fmt.Printf("[INFO] Preparing UPDATE VarianBarang | WHERE: {IdTransaksi:%d, Status:'Diproses', HoldBy:%d} | UPDATE: {Status:'Terjual'} | Limit=%d\n",
				data.ID, data.IdPengguna, data.Kuantitas)

			var biayaongkir int16
			err_bk := tx.Model(&models.Ongkir{}).
				Where(&models.Ongkir{Nama: data.JenisPengiriman}).
				Select("value").Take(&biayaongkir).Error
			if err_bk != nil {
				fmt.Printf("[ERROR] Gagal ambil biaya ongkir | TransaksiID=%d | Err=%v\n", data.ID, err_bk)
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

			if err := tx.Create(&pengiriman).Error; err != nil {
				go func() {
					_ = db.Model(&models.Transaksi{}).Where(&models.Transaksi{
						ID: data.ID,
					}).Update("status", "Dibayar")
				}()
				fmt.Printf("[ERROR] Gagal membuat data pengiriman | TransaksiID=%d | Err=%v\n", data.ID, err)
				return err
			}

			q := tx.Model(&models.VarianBarang{}).
				Where(&models.VarianBarang{
					IdTransaksi: data.ID,
					Status:      "Diproses",
					HoldBy:      data.IdPengguna,
				}).
				Limit(int(data.Kuantitas)).
				Updates(&models.VarianBarang{Status: "Terjual"})

			if q.Error != nil {
				fmt.Printf("[ERROR] Gagal update status VarianBarang | TransaksiID=%d | User=%d | Kuantitas=%d | Err=%v\n",
					data.ID, data.IdPengguna, data.Kuantitas, q.Error)
				return q.Error
			}

			if q.RowsAffected == 0 {
				fmt.Printf("[WARN] UPDATE executed but no rows affected | TransaksiID=%d | User=%d | Kuantitas=%d\n",
					data.ID, data.IdPengguna, data.Kuantitas)
			} else {
				fmt.Printf("[INFO] UPDATE success | TransaksiID=%d | RowsAffected=%d | User=%d | Kuantitas=%d\n",
					data.ID, q.RowsAffected, data.IdPengguna, data.Kuantitas)
			}

		} else {
			fmt.Printf("[INFO] Status transaksi bukan 'Diproses' (Status=%s), tidak ada aksi update | TransaksiID=%d\n",
				data.Status, data.ID)
		}

		fmt.Printf("[INFO] Transaction about to COMMIT | TransaksiID=%d\n", data.ID)
		return nil
	}); err != nil {
		fmt.Printf("[ERROR] Transaction ROLLBACK | TransaksiID=%d | Err=%v\n", data.ID, err)
	} else {
		fmt.Printf("[INFO] Transaction COMMIT | TransaksiID=%d\n", data.ID)
		notif.UserTransaksi("Pesanan", "Seller telah approve pesananmu.", nil)
		if err := notif.PublishMessage(helper.Getenvi("RMQ_NOTIF_EXCHANGE", "NIL"), fmt.Sprintf("user.%v", data.IdPengguna), conn); err != nil {
			fmt.Printf("[ERROR] Gagal mengirim notifikasi approve transaksi ke user (ID: %v): %v\n", data.IdPengguna, err)
		} else {
			fmt.Printf("[INFO] Notifikasi approve transaksi berhasil dikirim ke user (ID: %v)\n", data.IdPengguna)
		}
	}

	end := time.Now()
	fmt.Printf("[INFO] [END] ApprovedTransaksiChange | TransaksiID=%d | Duration=%v ms\n\n",
		data.ID, end.Sub(start).Milliseconds())
}

func UnapproveTransaksiChange(data notify_payload.NotifyResponsePayloadTransaksi, db *gorm.DB, conn *amqp091.Connection) {
	var id_varian_barangs []int64

	fmt.Printf("[INFO] [START] UnapproveTransaksiChange | TransaksiID=%d | User=%d | Kuantitas=%d\n", data.ID, data.IdPengguna, data.Kuantitas)

	if err_ambil_id := db.Model(&models.VarianBarang{}).
		Where(&models.VarianBarang{
			IdTransaksi:   data.ID,
			IdBarangInduk: data.IdBarangInduk,
		}).
		Limit(int(data.Kuantitas)).
		Pluck("id", &id_varian_barangs).Error; err_ambil_id != nil {
		fmt.Printf("[ERROR] Gagal ambil ID varian barang untuk unapprove transaksi | TransaksiID=%d | Err=%v\n", data.ID, err_ambil_id)
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, id_varian := range id_varian_barangs {
			if err_update := tx.Model(&models.VarianBarang{}).
				Where(models.VarianBarang{
					ID: id_varian,
				}).
				Updates(map[string]interface{}{
					"status":        "Down",
					"hold_by":       0,
					"holder_entity": nil,
				}).Error; err_update != nil {
				fmt.Printf("[ERROR] Gagal update status varian barang (ID: %d) | TransaksiID=%d | Err=%v\n", id_varian, data.ID, err_update)
				return err_update
			}
		}
		return nil
	}); err != nil {
		fmt.Printf("[ERROR] Gagal menjalankan Unapprove Transaksi Change | TransaksiID=%d | Err=%v\n", data.ID, err)
	} else {
		// Kirim notifikasi ke user
		var notif notification.Notification
		notif.UserTransaksi("Pesanan", "Transaksi kamu tidak di-approve oleh seller.", nil)
		if err := notif.PublishMessage(helper.Getenvi("RMQ_NOTIF_EXCHANGE", "NIL"), fmt.Sprintf("user.%v", data.IdPengguna), conn); err != nil {
			fmt.Printf("[ERROR] Gagal mengirim notifikasi unapprove transaksi ke user (ID: %v): %v\n", data.IdPengguna, err)
		} else {
			fmt.Printf("[INFO] Notifikasi unapprove transaksi berhasil dikirim ke user (ID: %v)\n", data.IdPengguna)
		}
	}
	fmt.Printf("[INFO] [END] UnapproveTransaksiChange | TransaksiID=%d\n", data.ID)
}

func WaitingConfirmation(IdTransaksi, IdKurir int64, status string, db *gorm.DB) {
	if status != "Packaging" {
		return
	}

	if IdTransaksi == 0 {
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err_update_transaksi := tx.Model(&models.Transaksi{}).Where(&models.Transaksi{
			ID: IdTransaksi,
		}).Update("status", "Waiting").Error; err_update_transaksi != nil {
			return err_update_transaksi
		}

		return nil
	}); err != nil {
		return
	}
}
