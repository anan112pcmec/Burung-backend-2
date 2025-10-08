package services

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/message_broker/notification"
	producer_mb "github.com/anan112pcmec/Burung-backend-2/watcher_app/message_broker/producer"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

func UpUser(ctx context.Context, data notify_payload.NotifyResponsesPayloadPengguna, conn *amqp091.Connection) {
	var notif notification.Notification
	NamaQueue, RoutingKey := producer_mb.UserQueueRoutingKeyGenerate(data.Username, data.ID)

	err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn)
	if err != nil {
		fmt.Println("Gagal Up New Notification")
		fmt.Println(err)
	}

	notif.UserAccount("Bergabung", "Hai Selamat Bergabung Di Burung ya", nil)
	if err := notif.PublishMessageCritical(RoutingKey, conn); err != nil {
		fmt.Println("Gagal Publish Pesan Pada User: ", data.Nama, err)
	}
}

func OnlinePengguna(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadPengguna, rds *redis.Client, conn *amqp091.Connection) {
	var notif notification.Notification
	fmt.Println("Caching Online User")

	key := fmt.Sprintf("pengguna_online:%v", data.ID)

	fields := map[string]interface{}{
		"nama":     data.Nama,
		"username": data.Username,
		"email":    data.Email,
	}

	for field, value := range fields {
		if err := rds.HSet(ctx, key, field, value).Err(); err != nil {
			fmt.Println("Gagal Set Redis:", err)
		}
	}

	notif.UserAccount("Login", "Kamu Login ya", nil)
	_, Routingkey := producer_mb.UserQueueRoutingKeyGenerate(data.Username, data.ID)
	if err := notif.PublishMessageCritical(Routingkey, conn); err != nil {
		fmt.Println("Gagal Notif User Login")
	}
}

func OfflinePengguna(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadPengguna, rds *redis.Client) {
	key := fmt.Sprintf("pengguna_online:%v", data.ID)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Println("Gagal Hapus Redis Key:", err)
	} else {
		fmt.Println("✅ User offline, key dihapus:", key)
	}
}

func UpSeller(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsePayloadSeller, rds *redis.Client, conn *amqp091.Connection) {
	fmt.Println("Caching NEW Seller")

	key := fmt.Sprintf("seller_data:%v", data.ID)

	var fields = map[string]interface{}{
		"id_seller":                data.ID,
		"username_seller":          data.Username,
		"nama_seller":              data.Nama,
		"email_seller":             data.Email,
		"jam_operasional_seller":   data.JamOperasional,
		"seller_dedication_seller": data.SellerDedication,
		"jenis_seller":             data.Jenis,
		"punchline_seller":         data.Punchline,
		"deskripsi_seller":         data.Deskripsi,
		"follower_total_seller":    data.FollowerTotal,
	}

	if err := rds.HSet(ctx, key, fields).Err(); err != nil {
		fmt.Println("Gagal Set Redis:", err)
	}

	NamaQueue, RoutingKey := producer_mb.SellerQueueRoutingKeyGenerate(data.Username, data.ID)

	err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn)
	if err != nil {
		fmt.Println("Gagal Up New Notification")
	}

}

func HapusSeller(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsePayloadSeller, rds *redis.Client) {
	key := fmt.Sprintf("seller_data:%v", data.ID)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Println("Gagal Hapus Redis Key:", err)
	} else {
		fmt.Println("✅ seller dihapus, key dihapus:", key)
	}
}
