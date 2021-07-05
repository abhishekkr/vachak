package book

type Page interface {
	Markdown() string
	Text() string
}
