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

func EpubReader(epubFilepath string, reader func(int, book.Page)) {
	fyl := Epub{Filepath: epubFilepath}
	fyl.open()
	defer fyl.close()
	fyl.metadata()
	fmt.Printf("Rootfile Path:%v\n\n", fyl.Book.Container.Rootfile.Path)
	fyl.navPointsReader(fyl.Book.Ncx.Points, reader)
}

func (e *Epub) navPointsReader(points []epub.NavPoint, reader func(int, book.Page)) {
	for idx, point := range points {
		fmt.Printf("Text: %v\n", point.Text)
		fmt.Printf("Content Src: %v\n", point.Content.Src)

		page := e.pageReader(point.Content.Src)
		reader(idx, page)
		e.navPointsReader(point.Points, reader)
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
