package user

import "github.com/convictional/template-embed-example/email"

func init() {
	forgotPasswordTemplate = email.MustParseContentFS(forgotPasswordFS, forgotPasswordPath)
}
