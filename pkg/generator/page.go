package generator

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gotti/gomd-blog/pkg/framework"
)

type standalonePageAddon struct {
	generator
}

func (t *standalonePageAddon) GeneratePage(page *framework.Page, pagehtml template.HTML) (template.HTML, error) {
	meta := page.Contents.MetaData
	if page.Path == "index.html" {
		fmt.Println(meta)
	}
	buf := new(bytes.Buffer)
	t.templates.ExecuteTemplate(buf, "layout", struct {
		Head template.HTML
		Body template.HTML
	}{
		Head: t.generateHeader(t.config.DefaultMetaData.WithDefault(meta.Title, meta.Thumbnail)),
		Body: t.generateMenu(page) + t.generatePost(meta.Title, template.HTML(pagehtml)),
	})
	return template.HTML(buf.String()), nil
}
