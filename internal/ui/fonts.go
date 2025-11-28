package ui

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/default.ttf
var defaultFontData []byte

func getDefaultFontSource() *text.GoTextFaceSource {
	if defaultFontSource == nil {
		source, err := text.NewGoTextFaceSource(bytes.NewReader(defaultFontData))
		if err != nil {
			panic(err)
		}
		defaultFontSource = source
	}
	return defaultFontSource
}