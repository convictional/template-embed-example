package email

import (
	"io/ioutil"
	"testing"
)

func BenchmarkSendForgotPasswordEmailTemplatesInitFS(b *testing.B) {
	benchSender := Sender{ioutil.Discard}
	for n:=0;n<b.N;n++ {
		_ = benchSender.SendForgotPasswordEmail("test@email.com")
	}
}

func TestSendForgotPasswordEmail(t *testing.T) {
	testSender := Sender{ioutil.Discard}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Default Pass",
			args: args{"test@email.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testSender.SendForgotPasswordEmail(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("SendForgotPasswordEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}