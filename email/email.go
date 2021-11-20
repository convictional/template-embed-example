package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
)

const templateLayout = "templates/layout.html"

var (
	//go:embed templates/layout.html
	baseLayoutFS embed.FS
	baseTemplate *template.Template
)


func MustParseContentFS(fsys fs.FS, patterns ...string) template.Template {
	t, err := baseTemplate.Clone()
	if err != nil {
		panic(err)
	}
	return *template.Must(t.ParseFS(fsys, patterns...))
}

type Sender struct {
	Writer io.Writer
}

type SendConfig struct {
	Subject string
	To string
	Template template.Template
	TemplateData interface{}
	TemplateName string // Optional. Override the template to be executed with the name of a sub-template.
}

func (s Sender) SendEmail(cfg SendConfig) error {
	var buf bytes.Buffer
	var err error

	if cfg.TemplateName == "" {
		err = cfg.Template.Execute(&buf, cfg.TemplateData)
	} else {
		err = cfg.Template.ExecuteTemplate(&buf, cfg.TemplateName, cfg.TemplateData)
	}
	email := fmt.Sprintf("Receipient: %s\nSubject:%s\nBody:\n%s", cfg.To, cfg.Subject, buf.String())
	_, _ = s.Writer.Write([]byte(email))
	return err
}