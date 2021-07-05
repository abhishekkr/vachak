package epubreader

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"

	goquery "github.com/PuerkitoBio/goquery"
	epub "github.com/kapmahc/epub"
	html2text "jaytaylor.com/html2text"
)

type Epub struct {
	Book     *epub.Book
	Filepath string
}

type Page struct {
	Epub *Epub
	Path string
	XML  struct {
		XMLName xml.Name `xml:"html"`
		Title   string   `xml:"head>title"`
		Body    struct {
			InnerXml string `xml:",innerxml"`
		} `xml:"body"`
	}
}

func Read(e *Epub, xmlpath string) *Page {
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

func (p *Page) read() error {
	fd, err := p.open()
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	b, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	return xml.Unmarshal(b, &p.XML)
}

func (p *Page) open() (io.ReadCloser, error) {
	xmlpath := p.filename()
	fd, err := zip.OpenReader(p.Epub.Filepath)
	if err != nil {
		panic(err)
	}
	for _, f := range fd.File {
		if f.Name == xmlpath {
			return f.Open()
		}
	}
	return nil, fmt.Errorf("file %s not exist", xmlpath)
}

func (p *Page) filename() string {
	return path.Join(path.Dir(p.Epub.Book.Container.Rootfile.Path), p.Path)
}

func (p *Page) Text() string {
	reader := strings.NewReader(p.XML.Body.InnerXml)
	doc, _ := goquery.NewDocumentFromReader(reader)

	doc.Find(".").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})
	return doc.Text()
}

func (p *Page) Markdown() string {
	text, err := html2text.FromString(p.XML.Body.InnerXml, html2text.Options{PrettyTables: true})
	if err != nil {
		panic(err)
	}
	return text
}
