package framework

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gotti/gomd-blog/pkg/mdparser"
)

type IndexingAddon interface {
	GenerateIndexes(pages *Pages) (*Pages, error)
}

type ReplaceAddon interface {
	ReplacePages(page *Pages) (*Pages, error)
}

type TemplateAddon interface {
	GeneratePage(page *Page, pagehtml template.HTML) (template.HTML, error)
}

// Generator has the templates, and generates the pages
type Generator struct {
	templates     *template.Template
	indexingAddon []IndexingAddon
	replaceAddon  []ReplaceAddon
	templateAddon []TemplateAddon
}

// NewGenerator loads the templates
func NewGenerator(iadd []IndexingAddon, radd []ReplaceAddon, tadd []TemplateAddon) (*Generator, error) {
	if iadd == nil {
		iadd = []IndexingAddon{}
	}
	if radd == nil {
		radd = []ReplaceAddon{}
	}
	if tadd == nil {
		tadd = []TemplateAddon{}
	}
	return &Generator{indexingAddon: iadd, replaceAddon: radd, templateAddon: tadd}, nil
}

// Generate generates the pages
func (g *Generator) Generate() (map[string]string, error) {
	ret := make(map[string]string)
	pages, err := Load()
	if err != nil {
		return nil, fmt.Errorf("error loading pages: %w", err)
	}
	for _, i := range g.indexingAddon {
		npages, err := i.GenerateIndexes(&pages)
		if err != nil {
			return nil, fmt.Errorf("error generating indexes: %w", err)
		}
		pages = *npages
	}
	for _, r := range g.replaceAddon {
		npages, err := r.ReplacePages(&pages)
		if err != nil {
			return nil, fmt.Errorf("error replacing pages: %w", err)
		}
		pages = *npages
	}
	for _, p := range pages {
		html := p.Contents.ToHTML()
		for _, t := range g.templateAddon {
			npage, err := t.GeneratePage(p, template.HTML(html))
			if err != nil {
				fmt.Printf("error generating pages, skipping...: %v\n", err)
				continue
			}
			html = string(npage)
			break
		}
		ret[p.Path] = html
	}
	return ret, nil
}

type Pages map[string]*Page

func NewPages() Pages {
	return make(map[string]*Page)
}

func (p *Pages) AddPage(path string, page *Page) {
	if page.Path != path {
		fmt.Printf("filename and page.filename are different!!!, path=%v, page.path=%v\n", path, page.Path)
	}
	(*p)[path] = page
}

type Page struct {
	Contents *mdparser.Root
	Filename string
	Path     string
}

func Load() (Pages, error) {
	parser := mdparser.NewLineParser(nil)
	pages := NewPages()
	filepath.WalkDir("pages", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}
		if d.IsDir() {
			return err
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		lines, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("error reading file, file=%v; %v\n", path, err)
			return nil
		}
		md, err := parser.Parse(string(lines))
		if err != nil {
			fmt.Printf("error parsing markdown, file=%v: %v\n", path, err)
			return nil
		}

		//remove pages/ from path
		path = path[len("pages/"):]

		//remove the extension from the path
		path = path[:len(path)-len(filepath.Ext(path))]

		if filepath.Base(path) != "index" {
			path = filepath.Join(filepath.Dir(path), path, "index")
		}

		//add .html to the path
		path += ".html"

		//get filename
		filename := filepath.Base(path)


		_, ok := pages[path]
		if ok {
			fmt.Printf("duplicated filename!!!, filename=%v\n", filename)
		}

		p := &Page{
			Contents: md,
			Filename: filename,
			Path:     path,
		}
		pages.AddPage(path, p)

		return nil
	})
	return pages, nil
}
