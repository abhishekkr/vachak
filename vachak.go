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
	isCui        = flag.Bool("cui", true, "for Markdown pretty print per page")
	isTts        = flag.Bool("tts", false, "for Text to Speech per page")
	epubFilepath = flag.String("epub", "", "epub file path for markdown show or read aloud")
)

func main() {
	flag.Parse()
	if *epubFilepath != "" {
		epubReader()
	} else {
		fmt.Println("run with `-h` to see usage")
	}
}

func epubReader() {
	book := epubreader.Epub{Filepath: *epubFilepath}
	book.Open()
	book.Metadata()
	fmt.Printf("Rootfile Path:%v\n\n", book.Book.Container.Rootfile.Path)
	for idx, point := range book.Book.Ncx.Points {
		fmt.Printf("Text: %v\n", point.Text)
		fmt.Printf("Content Src: %v\n", point.Content.Src)

		page := epubreader.Read(&book, point.Content.Src)
		reader(idx, page)
	}
}

func reader(count int, page book.Page) {
	if *isCui {
		cui.Slides(page)
	}
	if *isTts {
		tts.Espeak(page)
	}
}
