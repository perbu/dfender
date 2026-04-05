# dFender

Top-down 2D arena shooter built with Ebitengine (Go). See `SPEC.md` for full game design.
No external assets required — all graphics are procedural (vector shapes + shaders).

## Architecture

Single-package game under `game/`, entry point in `main.go`.

**Frame loop** (`game.go`): Systems run in order (input → physics → collision → wave → turret), then a single event drain loop dispatches reactions. Systems communicate via shared state on `*Game` and an `[]Event` queue — no dependency injection, no pub/sub.

Key files:
- `game.go` — Game struct, Update/Draw/Layout, arena drawing, event drain
- `player.go` — ship physics, WASD thruster input, wall collision
- `turret.go` — turret rotation, firing, heat/overheat
- `projectile.go` — bullet pool, movement
- `enemy.go` — enemy types, homing behavior
- `wave.go` — wave scheduling, enemy spawning from gates
- `collision.go` — projectile-enemy, enemy-player (squared distance, no sqrt)
- `particle.go` — particle system for explosions and thrust exhaust
- `score.go` — scoring with combo multiplier
- `event.go` — event types (value structs, no interfaces)
- `shaders.go` — shader loading, bloom pipeline, gate portals, heat distortion
- `ui.go` — HUD (score, heat bar, wave indicator, game over)
- `util.go` — drawPolygon, lerpColor helpers

**Shaders** (`game/shader/*.kage`): Ebitengine's Kage language. `//kage:unit pixels` — all coordinates are in pixels, not normalized UVs. This matters for blur offsets.

## Conventions

- All game state lives on the `Game` struct. No globals except cached constants (`gates`, color palette).
- Slices (projectiles, enemies, particles) are compacted in-place each frame — dead entries removed, no allocation.
- Events are value types appended to `g.Events` during system updates, drained once per frame in `updatePlaying()`.
- Collision uses squared distance to avoid sqrt in hot loops.
- Draw functions assume slices are already compacted (no alive checks needed).
- Colors follow the art deco palette defined in `game.go`. Player = gold/amber, enemies = teal, background = deep navy.

## Shader Pipeline

The render path in `Draw()`:
1. Draw background shader (animated grid) to offscreen `SceneImage`
2. Draw gate portal shaders (swirl effect per gate)
3. Draw arena border, particles, projectiles, enemies, player, turret to `SceneImage`
4. If heat > 0.3: apply heat distortion post-process
5. Bloom: extract bright pixels → 3 blur passes at increasing spread → additive composite to screen
6. Draw UI on top (not bloomed)

## Tuning Constants

Gameplay constants are at the top of their respective files (e.g., `ThrustForce`, `WallDeathSpeed` in `player.go`, `FireRate`, `HeatPerShot` in `turret.go`, `EnemyBaseSpeed` in `enemy.go`). Bloom threshold is in `shaders.go` (`ApplyBloom`).
