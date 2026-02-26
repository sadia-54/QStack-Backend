package queue

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sadia-54/qstack-backend/internal/config"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func Connect() error {
	env := config.Load()
	url := env.RabbitURL

	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	RabbitConn = conn
	RabbitChannel = ch

	log.Println("RabbitMQ connected!")
	return nil
}

func Close() {
	if RabbitChannel != nil {
		RabbitChannel.Close()
	}
	if RabbitConn != nil {
		RabbitConn.Close()
	}
}