package tts

import (
	"strings"

	"github.com/abhishekkr/hanasuhon/book"

	htgotts "github.com/hegedustibor/htgo-tts"
)

func Espeak(page book.Page) {
	for _, w := range strings.Split(page.Text(), "\n") {
		speech := htgotts.Speech{Folder: "/tmp/hanasuhon/audio", Language: "en"}
		speech.Speak(w)

	}
}
