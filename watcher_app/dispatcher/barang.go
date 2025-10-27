package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
	seller_barang_watcher "github.com/anan112pcmec/Burung-backend-2/watcher_app/services/seller/barang_services"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Fungsi Prosedur Pengguna Watcher
// Melihat dan mengawasi seluruh perubahan di dalam table pengguna
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func Barang_Induk_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB, barangCache *redis.Client, SE meilisearch.ServiceManager) {
	fmt.Println("Mengawasi Perubahan Seluruh Data Barang Induk, Kategori, dan Varian Barang")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	// Listener ke Postgres
	listener := pq.NewListener(dsn, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("barang_induk_channel"); err != nil {
		log.Fatalf("Gagal listen barang_induk_channel: %v", err)
	}

	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case n := <-listener.Notify:
			if n == nil {
				continue
			}

			fmt.Printf("🔔 Dapat notify Barang: %s\n", n.Extra)

			var data notify_payload.NotifyResponsesPayloadBarang
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("Gagal Parse JSON:", err)
				continue
			}

			if data.Action == "INSERT" {
				go seller_barang_watcher.BarangMasuk(ctx, dbQuery, data, barangCache, SE)
			}

			if data.Action == "DELETE" {
				go seller_barang_watcher.HapusBarang(ctx, dbQuery, data, barangCache)
			}

		case <-ticker.C:
			if err := listener.Ping(); err != nil {
				log.Printf("[Ping Listener] error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Barang_Watcher dihentikan")
			return
		}
	}
}

func Varian_Barang_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB) {
	fmt.Println("Mengawasi Perubahan Seluruh Data Varian Barang, Kategori, dan Varian Barang")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	// Listener ke Postgres
	listener := pq.NewListener(dsn, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("varian_barang_channel"); err != nil {
		log.Fatalf("Gagal listen varian barang channel: %v", err)
	}

	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case n := <-listener.Notify:
			if n == nil {
				continue
			}

			fmt.Printf("🔔 Dapat notify Barang: %s\n", n.Extra)

			var data notify_payload.NotifyResponsePayloadVarianBarang
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("Gagal Parse JSON:", err)
				continue
			}

		case <-ticker.C:
			if err := listener.Ping(); err != nil {
				log.Printf("[Ping Listener] error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Varian_Barang_Watcher dihentikan")
			return
		}
	}
}
