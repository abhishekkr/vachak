package epubreader

import (
	"fmt"
	"strings"

	"github.com/abhishekkr/vachak/book"

	epub "github.com/kapmahc/epub"
)

type Epub struct {
	Book     *epub.Book
	Filepath string
}

func EpubReader(epubFilepath string, reader func(book.Page)) {
	fyl := Epub{Filepath: epubFilepath}
	fyl.open()
	defer fyl.close()
	fyl.metadata()
	fmt.Printf("Rootfile Path:%v\n\n", fyl.Book.Container.Rootfile.Path)

	pageQueue := fyl.navPointsLoader([]string{}, fyl.Book.Ncx.Points)
	fyl.navPointsReader(pageQueue, reader)
}

func (e *Epub) navPointsLoader(pages []string, points []epub.NavPoint) []string {
	for _, point := range points {
		xmlfile := strings.Split(point.Content.Src, "#")[0]
		if pageQueued(pages, xmlfile) {
			continue
		}
		pages = append(pages, xmlfile)
		pages = e.navPointsLoader(pages, point.Points)
	}
	return pages
}

func pageQueued(pages []string, xmlfile string) bool {
	for _, pagePath := range pages {
		if xmlfile == pagePath {
			return true
		}
	}
	return false
}

func (e *Epub) navPointsReader(pageQueue []string, reader func(book.Page)) {
	for _, xmlfile := range pageQueue {
		fmt.Printf("Content Src: %v\n", xmlfile)

		page := e.pageReader(xmlfile)
		reader(page)
	}
}

func (e *Epub) pageReader(xmlpath string) *Page {
	page := &Page{Epub: e, Path: xmlpath}
	err := page.read()
	if err != nil {
		panic(err)
	}
	return page //page.XML.Title, bodyText(page.XML.Body.InnerXml)
}

func (e *Epub) open() {
	var err error
	e.Book, err = epub.Open(e.Filepath)
	if err != nil {
		panic(err)
	}
}

func (e *Epub) close() {
	e.Book.Close()
}

func (e *Epub) metadata() {
	fmt.Printf("title: %v\ncreator(s): %v\ndescription: %v\n",
		strings.Join(e.Book.Opf.Metadata.Title, ", "),
		strings.Join(e.creators(), ", "),
		strings.Join(e.Book.Opf.Metadata.Description, ", "))
}

func (e *Epub) creators() []string {
	creatorList := make([]string, len(e.Book.Opf.Metadata.Creator))
	for idx, creator := range e.Book.Opf.Metadata.Creator {
		creatorList[idx] = creator.Data
	}
	return creatorList
}
