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

	pp.Println(b.Stats())

	for _, e := range b.Entries() {
		fmt.Println(e.ToJson())
	}
}
