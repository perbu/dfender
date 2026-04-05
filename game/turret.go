package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TurretRotSpeed   = math.Pi // radians per second → 180 deg/s
	TurretLength     = 38.0
	FireRate         = 6  // frames between shots (10/sec at 60fps)
	ProjectileSpeed  = 12.0
	HeatPerShot      = 0.05  // ~2 sec to overheat at 10 rps
	HeatDecay        = 0.004 // ~4 sec to cool from full
	CooldownTime     = 120   // frames (2 sec)
)

type Turret struct {
	Angle     float64 // radians, 0 = up
	Heat      float64 // 0..1
	Cooldown  int     // frames remaining in overheat lockout
	FireTimer int     // frames until next shot allowed
}

func NewTurret() Turret {
	return Turret{Angle: -math.Pi / 2} // pointing up
}

func (t *Turret) Update(g *Game) {
	if !g.Player.Alive {
		return
	}

	// Rotation.
	rotSpeed := TurretRotSpeed / 60.0 // per frame
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		t.Angle -= rotSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		t.Angle += rotSpeed
	}

	// Cooldown.
	if t.Cooldown > 0 {
		t.Cooldown--
		t.Heat = float64(t.Cooldown) / float64(CooldownTime)
		if t.Cooldown == 0 {
			t.Heat = 0
		}
		return
	}

	// Heat decay.
	if t.Heat > 0 {
		t.Heat -= HeatDecay
		if t.Heat < 0 {
			t.Heat = 0
		}
	}

	// Fire timer.
	if t.FireTimer > 0 {
		t.FireTimer--
	}

	// Firing.
	if ebiten.IsKeyPressed(ebiten.KeySpace) && t.FireTimer == 0 && t.Cooldown == 0 {
		t.FireTimer = FireRate
		t.Heat += HeatPerShot
		if t.Heat >= 1.0 {
			t.Heat = 1.0
			t.Cooldown = CooldownTime
			g.Events = append(g.Events, Event{Type: EventOverheat, X: g.Player.X, Y: g.Player.Y})
			return
		}
		// Spawn projectile.
		dx := math.Cos(t.Angle)
		dy := math.Sin(t.Angle)
		spawnX := g.Player.X + dx*TurretLength
		spawnY := g.Player.Y + dy*TurretLength
		g.Projectiles = append(g.Projectiles, Projectile{
			X: spawnX, Y: spawnY,
			VX: dx * ProjectileSpeed, VY: dy * ProjectileSpeed,
			Alive: true,
		})
		g.Events = append(g.Events, Event{Type: EventFired, X: spawnX, Y: spawnY})
	}
}

func (t *Turret) Draw(screen *ebiten.Image, g *Game, ox, oy float64) {
	if !g.Player.Alive {
		return
	}
	cx := float32(g.Player.X + ox)
	cy := float32(g.Player.Y + oy)
	dx := float32(math.Cos(t.Angle))
	dy := float32(math.Sin(t.Angle))

	// Barrel.
	endX := cx + dx*TurretLength
	endY := cy + dy*TurretLength

	// Color shifts with heat.
	barrelColor := lerpColor(ColorPlayer, ColorHeatHot, float32(t.Heat))
	vector.StrokeLine(screen, cx, cy, endX, endY, 5, barrelColor, false)

	// Muzzle dot.
	vector.DrawFilledCircle(screen, endX, endY, 5, barrelColor, false)

	// Aim line (subtle).
	aimLen := float32(60.0)
	aimEndX := endX + dx*aimLen
	aimEndY := endY + dy*aimLen
	aimColor := ColorBorderDim
	vector.StrokeLine(screen, endX, endY, aimEndX, aimEndY, 1, aimColor, false)
}
