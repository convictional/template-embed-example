package user

import "github.com/convictional/template-embed-example/email"

type User struct {
	Name string
	Email string
}

func (u *User) ResetPassword() {
	err := email.SendForgotPasswordEmail(u.Email)
	if err != nil {
		panic(err)
	}
}