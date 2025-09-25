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
			// log kondisi sebelum update
			fmt.Printf("📝 Preparing UPDATE VarianBarang | WHERE: {IdTransaksi:%d, Status:'Diproses', HoldBy:%d} | UPDATE: {Status:'Terjual'} | Limit=%d\n",
				data.ID, data.IdPengguna, data.Kuantitas)

			q := tx.Model(&models.VarianBarang{}).
				Where(&models.VarianBarang{
					IdTransaksi: data.ID,
					Status:      "Diproses",
					HoldBy:      data.IdPengguna,
				}).
				Limit(int(data.Kuantitas)).
				Updates(&models.VarianBarang{Status: "Terjual"})

			if q.Error != nil {
				// ❌ TRACE error
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
