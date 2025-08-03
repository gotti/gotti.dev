package main

import (
	"net/http"

	"github.com/gotti/gomd-blog/pkg/framework"
	"github.com/gotti/gomd-blog/pkg/generator"
)

func main() {
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
	pages, err := framework.LoadLocalPages()
	if err != nil {
		panic(err)
	}
	data, err := g.Generate(&pages)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		// if path ended with /, add index.html
		if path == "" || path[len(path)-1] == '/' {
			path += "index.html"
		}
		d, ok := data[path]
		if !ok {
			h := http.FileServer(http.Dir("./pages"))
			h.ServeHTTP(w, r)
			return
		}
		w.Write([]byte(d))
	})
	http.ListenAndServe(":8080", nil)
}
