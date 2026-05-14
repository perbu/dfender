package game

// Event → Juice mapping:
//
//   EnemyKilled       → 30-particle teal explosion
//   EnemyHit          →  8-particle white flash
//   PlayerDied        → death explosion + 120-frame/12px shake
//   WallDeath         → same as PlayerDied
//   WallBounce        →  5-frame/3px shake + 12-particle gold explosion
//   EnemyWallDeath    →  4-frame/2px shake + 22-particle teal explosion
//   Overheat          → 15-particle red explosion
//   ProjectileWallHit → 18-particle gold explosion
//   MissileExploded   → 12-frame/6px shake + missile blast
//   MissileFired      → 12-particle red explosion
//   MinePlaced        → mine placed effect
//   MineExploded      → 16-frame/8px shake + mine blast
//   ShieldAbsorb      → 10-frame/5px shake + 45-particle gold explosion
//   PowerUpPickedUp   → 22-particle gold explosion
//
// applyJuice handles the visual and audio reactions to game events:
// particle effects, screen shake, and sound. Game mechanics (score,
// lives, state transitions) stay in the event drain in game.go.
func applyJuice(g *Game, e Event) {
	if !e.Silent {
		g.Sound.HandleEvent(e)
	}

	switch e.Type {
	case EventEnemyKilled:
		spawnExplosion(g, e.X, e.Y, ColorEnemy, 30)
	case EventEnemyHit:
		spawnExplosion(g, e.X, e.Y, ColorUI, 8)
	case EventPlayerDied, EventWallDeath:
		g.shake(120, 12)
		spawnDeathExplosion(g, e.X, e.Y)
	case EventWallBounce:
		g.shake(5, 3)
		spawnExplosion(g, e.X, e.Y, ColorBorder, 12)
	case EventEnemyWallDeath:
		g.shake(4, 2)
		spawnExplosion(g, e.X, e.Y, ColorEnemy, 22)
	case EventOverheat:
		spawnExplosion(g, e.X, e.Y, ColorHeatHot, 15)
	case EventProjectileWallHit:
		spawnExplosion(g, e.X, e.Y, ColorProjectile, 18)
	case EventMissileExploded:
		g.shake(12, 6)
		spawnMissileBlast(g, e.X, e.Y)
	case EventMissileFired:
		spawnExplosion(g, e.X, e.Y, ColorHeatHot, 12)
	case EventMinePlaced:
		spawnMinePlacedEffect(g, e.X, e.Y)
	case EventMineExploded:
		g.shake(16, 8)
		spawnMineBlast(g, e.X, e.Y)
	case EventShieldAbsorb:
		g.shake(10, 5)
		spawnExplosion(g, e.X, e.Y, ColorBorder, 45)
	case EventPowerUpPickedUp:
		spawnExplosion(g, e.X, e.Y, ColorBorder, 22)
	}
}

func (g *Game) shake(frames int, amount float32) {
	g.ShakeFrames = frames
	g.ShakeAmount = amount
}
