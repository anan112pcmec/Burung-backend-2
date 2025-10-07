package producer_mb

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"
)

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

func PublishMessageChannel(exchange, routingKey string, conn *amqp091.Connection, body []byte) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	return ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func UpNewNotificationQueue(NamaQueue, RoutingKey string, conn *amqp091.Connection) error {
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

func SellerQueueRoutingKeyGenerate(username string, id int32) (NamaQueue string, RoutingKey string) {
	NamaQueue = fmt.Sprintf("notification_seller_%v_%s", id, username)
	RoutingKey = fmt.Sprintf("seller.%v", id)

	return
}

func UserQueueRoutingKeyGenerate(username string, id int64) (NamaQueue string, RoutingKey string) {
	NamaQueue = fmt.Sprintf("notification_user_%v_%s", id, username)
	RoutingKey = fmt.Sprintf("user.%v", id)

	return
}

func KurirQueueRoutingKeyGenerate(username string, id int64) (NamaQueue string, RoutingKey string) {
	NamaQueue = fmt.Sprintf("notification_kurir_%v_%s", id, username)
	RoutingKey = fmt.Sprintf("kurir.%v", id)

	return
}
