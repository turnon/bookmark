package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"github/turnon/bookmark/bookmark"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

//go:embed views
var views embed.FS

func main() {
	const noFile = "no-file"
	filePtr := flag.String("file", noFile, "bookmark file location")
	flag.Parse()

	if *filePtr == noFile {
		panic("no bookmark file given")
	}

	b := bookmark.Load(*filePtr)

	pp.Println(b.Stats())

	for _, e := range b.Entries() {
		fmt.Println(e.ToJson())
	}

	t, _ := template.ParseFS(views, "views/index.html")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		t.Execute(&buffer, b)
		c.Data(200, "text/html; charset=utf-8", buffer.Bytes())
	})
	r.Run()
}
