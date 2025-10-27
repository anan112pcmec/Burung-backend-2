package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/services"
)

func Pengguna_Watcher(ctx context.Context, dsn string, db_query *gorm.DB, entity_cache *redis.Client, conn *amqp091.Connection) {
	// Objektif
	// 1.Watcher Pengguna memberikan notifikasi jika ada field sensitif yang berubah
	// (username, email, password, pin)

	fmt.Println("🟢 Mulai mengawasi pengguna_channel")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	listener := pq.NewListener(dsn, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("pengguna_channel"); err != nil {
		log.Fatalf("Gagal listen pengguna_channel: %v", err)
	}

	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case n := <-listener.Notify:
			if n != nil {
				fmt.Printf("🔔 Dapat notify: %s\n", n.Extra)
				var data notify_payload.NotifyResponsesPayloadPengguna
				err := json.Unmarshal([]byte(n.Extra), &data)
				if err != nil {
					fmt.Println("Gagal Parse JSON:", err)

				}

				if data.Action == "UPDATE" {

				}
			}

		case <-ticker.C:
			if err := listener.Ping(); err != nil {
				log.Printf("[Ping Listener] error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Entity_Watcher dihentikan")
			return
		}
	}
}

func Seller_Watcher(ctx context.Context, dsn string, db_query *gorm.DB, entity_cache *redis.Client, conn *amqp091.Connection) {
	fmt.Println("🟢 Mulai mengawasi seller_channel")

	minReconn := 10 & time.Second
	maxReconn := time.Minute

	listener_seller := pq.NewListener(dsn, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener_seller.Listen("seller_channel"); err != nil {
		log.Fatalf("Gagal listen seller_channel: %v", err)
	}

	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case n := <-listener_seller.Notify:
			if n != nil {
				fmt.Printf("🔔 Dapat notify: %s\n", n.Extra)
				var data notify_payload.NotifyResponsePayloadSeller
				err := json.Unmarshal([]byte(n.Extra), &data)
				if err != nil {
					fmt.Println("Gagal Parse JSON:", err)

				}

				if data.Action == "INSERT" {
					go services.UpSeller(ctx, db_query, data, entity_cache, conn)
				}

				if data.Action == "DELETE" {
					go services.HapusSeller(ctx, db_query, data, entity_cache)
				}
			}

		case <-ticker.C:
			if err := listener_seller.Ping(); err != nil {
				log.Printf("[Ping Listener] error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Entity_Watcher dihentikan")
			return
		}
	}
}
