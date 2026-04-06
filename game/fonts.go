package game

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Font faces at different sizes for menu/UI rendering.
var (
	FontTitle    *text.GoTextFace // Large title
	FontMenu     *text.GoTextFace // Menu items
	FontMenuSmall *text.GoTextFace // Hints, smaller text
	FontHUD      *text.GoTextFace // In-game HUD
)

func InitFonts(fontData []byte) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	FontTitle = &text.GoTextFace{Source: source, Size: 72}
	FontMenu = &text.GoTextFace{Source: source, Size: 32}
	FontMenuSmall = &text.GoTextFace{Source: source, Size: 20}
	FontHUD = &text.GoTextFace{Source: source, Size: 18}
}
