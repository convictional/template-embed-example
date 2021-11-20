package user

import (
	"github.com/convictional/template-embed-example/email"
	"os"
)

type User struct {
	Email string
}

func (u *User) ResetPassword() {
	emailSender := email.Sender{Writer: os.Stdout}
	err := emailSender.SendForgotPasswordEmailInit(u.Email)
	if err != nil {
		panic(err)
	}
}