package cui

import (
	"fmt"

	"github.com/abhishekkr/hanasuhon/book"

	"github.com/charmbracelet/glamour"
)

func Slides(page book.Page) error {
	out, err := glamour.Render(page.Markdown(), "dark")
	fmt.Print(out)
	return err
}
