package maintain_mb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	producer_mb "github.com/anan112pcmec/Burung-backend-2/watcher_app/message_broker/producer"
)

func NotificationMaintainLoop(ctx context.Context, db *gorm.DB, conn *amqp091.Connection, Exchange string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("❌ Entity Maintain Notification dihentikan")
			return
		default:
			MaintainEntityNotifQueue(ctx, db, conn, Exchange)
			time.Sleep(10 * time.Minute)
		}
	}
}

func MaintainEntityNotifQueue(ctx context.Context, db *gorm.DB, conn *amqp091.Connection, Exchange string) {
	var (
		users   []models.Pengguna
		sellers []models.Seller
		kurirs  []models.Kurir
		wg      sync.WaitGroup
	)

	_ = db.Model(&models.Pengguna{}).Select("id", "username").Find(&users)
	_ = db.Model(&models.Seller{}).Select("id", "username").Find(&sellers)
	_ = db.Model(&models.Kurir{}).Select("id", "username").Find(&kurirs)

	if len(users) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, u := range users {
				var NamaQueue string = ""
				var RoutingKey string = ""

				NamaQueue, RoutingKey = producer_mb.UserQueueRoutingKeyGenerate(u.Username, u.ID)

				if NamaQueue == "" && RoutingKey == "" {
					continue
				}

				status, routingkeydefault := producer_mb.CheckQueueExists(NamaQueue, conn)

				if !status {
					if err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn); err != nil {
						fmt.Println("Queue Ini Gagal Di up:", NamaQueue)
					}

					continue
				}

				if routingkeydefault != RoutingKey {
					if err := producer_mb.DownNotificationQueue(NamaQueue, conn); err != nil {
						continue
					}

					if err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn); err != nil {
						fmt.Println("Queue Ini Gagal Di up:", NamaQueue)
					}
				}

			}
		}()
	}

	if len(sellers) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, u := range sellers {
				var NamaQueue string = ""
				var RoutingKey string = ""

				NamaQueue, RoutingKey = producer_mb.SellerQueueRoutingKeyGenerate(u.Username, u.ID)

				if NamaQueue == "" && RoutingKey == "" {
					continue
				}

				status, routingkeydefault := producer_mb.CheckQueueExists(NamaQueue, conn)

				if !status {
					if err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn); err != nil {
						fmt.Println("Queue Ini Gagal Di up:", NamaQueue)
					}

					continue
				}

				if routingkeydefault != RoutingKey {
					if err := producer_mb.DownNotificationQueue(NamaQueue, conn); err != nil {
						continue
					}

					if err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn); err != nil {
						fmt.Println("Queue Ini Gagal Di up:", NamaQueue)
					}
				}

			}
		}()
	}

	if len(kurirs) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, u := range kurirs {
				var NamaQueue string = ""
				var RoutingKey string = ""

				NamaQueue, RoutingKey = producer_mb.KurirQueueRoutingKeyGenerate(u.Username, u.ID)

				if NamaQueue == "" && RoutingKey == "" {
					continue
				}

				status, routingkeydefault := producer_mb.CheckQueueExists(NamaQueue, conn)

				if !status {
					if err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn); err != nil {
						fmt.Println("Queue Ini Gagal Di up:", NamaQueue)
					}

					continue
				}

				if routingkeydefault != RoutingKey {
					if err := producer_mb.DownNotificationQueue(NamaQueue, conn); err != nil {
						continue
					}

					if err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn); err != nil {
						fmt.Println("Queue Ini Gagal Di up:", NamaQueue)
					}
				}

			}
		}()
	}

	wg.Wait()

	fmt.Println("Berhasil Maintenance Entity NotifQueue")
}
