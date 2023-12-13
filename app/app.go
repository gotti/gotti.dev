package main

import (
	"net/http"

	"github.com/gotti/gomd-blog/pkg/framework"
)

func main() {
	g, err := framework.LoadTemplates()
	if err != nil {
		panic(err)
	}
	data, err := g.Generate()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if path == "" || path[len(path)-1] == '/' {
			path += "index.html"
		}
		// if path ended with /, add index.html
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
