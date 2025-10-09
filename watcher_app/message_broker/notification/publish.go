package notification

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/helper"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Publish Message
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// :Berfungsi Publish Message Ke Routing Key Tertentu, dan LevelMessage

func (n *Notification) PublishMessageCritical(routingKey string, conn *amqp091.Connection) error {
	exchange := helper.Getenvi("RMQ_NOTIF_EXCHANGE", "NULL")
	if n.Level != "critical" {
		return fmt.Errorf("hanya pesan critical yang Boleh dipublish dengan method ini")
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}

	return ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
}

func (n *Notification) PublishMessage(exchange, routingKey string, conn *amqp091.Connection) error {
	if n.Level == "critical" {
		return fmt.Errorf("pesan critical tidak boleh di kirim dengan method ini")
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}

	return ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
