package main

import (
	"flag"

	"github/turnon/bookmark/bookmark"
	"github/turnon/bookmark/views"

	"github.com/gin-gonic/gin"
)

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
		bytes, _ := views.Render("index.html", b)
		c.Data(200, "text/html; charset=utf-8", bytes)
	})
	r.Run()
}
