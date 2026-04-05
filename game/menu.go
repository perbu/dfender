package game

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const menuItemCount = 3

var menuLabels = [menuItemCount]string{"NEW GAME", "CREDITS", "QUIT"}

func (g *Game) updateMenu() {
	if g.MenuKeyDelay > 0 {
		g.MenuKeyDelay--
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.MenuSelection = (g.MenuSelection + 1) % menuItemCount
		g.MenuKeyDelay = 12
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.MenuSelection = (g.MenuSelection + menuItemCount - 1) % menuItemCount
		g.MenuKeyDelay = 12
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.MenuKeyDelay = 15
		switch g.MenuSelection {
		case 0: // New Game
			g.reset()
		case 1: // Credits
			g.State = StateCredits
			g.MenuKeyDelay = 20
		case 2: // Quit
			os.Exit(0)
		}
	}
}

func (g *Game) updateCredits() {
	if g.MenuKeyDelay > 0 {
		g.MenuKeyDelay--
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.State = StateMenu
		g.MenuKeyDelay = 15
	}
}

func (g *Game) drawMenuScreen(screen *ebiten.Image) {
	s := g.Shaders

	// Starfield background → scene → bloom.
	s.SceneImage.Clear()
	s.DrawStarfield(s.SceneImage, g.Tick)

	// Draw a decorative border frame on the scene.
	borderInset := float32(60)
	bw := float32(ScreenWidth) - borderInset*2
	bh := float32(ScreenHeight) - borderInset*2
	vector.StrokeRect(s.SceneImage, borderInset, borderInset, bw, bh, 2, ColorBorder, false)

	// Corner accents — small diagonal lines at each corner.
	accent := float32(20)
	// Top-left
	vector.StrokeLine(s.SceneImage, borderInset, borderInset+accent, borderInset+accent, borderInset, 2, ColorBorder, false)
	// Top-right
	vector.StrokeLine(s.SceneImage, float32(ScreenWidth)-borderInset-accent, borderInset, float32(ScreenWidth)-borderInset, borderInset+accent, 2, ColorBorder, false)
	// Bottom-left
	vector.StrokeLine(s.SceneImage, borderInset, float32(ScreenHeight)-borderInset-accent, borderInset+accent, float32(ScreenHeight)-borderInset, 2, ColorBorder, false)
	// Bottom-right
	vector.StrokeLine(s.SceneImage, float32(ScreenWidth)-borderInset-accent, float32(ScreenHeight)-borderInset, float32(ScreenWidth)-borderInset, float32(ScreenHeight)-borderInset-accent, 2, ColorBorder, false)

	// Bloom the scene.
	s.ApplyBloom(screen)

	// Title — drawn on top of bloom.
	titleX := ScreenWidth/2 - 60
	titleY := 200
	ebitenutil.DebugPrintAt(screen, "d F E N D E R", titleX, titleY)

	// Subtitle.
	ebitenutil.DebugPrintAt(screen, "A R E N A   S H O O T E R", titleX-50, titleY+30)

	if g.State == StateCredits {
		g.drawCredits(screen)
		return
	}

	// Menu items.
	menuStartY := ScreenHeight/2 - 20
	for i, label := range menuLabels {
		y := menuStartY + i*40
		x := ScreenWidth/2 - 50

		if i == g.MenuSelection {
			// Selection indicator — gold chevron.
			indicator := fmt.Sprintf("> %s <", label)
			ebitenutil.DebugPrintAt(screen, indicator, x-20, y)
		} else {
			ebitenutil.DebugPrintAt(screen, label, x, y)
		}
	}

	// Nav hint.
	ebitenutil.DebugPrintAt(screen, "W/S or ARROWS to navigate  -  ENTER to select", ScreenWidth/2-160, ScreenHeight-120)
}

func (g *Game) drawCredits(screen *ebiten.Image) {
	cx := ScreenWidth/2 - 100
	cy := ScreenHeight/2 - 80

	ebitenutil.DebugPrintAt(screen, "C R E D I T S", cx+30, cy)
	cy += 50
	ebitenutil.DebugPrintAt(screen, "Game Design & Programming", cx, cy)
	cy += 20
	ebitenutil.DebugPrintAt(screen, "Per Buer", cx+30, cy)

	cy += 40
	ebitenutil.DebugPrintAt(screen, "AI Co-Pilot & Code Generation", cx-15, cy)
	cy += 20
	ebitenutil.DebugPrintAt(screen, "Claude (Anthropic)", cx+15, cy)

	cy += 40
	ebitenutil.DebugPrintAt(screen, "Built with Ebitengine", cx+5, cy)

	cy += 60
	ebitenutil.DebugPrintAt(screen, "PRESS ENTER OR ESC TO RETURN", cx-20, cy)
}
