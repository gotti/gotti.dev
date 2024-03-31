package generator

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

// Layout has the configuration for the layout
type Layout struct {
	Name string `json:"name"`
}

// MenuItem has the configuration for the menu
type MenuItem struct {
	Name string `json:"name"`
	Link string `json:"link"`
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
	newMetaData := MetaData{
		Title:     m.Title,
		TwitterID: m.TwitterID,
		SiteName:  m.SiteName,
		Image:     m.Image,
	}
	if title != nil {
		newMetaData.Title = title
	}
	if image != nil {
		newMetaData.Image = image
	}
	return &newMetaData
}

// Config has the configuration for the framework
type Config struct {
	Layout          Layout     `json:"layout"`
	BlogPath        string     `json:"blog_path"`
	DefaultMetaData MetaData   `json:"default_metadata"`
	Menu            []MenuItem `json:"menu"`
}

// ParseConfig parses the configuration
func ParseConfig(config []byte) (*Config, error) {
	var c Config
	err := json.Unmarshal(config, &c)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return &c, nil
}

type BlogGenerator struct {
	config    *Config
	templates *template.Template
}

func (b BlogGenerator) NewBlogIndexGenerator() *blogIndexGenerator {
	return &blogIndexGenerator{blogPath: b.config.BlogPath}
}

func (b BlogGenerator) NewBlogTemplateAddon() *postTemplateAddon {
	return &postTemplateAddon{generator: generator{config: b.config, templates: b.templates}}
}

func (b BlogGenerator) NewPageTemplateAddon() *standalonePageAddon {
	return &standalonePageAddon{generator: generator{config: b.config, templates: b.templates}}
}

//go:embed static/*
var tmps embed.FS

func NewBlogGenerator(configPath string) (*BlogGenerator, error) {
	config, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}
	c, err := ParseConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	t, err := template.ParseFS(tmps, "static/*.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}
	return &BlogGenerator{config: c, templates: t}, nil
}
