package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	EnemyRadius   = 22.0
	EnemyBaseSpeed = 1.5
)

type Enemy struct {
	X, Y   float64
	VX, VY float64
	Speed  float64
	HP     int
	MaxHP  int
	Alive  bool
	FlashFrames int
}

func updateEnemies(g *Game) {
	for i := range g.Enemies {
		e := &g.Enemies[i]
		if !e.Alive {
			continue
		}
		// Home toward player.
		dx := g.Player.X - e.X
		dy := g.Player.Y - e.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 1 {
			e.VX = dx / dist * e.Speed
			e.VY = dy / dist * e.Speed
		}
		e.X += e.VX
		e.Y += e.VY

		if e.FlashFrames > 0 {
			e.FlashFrames--
		}
	}
	// Compact.
	n := 0
	for i := range g.Enemies {
		if g.Enemies[i].Alive {
			g.Enemies[n] = g.Enemies[i]
			n++
		}
	}
	g.Enemies = g.Enemies[:n]
}

func drawEnemies(screen *ebiten.Image, g *Game, ox, oy float64) {
	for i := range g.Enemies {
		e := &g.Enemies[i]
		cx := float32(e.X + ox)
		cy := float32(e.Y + oy)

		col := ColorEnemy
		if e.FlashFrames > 0 {
			col = ColorUI // white flash
		}

		// Outer glow circle.
		r := float32(EnemyRadius)
		vector.StrokeCircle(screen, cx, cy, r+3, 4, col, false)

		angle := math.Atan2(e.VY, e.VX)
		drawPolygon(screen, cx, cy, r, 3, angle, 4, col)

		// HP indicator for tough enemies.
		if e.MaxHP > 1 {
			hpFrac := float32(e.HP) / float32(e.MaxHP)
			barW := r * 2
			barH := float32(3)
			barX := cx - r
			barY := cy - r - 6
			vector.DrawFilledRect(screen, barX, barY, barW*hpFrac, barH, col, false)
			vector.StrokeRect(screen, barX, barY, barW, barH, 1, ColorBorderDim, false)
		}
	}
}
