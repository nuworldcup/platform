package email

import (
	"html/template"
	"net/smtp"
	"os"
)

func sendEmail() error {
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "piotr@mailtrap.io", "extremely_secret_pass", "smtp.mailtrap.io")

	t, err := template.ParseFiles("teamRegistration.gohtml")
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string
	}{"John Smith"}

	err = t.Execute(os.Stdout, data)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"billy@microsoft.com"}
	msg := []byte("To: billy@microsoft.com\r\n" +
		"Subject: Why are you not using Mailtrap yet?\r\n" +
		"\r\n" +
		"Here’s the space for our great sales pitch\r\n")
	err = smtp.SendMail("smtp.mailtrap.io:25", auth, "piotr@mailtrap.io", to, msg)
	return err
}
