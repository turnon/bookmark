package main

import (
	"flag"

	"github/turnon/bookmark/bookmark"
	"github/turnon/bookmark/views"

	"github.com/gin-gonic/gin"
)

type query struct {
	Stat   string `form:"stat"`
	Name   string `form:"name"`
	URL    string `form:"url"`
	Folder string `form:"folder"`
}

type page struct {
	Stats    *bookmark.Stats
	StatOpts map[string]string
	Query    query
}

func main() {
	const noFile = "no-file"
	filePtr := flag.String("file", noFile, "bookmark file location")
	flag.Parse()

	if *filePtr == noFile {
		panic("no bookmark file given")
	}

	b := bookmark.Load(*filePtr)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		q := query{}
		c.ShouldBindQuery(&q)
		if q.Stat == "" {
			q.Stat = "dirs"
		}

		filter := &bookmark.EntryFilter{Name: q.Name, URL: q.URL, Folder: q.Folder}
		stats, _ := b.Filter(filter).VerboseStat(q.Stat)
		p := page{Stats: stats, StatOpts: bookmark.StatOpts(), Query: q}

		bytes, _ := views.Render("index.html", p)
		c.Data(200, "text/html; charset=utf-8", bytes)
	})
	r.Run()
}
