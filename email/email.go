package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
)

const templateForgotPassword = "templates/forgot_password.html"
const templateLayout = "templates/layout.html"

var (
	//go:embed templates/layout.html
	baseLayoutFS embed.FS
	//go:embed templates/forgot_password.html
	passwordTemplateFS embed.FS
	passwordTemplate *template.Template
)

func init() {
	baseLayout := template.Must(template.New("layout").ParseFS(baseLayoutFS, templateLayout))
	passwordTemplate = template.Must(baseLayout.ParseFS(passwordTemplateFS, templateForgotPassword))
}

type ForgotPasswordData struct {
	Link string
}

type Sender struct {
	Writer io.Writer
}

func (s Sender) SendForgotPasswordEmail(address string) error {
	// Execute template with data and store in a bytes.Buffer for use in email
	var body bytes.Buffer
	err := passwordTemplate.ExecuteTemplate(&body, "layout", &ForgotPasswordData{Link: "https://httpbin.org"})
	if err != nil {
		panic(err)
	}
	return s.sendEmail(address, "Reset Password", body.String())
}

func (s Sender) sendEmail(address, subject, body string) error {
	email := fmt.Sprintf("Receipient: %s\nSubject:%s\nBody:\n%s", address, subject, body)
	_, _ = s.Writer.Write([]byte(email))
	return nil
}