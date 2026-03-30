package queue

import (
	"context"
	"encoding/json"
	"time"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const EmailQueueName = "email_verification_queue"

type EmailJob struct {
	Email string `json:"email"`
	Token string `json:"token"`
	Type  string `json:"type"` // verify or reset
}

func PublishEmailVerification(email, token string) error {
	body, _ := json.Marshal(EmailJob{Email: email, Token: token})

	_, err := RabbitChannel.QueueDeclare(
		EmailQueueName,
		true,  // durable
		false, // auto delete
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = RabbitChannel.PublishWithContext(ctx,
		"", // direct queue
		EmailQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}

	log.Println("Email verification job published")
	return nil
}

func PublishPasswordReset(email, token string) error {

	body, _ := json.Marshal(EmailJob{
		Email: email,
		Token: token,
		Type:  "reset",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return RabbitChannel.PublishWithContext(
		ctx,
		"",
		EmailQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}