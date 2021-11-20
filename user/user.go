package user

import (
	"embed"
	"github.com/convictional/template-embed-example/email"
	"html/template"
	"os"
)

const forgotPasswordPath = "templates/forgot_password.html"

var (
	//go:embed templates/forgot_password.html
	forgotPasswordFS embed.FS
	forgotPasswordTemplate template.Template
)

type User struct {
	Email string
}

func (u *User) ResetPassword() {
	err := u.sendPasswordResetEmail()
	if err != nil {
		panic(err)
	}
}

type ForgotPasswordData struct {
	Link string
}

func (u *User) sendPasswordResetEmail() error{
	emailSender := email.Sender{Writer: os.Stdout}
	emailConf := email.SendConfig{
		Subject:      "Reset Password",
		To:           u.Email,
		Template:     forgotPasswordTemplate,
		TemplateData: ForgotPasswordData{Link: "https://httpbin.org"},
	}
	return emailSender.SendEmail(emailConf)
}