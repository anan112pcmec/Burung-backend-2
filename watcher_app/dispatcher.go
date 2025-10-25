package watcher_app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/meilisearch/meilisearch-go"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/services"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Fungsi Prosedur Pengguna Watcher
// Melihat dan mengawasi seluruh perubahan di dalam table pengguna
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
				go services.BarangMasuk(ctx, dbQuery, data, barangCache, SE)
			}

			if data.Action == "DELETE" {
				go services.HapusBarang(ctx, dbQuery, data, barangCache)
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

			var data notify_payload.NotifyResponseTransaksi
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("❌ Gagal Parse JSON:", err)
				continue
			}

			switch data.Action {
			case "UPDATE":
				if data.Status != "Dibatalkan" {
					go services.ApprovedTransaksiChange(data, dbQuery, conn)
				} else {
					go services.UnapproveTransaksiChange(data, dbQuery, conn)
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

func Informasi_Kurir_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB) {
	fmt.Println("Mengawasi Perubahan Informasi Kurir")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	listener := pq.NewListener(dsn, minReconn, maxReconn, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("informasi_kurir_channel"); err != nil {
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

			fmt.Printf("🔔 Dapat notify Informasi Kurir: %s\n", n.Extra)

			var data notify_payload.NotifyResponseInformasiKurir
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("❌ Gagal Parse JSON:", err)
				continue
			}

			switch data.Action {
			case "UPDATE":
				go services.VerifiedKurir(ctx, data.IdKurir, data.StatusPerizinan, data.JenisKendaraan, dbQuery)
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

func Informasi_Pengiriman_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB) {
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
				go services.WaitingConfirmation(data.IdTransaksi, data.IdKurir, data.Status, dbQuery)
			case "UPDATE":
				if data.Status == "Packaging" {
					go services.WaitingConfirmation(data.IdTransaksi, data.IdKurir, data.Status, dbQuery)
				}
				if data.Status == "Picked Up" {
					go services.DiperjalananConfirmation(data.IdTransaksi, data.Status, dbQuery)
				}

				if data.Status == "Sampai" {
					go services.SampaiConfirmation(data.IdTransaksi, data.Status, dbQuery)
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

func Follower_Watcher(ctx context.Context, dsn string, dbQuery *gorm.DB, entity_cache *redis.Client) {
	fmt.Println("Menjalankan Follower Watcher")

	minReconn := 10 * time.Second
	maxReconn := time.Minute

	listener := pq.NewListener(dsn, minReconn, maxReconn, func(event pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("[Listener Error] %v", err)
		}
	})

	if err := listener.Listen("follower_channel"); err != nil {
		log.Printf("❌ Gagal listen follower_channel: %v", err)
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

			fmt.Printf("🔔 Dapat notify informasi pengiriman: %s\n", n.Extra)

			var data notify_payload.NotifyResponseFollower
			if err := json.Unmarshal([]byte(n.Extra), &data); err != nil {
				fmt.Println("❌ Gagal parse JSON:", err)
				continue
			}

			switch data.Action {
			case "INSERT":
				go services.SellerFollowed(ctx, data, dbQuery, entity_cache)
			case "DELETE":
				go services.SellerUnfollowed(ctx, data, dbQuery, entity_cache)
			default:
				fmt.Println("⚠️ Aksi pengiriman tidak dikenali:", data.Action)
			}

		case <-ticker.C:
			if err := listener.Ping(); err != nil {
				log.Printf("[Ping Listener] Error: %v", err)
			}

		case <-ctx.Done():
			fmt.Println("🔴 Follower_Watcher dihentikan")
			return
		}
	}
}
