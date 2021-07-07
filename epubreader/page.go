package epubreader

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"

	goquery "github.com/PuerkitoBio/goquery"
	html2text "jaytaylor.com/html2text"
)

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

func (p *Page) BookName() string {
	return strings.Join(p.Epub.Book.Opf.Metadata.Title, ", ")
}

func (p *Page) Creators() string {
	return strings.Join(p.Epub.creators(), ", ")
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
