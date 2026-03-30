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

	// determine email type
	subject := "Verify Your Email"
	link := fmt.Sprintf("%s/verify-email?token=%s", env.AppBaseURL, job.Token)
	text := "Click the link below to verify your email:"

	if job.Type == "reset" {
		subject = "Reset Your Password"
		link = fmt.Sprintf("%s/reset-password?token=%s", env.AppBaseURL, job.Token)
		text = "Click the link below to reset your password:"
	}

	// Build email content
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
			"%s\n%s\n",
		from, to, subject, text, link,
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