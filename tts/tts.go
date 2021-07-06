package tts

import (
	"strings"

	"github.com/abhishekkr/vachak/book"

	htgotts "github.com/hegedustibor/htgo-tts"
)

func Espeak(page book.Page) {
	for _, w := range strings.Split(page.Text(), "\n") {
		speech := htgotts.Speech{Folder: "/tmp/vachak/audio", Language: "en"}
		speech.Speak(w)

	}
}
