package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
	seller_order_processing_watcher "github.com/anan112pcmec/Burung-backend-2/watcher_app/services/seller/order_processing_services"
)

func Transaksi_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB, conn *amqp091.Connection) {
	fmt.Println("Mengawasi Perubahan Transaksi")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	listener := pq.NewListener(dsn, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("transaksi_channel"); err != nil {
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

			fmt.Printf("🔔 Dapat notify Transaksi: %s\n", n.Extra)

			var data notify_payload.NotifyResponsePayloadTransaksi
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("❌ Gagal Parse JSON:", err)
				continue
			}

			switch data.Action {
			case "UPDATE":
				if data.Status != "Dibatalkan" {
					go seller_order_processing_watcher.ApprovedTransaksiChange(data, dbQuery, conn)
				} else {
					go seller_order_processing_watcher.UnapproveTransaksiChange(data, dbQuery, conn)
				}
			default:
				fmt.Println("⚠️ Aksi komentar tidak dikenali:", data.Action)
			}

		case <-ticker.C:
			if err := listener.Ping(); err != nil {
				log.Printf("[Ping Listener] error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Transaksi_Watcher dihentikan")
			return
		}
	}
}

func Pembayaran_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB, conn *amqp091.Connection) {
	fmt.Println("Mengawasi Perubahan Pembayaran")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	listener := pq.NewListener(dsn, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("pembayaran_channel"); err != nil {
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

			fmt.Printf("🔔 Dapat notify Transaksi: %s\n", n.Extra)

			var data notify_payload.NotifyResponsePayloadTransaksi
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("❌ Gagal Parse JSON:", err)
				continue
			}

		case <-ticker.C:
			if err := listener.Ping(); err != nil {
				log.Printf("[Ping Listener] error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Transaksi_Watcher dihentikan")
			return
		}
	}
}
