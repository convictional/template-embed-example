package email

import "html/template"

func init() {
	baseTemplate = template.Must(template.New("layout").ParseFS(baseLayoutFS, templateLayout))
}
