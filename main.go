package main

import (
	"github.com/convictional/template-embed-example/user"
	"os"
)

func main() {
	forgetfulUser := user.User{Email: "test@example.com"}
	forgetfulUser.ResetPassword()
	os.Exit(0)
}
