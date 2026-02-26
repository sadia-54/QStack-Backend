package queue

import (
	"log"
)

func StartConsumer(worker func([]byte) error) {

	_, err := RabbitChannel.QueueDeclare(
		EmailQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	msgs, err := RabbitChannel.Consume(
		EmailQueueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}

	log.Println("Email worker listening for messages...")

	for msg := range msgs {
		err := worker(msg.Body)
		if err != nil {
			log.Println("Worker error:", err)
			msg.Nack(false, true)
			continue
		}
		msg.Ack(false)
	}
}