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
		fmt.Printf("[ERROR] Gagal membuat queue notifikasi untuk user '%s' (ID: %v): %v\n", data.Username, data.ID, err)
	}

	notif.UserAccount("Bergabung", "Hai Selamat Bergabung di Burung!", nil)
	if err := notif.PublishMessageCritical(RoutingKey, conn); err != nil {
		fmt.Printf("[ERROR] Gagal publish pesan notifikasi ke user '%s' (ID: %v): %v\n", data.Username, data.ID, err)
	} else {
		fmt.Printf("[INFO] Notifikasi bergabung berhasil dikirim ke user '%s' (ID: %v)\n", data.Username, data.ID)
	}
}

func OnlinePengguna(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadPengguna, rds *redis.Client, conn *amqp091.Connection) {
	var notif notification.Notification
	fmt.Printf("[INFO] Menyimpan status online user '%s' (ID: %v) ke Redis\n", data.Username, data.ID)

	key := fmt.Sprintf("pengguna_online:%v", data.ID)

	fields := map[string]interface{}{
		"nama":     data.Nama,
		"username": data.Username,
		"email":    data.Email,
	}

	for field, value := range fields {
		if err := rds.HSet(ctx, key, field, value).Err(); err != nil {
			fmt.Printf("[ERROR] Gagal set field '%s' untuk user '%s' (ID: %v) di Redis: %v\n", field, data.Username, data.ID, err)
		}
	}

	notif.UserAccount("Login", "Kamu Telah login pada pukul dan jam .", nil)
	_, Routingkey := producer_mb.UserQueueRoutingKeyGenerate(data.Username, data.ID)
	if err := notif.PublishMessageCritical(Routingkey, conn); err != nil {
		fmt.Printf("[ERROR] Gagal mengirim notifikasi login ke user '%s' (ID: %v): %v\n", data.Username, data.ID, err)
	} else {
		fmt.Printf("[INFO] Notifikasi login berhasil dikirim ke user '%s' (ID: %v)\n", data.Username, data.ID)
	}
}

func OfflinePengguna(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsesPayloadPengguna, rds *redis.Client) {
	key := fmt.Sprintf("pengguna_online:%v", data.ID)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Printf("[ERROR] Gagal menghapus key Redis untuk user offline (key: %s): %v\n", key, err)
	} else {
		fmt.Printf("[INFO] User offline, key Redis dihapus: %s\n", key)
	}
}

func UpSeller(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsePayloadSeller, rds *redis.Client, conn *amqp091.Connection) {
	fmt.Printf("[INFO] Menyimpan data seller baru '%s' (ID: %v) ke Redis\n", data.Username, data.ID)

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
		fmt.Printf("[ERROR] Gagal menyimpan data seller '%s' (ID: %v) ke Redis: %v\n", data.Username, data.ID, err)
	}

	NamaQueue, RoutingKey := producer_mb.SellerQueueRoutingKeyGenerate(data.Username, data.ID)

	err := producer_mb.UpNotificationQueue(NamaQueue, RoutingKey, conn)
	if err != nil {
		fmt.Printf("[ERROR] Gagal membuat queue notifikasi untuk seller '%s' (ID: %v): %v\n", data.Username, data.ID, err)
	} else {
		fmt.Printf("[INFO] Queue notifikasi seller '%s' (ID: %v) berhasil dibuat\n", data.Username, data.ID)
	}

	var notif notification.Notification

	notif.SellerAccount("Bergabung", "Hai Selamat datang dan bergabung di burung ya", nil)

}

func HapusSeller(ctx context.Context, db *gorm.DB, data notify_payload.NotifyResponsePayloadSeller, rds *redis.Client) {
	key := fmt.Sprintf("seller_data:%v", data.ID)

	if err := rds.Del(ctx, key).Err(); err != nil {
		fmt.Printf("[ERROR] Gagal menghapus key Redis untuk seller (key: %s): %v\n", key, err)
	} else {
		fmt.Printf("[INFO] Seller dihapus, key Redis dihapus: %s\n", key)
	}
}
