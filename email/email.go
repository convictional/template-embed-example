package email

import (
	"bytes"
	"fmt"
	"html/template"
)

const templateForgotPassword = "../templates/forgot_password.html"
const templateLayout = "../templates/layout.html"

type ForgotPasswordData struct {
	Link string
}

func SendForgotPasswordEmail(address string) error {
	// Read in template (in this example we are sending a forgot password email)
	passwordTemplate, err := template.ParseFiles(templateLayout, templateForgotPassword)
	if err != nil {
		panic(err)
	}

	// Execute template with data and store in a bytes.Buffer for use in email
	var body bytes.Buffer
	err = passwordTemplate.ExecuteTemplate(&body, "layout", &ForgotPasswordData{Link: "https://httpbin.org"})
	if err != nil {
		panic(err)
	}
	return sendEmail(address, "Reset Password", body.String())
}

func sendEmail(address, subject, body string) error {
	fmt.Printf("Receipient: %s\nSubject:%s\nBody:\n%s", address, subject, body)
	return nil
}