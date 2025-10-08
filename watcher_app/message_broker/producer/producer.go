package producer_mb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rabbitmq/amqp091-go"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Up Connection
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// :Berfungsi Saat Sistem Pertama Kali Jalan Akan Auto Membuat Koneksi Dan Beberapa Queue Default

func UpConnectionDefaults(username, password, port, exchange string) (*amqp091.Connection, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@localhost:%s", username, password, port)
	connection, err := amqp091.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect RabbitMQ: %w", err)
	}

	ch, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		connection.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	staticQueues := map[string]string{
		"global":                     "global",
		"notification_user_global":   "user.global",
		"notification_seller_global": "seller.global",
		"notification_kurir_global":  "kurir.global",
	}

	for qName, routingKey := range staticQueues {
		_, err := ch.QueueDeclare(qName, true, false, false, false, nil)
		if err != nil {
			ch.Close()
			connection.Close()
			return nil, fmt.Errorf("failed to declare queue %s: %w", qName, err)
		}

		err = ch.QueueBind(qName, routingKey, exchange, false, nil)
		if err != nil {
			ch.Close()
			connection.Close()
			return nil, fmt.Errorf("failed to bind queue %s: %w", qName, err)
		}
	}

	log.Println("✅ RabbitMQ channel, exchange, & queues ready")

	return connection, nil
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Up Queue Dan Down Queue
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// :Berfungsi Untuk Membuat Queue Baru Dan Menghapus Queue Yang Sudah AdA

func UpNotificationQueue(NamaQueue, RoutingKey string, conn *amqp091.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("gagal membuat channel RabbitMQ: %w", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(
		NamaQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("gagal mendeklarasikan queue %s: %w", NamaQueue, err)
	}

	exchange := helper.Getenvi("NOTIF_EXCHANGE", "gaada")
	if err := ch.QueueBind(
		NamaQueue,
		RoutingKey,
		exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("gagal mengikat queue %s ke exchange %s dengan routing key %s: %w",
			NamaQueue, exchange, RoutingKey, err)
	}

	return nil
}

func DownNotificationQueue(NamaQueue string, conn *amqp091.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("gagal membuka channel RabbitMQ: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDelete(
		NamaQueue,
		false,
		false,
		false,
	)
	if err != nil {
		return fmt.Errorf("gagal menghapus queue %s: %w", NamaQueue, err)
	}

	return nil
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Entity Queue And Routing Key
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// :Berfungsi Merilis Nama Queue dan Routing Key sesuai jenis Entity mereka bertujuan supaya semua sama rata dan
// mencegah boilerplate

func UserQueueRoutingKeyGenerate(username string, id int64) (NamaQueue string, RoutingKey string) {
	NamaQueue = fmt.Sprintf("notification_user_%v_%s", id, username)
	RoutingKey = fmt.Sprintf("user.%v", id)

	return
}

func SellerQueueRoutingKeyGenerate(username string, id int32) (NamaQueue string, RoutingKey string) {
	NamaQueue = fmt.Sprintf("notification_seller_%v_%s", id, username)
	RoutingKey = fmt.Sprintf("seller.%v", id)

	return
}

func KurirQueueRoutingKeyGenerate(username string, id int64) (NamaQueue string, RoutingKey string) {
	NamaQueue = fmt.Sprintf("notification_kurir_%v_%s", id, username)
	RoutingKey = fmt.Sprintf("kurir.%v", id)

	return
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Queue Check
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// :Berfungsi Untuk Mengecek apakah sebuah Queue sudah eksis dan routing key-nya.

func CheckQueueExists(NamaQueue string, conn *amqp091.Connection) (bool, string) {
	ch, err := conn.Channel()
	if err != nil {
		return false, ""
	}
	defer ch.Close()

	_, err = ch.QueueDeclarePassive(
		NamaQueue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		if amqpErr, ok := err.(*amqp091.Error); ok && amqpErr.Code == 404 {
			return false, ""
		}
		return false, ""
	}

	// Ambil binding info dari RabbitMQ Management API
	url := fmt.Sprintf("http://localhost:15672/api/queues/%%2F/%s/bindings", NamaQueue)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(
		helper.Getenvi("RMQ_USER", "guest"),
		helper.Getenvi("RMQ_PASS", "guest"),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return true, ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, ""
	}

	var bindings []QueueBinding
	if err := json.Unmarshal(body, &bindings); err != nil {
		return true, ""
	}

	if len(bindings) == 0 {
		return true, ""
	}

	for _, b := range bindings {
		if b.RoutingKey != "" {
			return true, b.RoutingKey
		}
	}

	return true, ""
}
