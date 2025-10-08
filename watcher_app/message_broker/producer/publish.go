package producer_mb

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Publish Message
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// :Berfungsi Publish Message Ke Routing Key Tertentu

// func PublishMessage(exchange, routingKey string, conn *amqp091.Connection) error {
// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return fmt.Errorf("failed to create channel: %w", err)
// 	}
// 	defer ch.Close()

// 	body, err := json.Marshal(event)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal event to JSON: %w", err)
// 	}

// 	return ch.Publish(
// 		exchange,
// 		routingKey,
// 		false, // mandatory
// 		false, // immediate
// 		amqp091.Publishing{
// 			ContentType: "application/json", // penting: json
// 			Body:        body,
// 		},
// 	)
// }
