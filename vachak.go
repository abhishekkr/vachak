package main

import (
	"flag"
	"fmt"

	"vachak/book"
	"vachak/cui"
	"vachak/epubreader"
	"vachak/tts"
)

var (
	isTts        = flag.Bool("tts", false, "for Text to Speech per page")
	epubFilepath = flag.String("epub", "", "epub file path for markdown show or read aloud")
)

func main() {
	flag.Parse()
	if *epubFilepath != "" {
		epubreader.EpubReader(*epubFilepath, reader)
	} else {
		fmt.Println("run with `-h` to see usage")
	}
}

func reader(page book.Page) {
	cui.Slides(page)
	if *isTts {
		tts.Espeak(page)
	}
}
