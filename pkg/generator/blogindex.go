package generator

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gotti/gomd-blog/pkg/framework"
	"github.com/gotti/gomd-blog/pkg/mdparser"
)

type blogIndexGenerator struct {
	blogPath string
}

func (g *blogIndexGenerator) GenerateIndexes(pages *framework.Pages) (*framework.Pages, error) {
	blogs := framework.NewPages()
	for k, v := range *pages {
		if strings.HasPrefix(k, g.blogPath) {
			blogs[k] = v
		}
	}

	fmt.Printf("blogs=%+v\n", blogs)

	posts := Posts{}

	for k, v := range blogs {
		t := v.Contents.MetaData.Title
		title := ""
		if t == nil {
			title = filepath.Base(k)
		} else {
			title = *t
		}
		d := v.Contents.MetaData.Date
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
			Tags:  v.Contents.MetaData.Tags,
		})
	}

	sort.Slice(posts.Posts, func(i, j int) bool {
		return posts.Posts[i].Date.After(posts.Posts[j].Date)
	})

	index := mdparser.Root{
		Objects: mdparser.Objects{
			mdparser.Heading{Level: 1, PlainObjectImpl: mdparser.PlainObjectImpl{
				InlineContainerObjectImpl: mdparser.InlineContainerObjectImpl{
					Contents: mdparser.InlineBlocks{
						Children: []mdparser.InlineBlock{
							mdparser.InlineText{Text: []rune("Blog")},
						},
					},
				},
			}},
		},
	}

	for _, p := range posts.Posts {
		index.Objects = append(index.Objects, mdparser.Divider{
			Objects: mdparser.Objects{
				mdparser.InlineLink{
					URL: []rune(p.Link),
					Text: mdparser.InlineBlocks{
						Children: []mdparser.InlineBlock{
							mdparser.InlineText{
								Text: []rune(p.Title),
							},
						},
					},
				},
				mdparser.InlineText{
					Text: []rune(" " + p.Date.Format("2006-01-02")),
				},
			},
		})
	}
	pages.AddPage("post/index.html", &framework.Page{
		Contents: &index,
		Filename: "index.html",
		Path:     "post/index.html",
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
