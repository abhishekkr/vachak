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
	fyl.Open()
	fyl.Metadata()
	fmt.Printf("Rootfile Path:%v\n\n", fyl.Book.Container.Rootfile.Path)
	for idx, point := range fyl.Book.Ncx.Points {
		fmt.Printf("Text: %v\n", point.Text)
		fmt.Printf("Content Src: %v\n", point.Content.Src)

		page := fyl.ReadPage(point.Content.Src)
		reader(idx, page)
	}
}

func (e *Epub) ReadPage(xmlpath string) *Page {
	page := &Page{Epub: e, Path: xmlpath}
	err := page.read()
	if err != nil {
		panic(err)
	}
	return page //page.XML.Title, bodyText(page.XML.Body.InnerXml)
}

func (e *Epub) Open() {
	var err error
	e.Book, err = epub.Open(e.Filepath)
	if err != nil {
		panic(err)
	}
}

func (e *Epub) Close() {
	e.Book.Close()
}

func (e *Epub) Metadata() {
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
