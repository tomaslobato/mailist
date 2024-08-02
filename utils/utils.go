package utils

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	DB "github.com/tomaslobato/mailist/db"
	"github.com/tomaslobato/mailist/models"
	"gopkg.in/gomail.v2"
)

func SendEmail(db *sql.DB, appPwd string, sendReq models.SendEmailReq) error {
	users, err := DB.GetUsers(db)
	if err != nil {
		return err
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, "tlobatodev@gmail.com", appPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", "tlobatodev@gmail.com")
	m.SetHeader("Subject", sendReq.Subject)
	m.SetBody("text/html", sendReq.Body)

	for _, user := range users {
		m.SetHeader("To", user.Email)
		err := d.DialAndSend(m)
		if err != nil {
			fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
			continue
		} else {
			fmt.Printf("Email sent to %s\n", user.Email)
		}
	}

	return nil
}
