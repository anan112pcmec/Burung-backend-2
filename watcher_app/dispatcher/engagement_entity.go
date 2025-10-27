package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
	kurir_informasi_watcher "github.com/anan112pcmec/Burung-backend-2/watcher_app/services/kurir/informasi_services"
	pengguna_social_media_watcher "github.com/anan112pcmec/Burung-backend-2/watcher_app/services/pengguna/social_media_services"
)

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
				go kurir_informasi_watcher.VerifiedKurir(ctx, data.IdKurir, data.StatusPerizinan, data.JenisKendaraan, dbQuery)
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
				go pengguna_social_media_watcher.SellerFollowed(ctx, data, dbQuery, entity_cache)
			case "DELETE":
				go pengguna_social_media_watcher.SellerUnfollowed(ctx, data, dbQuery, entity_cache)
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
