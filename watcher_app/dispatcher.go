package watcher_app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/services"
)

func Entity_Watcher(ctx context.Context, dsn string, db_query *gorm.DB, entity_cache *redis.Client) {
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
					if data.ChangedColumns.Status == "Online" {
						go services.OnlinePengguna(ctx, db_query, data, entity_cache)
					} else if data.ChangedColumns.Status == "Offline" {
						go services.OfflinePengguna(ctx, db_query, data, entity_cache)
					}
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

func Barang_Watcher(ctx context.Context, dsn string, db_query *gorm.DB, barang_cache *redis.Client) {

	fmt.Println("Mengawasi Perubahan Seluruh Data Barang Induk, Kategori, dan Varian Barang")
}
