package framework

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gotti/gomd-blog/pkg/mdparser"
)

//go:embed static/*
var templates embed.FS

// Config has the configuration for the framework
type Config struct {
	Layout          Layout     `json:"layout"`
	BlogPath        string     `json:"blog_path"`
	DefaultMetaData MetaData   `json:"default_metadata"`
	Menu            []MenuItem `json:"menu"`
}

// MetaData has the configuration for the metadata
type MetaData struct {
	Title     *string `json:"title"`
	TwitterID *string `json:"twitter_id"`
	SiteName  *string `json:"site_name"`
	Image     *string `json:"image"`
}

// WithDefault returns the default value for the metadata
func (m *MetaData) WithDefault(title *string, image *string) *MetaData {
	if title != nil {
		tmp := *title
		m.Title = &tmp
	}
	if image != nil {
		tmp := *image
		m.Image = &tmp
	}
	return m
}

// Layout has the configuration for the layout
type Layout struct {
	Name string `json:"name"`
}

// MenuItem has the configuration for the menu
type MenuItem struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// Parse parses the configuration
func Parse(config []byte) (*Config, error) {
	var c Config
	err := json.Unmarshal(config, &c)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return &c, nil
}

// Generator has the templates, and generates the pages
type Generator struct {
	templates *template.Template
	config    *Config
}

// LoadTemplates loads the templates
func LoadTemplates() (*Generator, error) {
	t, err := template.ParseFS(templates, "static/*.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	config, err := Parse(configFile)
	return &Generator{templates: t, config: config}, nil
}

// Generate generates the pages
func (g *Generator) Generate() (map[string]string, error) {
	pages, err := g.load()
	if err != nil {
		return nil, fmt.Errorf("error loading pages: %w", err)
	}
	data := make(map[string]string)
	for k, v := range pages {
		p := g.generateLayout(v)
		data[k] = string(p)
	}
	data["post/index.html"] = string(g.generateBlogIndex(pages))
	return data, nil
}

func (g *Generator) load() (map[string]*mdparser.Root, error) {
	pages := make(map[string]*mdparser.Root)
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
		md, err := mdparser.Parse(string(lines))
		if err != nil {
			fmt.Printf("error parsing markdown, file=%v: %v\n", path, err)
			return nil
		}

		//remove pages/ from path
		path = path[len("pages/"):]

		//change the extension to html, if is not index.md
		htmlPath := path[:len(path)-len(filepath.Ext(path))] + ".html"
		if filepath.Base(path) != "index.md" {
			htmlPath = path[:len(path)-len(filepath.Ext(path))] + "/index.html"
		}
		pages[htmlPath] = md
		return nil
	})
	return pages, nil
}

// Posts have the posts
type Posts struct {
	Posts []Post
}

// Post has the post
type Post struct {
	Title string
	Link  string
	Date  time.Time
	Tags  []string
}

func (g *Generator) generateBlogIndex(pages map[string]*mdparser.Root) template.HTML {
	blogs := map[string]*mdparser.Root{}
	for k, v := range pages {
		if strings.HasPrefix(k, g.config.BlogPath) {
			blogs[k] = v
		}
	}

	fmt.Printf("blogs=%+v\n", blogs)

	posts := Posts{}

	for k, v := range blogs {
		t := v.MetaData.Title
		title := ""
		if t == nil {
			title = filepath.Base(k)
		} else {
			title = *t
		}
		d := v.MetaData.Date
		var date time.Time
		if d == nil {
			date = time.Unix(0, 0)
		} else {
			date = *d
		}
		posts.Posts = append(posts.Posts, Post{
			Title: title,
			Link:  "/" + k,
			Date:  date,
			Tags:  v.MetaData.Tags,
		})
	}

	sort.Slice(posts.Posts, func(i, j int) bool {
		return posts.Posts[i].Date.After(posts.Posts[j].Date)
	})

	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "post_index", posts)
	return withDiv("content", g.generateMenu()+template.HTML(buf.String()))

}

func (g *Generator) generateHeader(meta *MetaData) template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "head", meta)
	return template.HTML(buf.String())
}

func (g *Generator) generateTitle(title *string) template.HTML {
	if title == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "title", title)
	return template.HTML(buf.String())
}

func (g *Generator) generatePost(title *string, child template.HTML) template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "post", child)
	return withDiv("content", g.generateTitle(title)+template.HTML(buf.String()))
}

func (g *Generator) generatePage() template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "page", g.config.Layout)
	return withDiv("content", template.HTML(buf.String()))
}

func (g *Generator) generateMenu() template.HTML {
	buf := new(bytes.Buffer)
	g.templates.ExecuteTemplate(buf, "menu", g.config.Menu)
	return withDiv("menu", template.HTML(buf.String()))
}

func withDiv(id string, s template.HTML) template.HTML {
	return template.HTML(fmt.Sprintf("<div id=%s>%s</div>", id, s))
}

func (g *Generator) generateLayout(md *mdparser.Root) template.HTML {
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
		Body: g.generateMenu() + g.generatePost(md.MetaData.Title, template.HTML(md.Objects.ToHTML())),
	})
	return template.HTML(buf.String())
}
