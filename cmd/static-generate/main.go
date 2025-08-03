package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/gotti/gomd-blog/pkg/framework"
	"github.com/gotti/gomd-blog/pkg/generator"
)

func main() {
	// flag to set output directory
	outputDir := flag.String("output", "./static", "output directory")
	flag.Parse()
	bi, err := generator.NewBlogGenerator("./config.json")
	if err != nil {
		panic(err)
	}
	g, err := framework.NewGenerator([]framework.IndexingAddon{
		bi.NewBlogIndexGenerator(),
	}, nil, []framework.TemplateAddon{
		bi.NewBlogTemplateAddon(),
		bi.NewPageTemplateAddon(),
	})
	if err != nil {
		panic(err)
	}
	// data is a map of path to content

	pages, err := framework.LoadLocalPages()
	if err != nil {
		panic(err)
	}
	data, err := g.Generate(&pages)
	if err != nil {
		panic(err)
	}
	for k, v := range data {
		// create directory if not exists
		os.MkdirAll(filepath.Join(*outputDir, filepath.Dir(k)), 0755)
		// write file
		os.WriteFile(filepath.Join(*outputDir, k), []byte(v), 0644)
	}
}
