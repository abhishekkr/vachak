package cui

import (
	"strings"

	"github.com/abhishekkr/vachak/book"

	glamour "github.com/charmbracelet/glamour"
)

func Slides(page book.Page) error {
	md, err := glamour.Render(page.Markdown(), "dark")
	if err != nil {
		return err
	}
	if strings.TrimSpace(md) != "" {
		render(page, md)
	}
	return nil
}
