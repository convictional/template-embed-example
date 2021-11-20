package user

import (
	"testing"
)

func BenchmarkSendForgotPasswordEmailTemplatesInitFS(b *testing.B) {
	testUser := User{Email: "test@example.com"}
	for n:=0;n<b.N;n++ {
		_ = testUser.sendPasswordResetEmail()
	}
}