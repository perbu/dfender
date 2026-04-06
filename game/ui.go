package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawUI(screen *ebiten.Image, g *Game) {
	// Score — top right.
	scoreStr := fmt.Sprintf("%d", g.Score.Score)
	drawTextAt(screen, scoreStr, FontHUD, float64(ScreenWidth-150), 10, ColorUI)

	// Combo.
	if g.Score.Combo > 1 {
		comboStr := fmt.Sprintf("x%d", g.Score.Combo)
		drawTextAt(screen, comboStr, FontHUD, float64(ScreenWidth-150), 34, ColorBorder)
	}

	// Wave — center screen announcement.
	if g.State == StateWaveIntro {
		waveStr := fmt.Sprintf("WAVE %d", g.Wave.Number)
		drawTextCentered(screen, waveStr, FontMenu, float64(ScreenHeight)/2-40, ColorBorder)
	}

	// Lives — top right, below combo.
	livesStr := fmt.Sprintf("LIVES: %d", g.Lives)
	drawTextAt(screen, livesStr, FontHUD, float64(ScreenWidth-150), 56, ColorUI)

	// Heat bar — below lives.
	drawHeatBar(screen, g)

	// Game over.
	if g.State == StateGameOver {
		cy := float64(ScreenHeight)/2 - 80
		drawTextCentered(screen, "GAME OVER", FontTitle, cy, color.RGBA{0xFF, 0x33, 0x33, 0xFF})
		cy += 90
		drawTextCentered(screen, fmt.Sprintf("SCORE: %d    WAVE: %d", g.Score.Score, g.Wave.Number), FontMenu, cy, ColorUI)
		cy += 50
		if g.HighScores.Qualifies(g.Score.Score) {
			drawTextCentered(screen, "NEW HIGH SCORE!", FontMenu, cy, ColorBorder)
			cy += 50
		}
		drawTextCentered(screen, "PRESS ENTER TO CONTINUE", FontMenuSmall, cy, ColorBorderDim)
	}

	// FPS (keep debug font — it's a debug readout).
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()), 10, 10)
}

func drawHeatBar(screen *ebiten.Image, g *Game) {
	barW := float32(150)
	barH := float32(8)
	barX := float32(ScreenWidth) - barW - 20
	barY := float32(78)

	// Background.
	vector.DrawFilledRect(screen, barX, barY, barW, barH, color.RGBA{0x1A, 0x1A, 0x2E, 0xFF}, false)

	// Fill.
	heat := float32(g.Turret.Heat)
	fillColor := lerpColor(ColorHeatCool, ColorHeatHot, heat)
	vector.DrawFilledRect(screen, barX, barY, barW*heat, barH, fillColor, false)

	// Border.
	vector.StrokeRect(screen, barX, barY, barW, barH, 1, ColorBorderDim, false)

	// Label if overheated.
	if g.Turret.Cooldown > 0 {
		drawTextAt(screen, "OVERHEAT", FontHUD, float64(barX), float64(barY)-20, ColorHeatHot)
	}
}
