package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"runtime"
)

const templateForgotPassword = "templates/forgot_password.html"
const templateLayout = "templates/layout.html"

var	passwordTemplateInit *template.Template

var (
	//go:embed templates/layout.html
	baseLayoutFS embed.FS
	//go:embed templates/forgot_password.html
	passwordTemplateFS embed.FS
	passwordTemplateInitFS *template.Template
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func init() {
	passwordTemplateInit = template.Must(template.ParseFiles(fmt.Sprintf("%s/%s", basepath, templateLayout), fmt.Sprintf("%s/%s", basepath, templateForgotPassword)))
	baseLayout := template.Must(template.New("layout").ParseFS(baseLayoutFS, templateLayout))
	passwordTemplateInitFS = template.Must(baseLayout.ParseFS(passwordTemplateFS, templateForgotPassword))
}

type ForgotPasswordData struct {
	Link string
}

type Sender struct {
	Writer io.Writer
}

func (s Sender) SendForgotPasswordEmailInitFS(address string) error {
	// Execute template with data and store in a bytes.Buffer for use in email
	var body bytes.Buffer
	err := passwordTemplateInitFS.ExecuteTemplate(&body, "layout", &ForgotPasswordData{Link: "https://httpbin.org"})
	if err != nil {
		panic(err)
	}
	return s.sendEmail(address, "Reset Password", body.String())
}

func (s Sender) SendForgotPasswordEmailInit(address string) error {
	// Execute template with data and store in a bytes.Buffer for use in email
	var body bytes.Buffer
	err := passwordTemplateInit.ExecuteTemplate(&body, "layout", &ForgotPasswordData{Link: "https://httpbin.org"})
	if err != nil {
		panic(err)
	}
	return s.sendEmail(address, "Reset Password", body.String())
}

func (s Sender) SendForgotPasswordEmail(address string) error {
	// Read in template (in this example we are sending a forgot password email)
	passwordTemplate, err := template.ParseFiles(fmt.Sprintf("%s/%s", basepath, templateLayout), fmt.Sprintf("%s/%s", basepath, templateForgotPassword))
	if err != nil {
		panic(err)
	}

	// Execute template with data and store in a bytes.Buffer for use in email
	var body bytes.Buffer
	err = passwordTemplate.ExecuteTemplate(&body, "layout", &ForgotPasswordData{Link: "https://httpbin.org"})
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