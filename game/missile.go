package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	MissileSpeed    = 6.0
	MissileRadius   = 8.0
	MissileTurnRate = 0.06 // radians/frame
	MissileMaxCount = 9
)

type Missile struct {
	X, Y  float64
	VX, VY float64
	Angle float64
	Alive bool
}

func updateMissiles(g *Game) {
	for i := range g.Missiles {
		m := &g.Missiles[i]
		if !m.Alive {
			continue
		}

		// Find nearest enemy in forward hemisphere.
		headX := math.Cos(m.Angle)
		headY := math.Sin(m.Angle)
		bestDist := math.MaxFloat64
		bestAngle := m.Angle
		hasTarget := false

		for j := range g.Enemies {
			e := &g.Enemies[j]
			if !e.Alive {
				continue
			}
			dx := e.X - m.X
			dy := e.Y - m.Y

			// Forward hemisphere check.
			if headX*dx+headY*dy <= 0 {
				continue
			}

			dist := dx*dx + dy*dy
			if dist < bestDist {
				bestDist = dist
				bestAngle = math.Atan2(dy, dx)
				hasTarget = true
			}
		}

		// Steer toward target.
		if hasTarget {
			diff := math.Remainder(bestAngle-m.Angle, 2*math.Pi)
			if diff > MissileTurnRate {
				diff = MissileTurnRate
			} else if diff < -MissileTurnRate {
				diff = -MissileTurnRate
			}
			m.Angle += diff
		}

		m.VX = math.Cos(m.Angle) * MissileSpeed
		m.VY = math.Sin(m.Angle) * MissileSpeed
		m.X += m.VX
		m.Y += m.VY

		// Wall collision.
		if m.X < ArenaLeft() || m.X > ArenaRight() ||
			m.Y < ArenaTop() || m.Y > ArenaBottom() {
			ix := math.Max(ArenaLeft(), math.Min(m.X, ArenaRight()))
			iy := math.Max(ArenaTop(), math.Min(m.Y, ArenaBottom()))
			g.Events = append(g.Events, Event{
				Type: EventMissileWallHit, X: ix, Y: iy,
			})
			m.Alive = false
		}
	}

	// Compact.
	n := 0
	for i := range g.Missiles {
		if g.Missiles[i].Alive {
			g.Missiles[n] = g.Missiles[i]
			n++
		}
	}
	g.Missiles = g.Missiles[:n]
}

func checkMissileCollisions(g *Game) {
	for i := range g.Missiles {
		m := &g.Missiles[i]
		if !m.Alive {
			continue
		}
		for j := range g.Enemies {
			e := &g.Enemies[j]
			if !e.Alive {
				continue
			}
			dx := m.X - e.X
			dy := m.Y - e.Y
			if dx*dx+dy*dy < missileEnemyRadSq {
				m.Alive = false
				e.Alive = false
				g.Events = append(g.Events, Event{
					Type:  EventEnemyKilled,
					X:     e.X,
					Y:     e.Y,
					Value: float64(e.MaxHP) * 100,
				})
				break
			}
		}
	}
}

func fireMissile(g *Game) {
	dx := math.Cos(g.Turret.Angle)
	dy := math.Sin(g.Turret.Angle)
	spawnX := g.Player.X + dx*TurretLength
	spawnY := g.Player.Y + dy*TurretLength

	g.Missiles = append(g.Missiles, Missile{
		X: spawnX, Y: spawnY,
		VX:    dx * MissileSpeed,
		VY:    dy * MissileSpeed,
		Angle: g.Turret.Angle,
		Alive: true,
	})
	g.Events = append(g.Events, Event{
		Type: EventMissileFired, X: spawnX, Y: spawnY,
	})
}

func drawMissiles(screen *ebiten.Image, g *Game, ox, oy float64) {
	for i := range g.Missiles {
		m := &g.Missiles[i]
		cx := float32(m.X + ox)
		cy := float32(m.Y + oy)

		// Outer glow.
		vector.DrawFilledCircle(screen, cx, cy, MissileRadius+3, ColorHeatCool, false)

		// Red core diamond.
		drawPolygon(screen, cx, cy, MissileRadius, 4, m.Angle, 2, ColorHeatHot)
		vector.DrawFilledCircle(screen, cx, cy, 3, ColorHeatHot, false)

		// Trail.
		tx := float32(m.X - m.VX*1.5 + ox)
		ty := float32(m.Y - m.VY*1.5 + oy)
		vector.StrokeLine(screen, cx, cy, tx, ty, 3, ColorHeatHot, false)
	}
}
