package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/sadia-54/qstack-backend/internal/config"
	"github.com/sadia-54/qstack-backend/internal/queue"
)

func EmailWorker(body []byte) error {
	// Load env 
	env := config.Load()

	// Deserialize RabbitMQ job
	var job queue.EmailJob
	if err := json.Unmarshal(body, &job); err != nil {
		return err
	}

	from := "no-reply@qstack.com"
	to := job.Email

	host := env.MailpitHost
	port := env.MailpitPort

	//Mailpit requires NO authentication
	var auth smtp.Auth = nil

	// Build email content
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Verify Your Email\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
			"Click the link below to verify your email:\n%s/verify-email?token=%s\n",
		from, to, env.AppBaseURL, job.Token,
	))

	// Send email
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, message)
	if err != nil {
		log.Println("Email send error:", err)
		return err
	}

	log.Println("Email sent to:", to)
	return nil
}