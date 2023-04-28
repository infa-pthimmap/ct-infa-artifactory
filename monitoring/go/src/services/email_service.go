package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(subject string, emailBody string) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "sender@example.com", "password", "mailqa-useast1.cloudtrust.rocks")

	// Set up the message headers and body.
	to := []string{"pthimmappa@informatica.com"}
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		emailBody)

	// Send the email message.
	err := smtp.SendMail("mailqa-useast1.cloudtrust.rocks", auth, "noreply@informaticacloud.com", to, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
