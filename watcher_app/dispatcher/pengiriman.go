package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
	kurir_pengiriman_watcher "github.com/anan112pcmec/Burung-backend-2/watcher_app/services/kurir/pengiriman_services"
	seller_order_processing_watcher "github.com/anan112pcmec/Burung-backend-2/watcher_app/services/seller/order_processing_services"
)

func Pengiriman_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB) {
	fmt.Println("Menjalankan Informasi Pengiriman Watcher")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	listener := pq.NewListener(dsn, minReconn, maxReconn, func(event pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("pengiriman_channel"); err != nil {
		log.Printf("❌ Gagal listen transaksi_channel: %v", err)
		return
	}

	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case n := <-listener.Notify:
			if n == nil {
				continue
			}

			fmt.Printf("🔔 Dapat notify Informasi Pengiriman: %s\n", n.Extra)

			var data notify_payload.NotifyResponsePengiriman
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("❌ Gagal Parse JSON:", err)
				continue
			}

			switch data.Action {
			case "INSERT":
				go seller_order_processing_watcher.WaitingConfirmation(data.IdTransaksi, data.IdKurir, data.Status, dbQuery)
			case "UPDATE":
				if data.Status == "Packaging" {
					go seller_order_processing_watcher.WaitingConfirmation(data.IdTransaksi, data.IdKurir, data.Status, dbQuery)
				}
				if data.Status == "Picked Up" {
					go kurir_pengiriman_watcher.DiperjalananConfirmation(data.IdTransaksi, data.Status, dbQuery)
				}

				if data.Status == "Sampai" {
					go kurir_pengiriman_watcher.SampaiConfirmation(data.IdTransaksi, data.Status, dbQuery)
				}
			case "DELETE":
			default:
				fmt.Println("⚠️ Aksi Pengiriman tidak dikenali:", data.Action)
			}

		case <-ticker.C:
			if err := listener.Ping(); err != nil {
				log.Printf("[Ping Listener] error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Pengiriman_Watcher dihentikan")
			return
		}
	}
}
