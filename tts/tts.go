package tts

import (
	"log"
	"strings"

	"vachak/book"
	//htgotts "github.com/hegedustibor/htgo-tts"
)

func Espeak(page book.Page) {
	for _, w := range strings.Split(page.Text(), "\n") {
		//speech := htgotts.Speech{Folder: "/tmp/vachak/audio", Language: "en"}
		//speech.Speak(w)
		_ = w
	}
	log.Println("[WIP] Kokoro or something. !Espeak")
}
