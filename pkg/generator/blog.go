package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/gotti/gomd-blog/pkg/framework"
	"github.com/gotti/gomd-blog/pkg/mdparser"
)

type generator struct {
	config    *Config
	templates *template.Template
}

type postTemplateAddon struct {
	generator
}

func NewTemplateAddon() *postTemplateAddon {
	return &postTemplateAddon{}
}

func (g *generator) generateHeader(meta *MetaData) template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "head", meta)
	return template.HTML(buf.String())
}

func (g *generator) generateTitle(title *string) template.HTML {
	if title == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "title", title)
	return template.HTML(buf.String())
}

func (g *generator) generatePost(title *string, child template.HTML) template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "post", child)
	return withDiv("content", g.generateTitle(title)+template.HTML(buf.String()))
}

func (g *generator) generatePage() template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "page", g.config.Layout)
	return withDiv("content", template.HTML(buf.String()))
}

func (g *generator) generateMenu(page *framework.Page) template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "menu", g.config.Menu)
	if page != nil {
		g.templates.ExecuteTemplate(buf, "request_changes", fmt.Sprintf("https://github.com/gotti/gotti.dev/blob/main/%s", page.OriginalPath))
	}
	return withDiv("menu", template.HTML(buf.String()))
}

func withDiv(id string, s template.HTML) template.HTML {
	return template.HTML(fmt.Sprintf("<div id=%s>%s</div>", id, s))
}

func (g *generator) generateLayout(md *mdparser.Root) template.HTML {
	title := g.config.DefaultMetaData.SiteName
	if md.MetaData.Title != nil {
		title = md.MetaData.Title
	}
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "layout", struct {
		Head template.HTML
		Body template.HTML
	}{
		Head: g.generateHeader(g.config.DefaultMetaData.WithDefault(title, md.MetaData.Thumbnail)),
		Body: g.generateMenu(nil) + g.generatePost(md.MetaData.Title, template.HTML(md.Objects.ToHTML())),
	})
	return template.HTML(buf.String())
}

func (t *postTemplateAddon) GeneratePage(page *framework.Page, pagehtml template.HTML) (template.HTML, error) {
	if !strings.HasPrefix(page.Path, t.config.BlogPath) {
		return "", fmt.Errorf("path is not a blog path: %v", page.Path)
	}
	meta := page.Contents.MetaData
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
