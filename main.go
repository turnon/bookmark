package main

import (
	"flag"
	"fmt"
	"github/turnon/bookmark/bookmark"

	"github.com/k0kubun/pp"
)

func main() {
	const noFile = "no-file"
	filePtr := flag.String("file", noFile, "bookmark file location")
	flag.Parse()

	if *filePtr == noFile {
		panic("no bookmark file given")
	}

	b := bookmark.Load(*filePtr)

	pp.Println("名字重复------------")
	pp.Println(b.DupName())
	pp.Println("URL重复------------")
	pp.Println(b.DupURL())
	pp.Println("网站统计------------")
	pp.Println(b.Hosts())
	pp.Println("目录统计------------")
	pp.Println(b.Folders())

	for _, e := range b.Entries() {
		fmt.Println(e.ToJson())
	}
}
