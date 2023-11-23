package helper

import (
	"auth/models"
	"fmt"
	"log"

	"net/smtp"
)

func SendVerifEmail(user models.User, code string) {

	auth := smtp.PlainAuth("Marko Markovic", "soanosqlibmrs@gmail.com", "lqsiryrgbrjiofdz", "smtp.gmail.com")

	to := []string{*user.Email}

	subject := "Subject: Verify AirBnb clone account\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body><a href='http://localhost:4200/account/" + *user.Username + "/" + code + "'>Verify account</a></body></html>"

	msg := []byte(subject + mime + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, "Airbnb clone", to, msg)

	if err != nil {

		log.Fatal(err)
		fmt.Print(err)

	}

}

func SendVerifPasswordCode(email string, code string) {

	auth := smtp.PlainAuth("Marko Markovic", "soanosqlibmrs@gmail.com", "lqsiryrgbrjiofdz", "smtp.gmail.com")

	to := []string{email}

	subject := "Subject: Change AirBnb clone password\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body><h3>Your password recovery code is " + code + "</h3>" +
		"<br><h4>You have only 60 seconds to use this code before it expires!</h4></body></html>"

	msg := []byte(subject + mime + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, "Airbnb clone", to, msg)

	if err != nil {

		log.Fatal(err)
		fmt.Print(err)

	}

}
