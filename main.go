package main

import (
	"fmt"
	"os"

	"github.com/gotti/gomd-blog/pkg/addon/obsidian"
	"github.com/gotti/gomd-blog/pkg/mdparser"
)

func main() {
	md, err := os.ReadFile("./pages/about.md")
	if err != nil {
		panic(err)
	}
	ob := obsidian.NewObsidianAddon()
	parser := mdparser.NewLineParser([]mdparser.ParserAddon{ob})
	p, err := parser.Parse(string(md))
	if err != nil {
		panic(err)
	}
  fmt.Println(p.ToHTML())
  for _, tag := range ob.GetTags() {
    fmt.Println(string(tag))
  }
}
