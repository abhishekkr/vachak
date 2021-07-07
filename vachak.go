package main

import (
	"flag"
	"fmt"

	"github.com/abhishekkr/vachak/book"
	"github.com/abhishekkr/vachak/cui"
	"github.com/abhishekkr/vachak/epubreader"
	"github.com/abhishekkr/vachak/tts"
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

func reader(count int, page book.Page) {
	cui.Slides(page)
	if *isTts {
		tts.Espeak(page)
	}
}
