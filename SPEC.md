# dFender — Game Specification

## Overview

Top-down 2D arena shooter. You pilot a spaceship defending against waves of enemies pouring through four gates. Newtonian physics, a rotating turret, and art deco visuals with heavy shader work.

Engine: **Ebitengine** (Go). Target: **1920x1080** windowed.

---

## Arena

- Rectangular, fills the entire window (1920x1080).
- Four gates centered on each wall: **North, South, West, East**.
- Gates are openings in the arena border — visually distinct (glowing portals).
- Interior is open — no obstacles, no cover.
- The arena border is a solid wall. Collision behavior depends on impact speed (see Ship).

---

## Player Ship

### Appearance
- Clean vector style. Body is a geometric shape (hexagon or diamond), always oriented upward on screen — it does not rotate.
- A turret sits on top, visually distinct, and rotates freely 360 degrees.
- Art deco gold/amber color with glow outline.

### Movement — Newtonian Physics
- Four thrusters: North (W), South (S), West (A), East (D). Each applies a force impulse in that cardinal direction.
- No friction in open space — the ship drifts indefinitely.
- Thruster firing produces visible exhaust particles and glow (shader effect).
- There is no artificial max speed cap.

### Wall Collision
- **Soft impact** (below speed threshold): ship bounces off the wall with velocity dampening (e.g. 0.5x reflection). Small particle burst.
- **Hard impact** (above speed threshold): ship is destroyed. Death explosion effect. Game over.
- The lethal speed threshold should be tuned so it's possible to ram a wall at full single-thruster speed after ~3 seconds of acceleration. Players must actively brake.

### Death
- One-hit-kill from enemy contact or projectiles.
- One-hit-kill from high-speed wall collision.
- Death triggers an explosion effect and game-over screen with final score.

---

## Turret & Weapon

### Rotation
- Turret rotates continuously while arrow key is held. Rotation speed: ~180 degrees/second (tunable).
- Visual: turret barrel with a subtle aim line or glow showing direction.

### Firing — Machine Gun
- Fires while SPACE is held.
- High rate of fire (~10 rounds/sec).
- Projectiles are small, bright, fast-moving. Gold/amber with a short trail (shader glow).
- Projectiles travel in a straight line until hitting an enemy or leaving the arena.
- Unlimited ammo.

### Overheating
- A heat gauge fills while firing (~4 seconds to overheat at continuous fire).
- When fully overheated, the gun **locks** for a cooldown period (~2 seconds).
- Heat dissipates when not firing (~3 seconds from full to zero).
- Visual: turret glows from amber to red as heat builds. At overheat, a brief flash/steam effect.

---

## Enemies

### General
- Enemies spawn from the four gates in waves.
- All enemies move toward the player (basic homing behavior).
- Contact with the player kills the player.

### Wave Progression
- **Wave 1**: Enemies have 1 HP. Small number (e.g. 8 enemies, 2 per gate).
- **Subsequent waves**: Enemy HP increases by 1 per wave. Count increases gradually.
- Spawn timing: enemies trickle in over the first few seconds of a wave, not all at once.
- Brief pause between waves (~3 seconds) with a "WAVE N" announcement.

### Enemy Appearance
- Geometric shapes — triangles or angular forms.
- Teal/cyan color with glow. Distinct from the player's gold palette.
- When damaged (HP > 1), a brief hit flash (white).
- Death: small particle explosion in teal/magenta.

### Enemy Behavior (Initial)
- Move directly toward the player at a fixed speed.
- Speed increases slightly per wave.
- No shooting, no evasion — just relentless approach.
- Future: add enemy variants (shooters, fast flankers, tanks). Out of scope for v1.

---

## Controls

| Key       | Action                          |
|-----------|---------------------------------|
| W         | Fire north thruster (push south → move up) |
| S         | Fire south thruster (push north → move down) |
| A         | Fire west thruster (push east → move left) |
| D         | Fire east thruster (push west → move right) |
| Left arrow | Rotate turret counter-clockwise |
| Right arrow | Rotate turret clockwise        |
| Space     | Fire weapon                     |
| Escape    | Pause / Quit                    |

Note: thrusters apply force in the opposite direction of the key's label. W fires the north thruster, pushing the ship southward on the thrust vector — but since "up" on screen is north, pressing W moves the ship upward. This is intuitive: WASD maps to movement direction, not thruster position.

---

## Visual Style

### Color Palette (Art Deco)
- **Background**: Deep navy (#0A0E27) to near-black.
- **Arena border**: Thin gold lines with subtle glow.
- **Player**: Gold/amber (#D4A843, #F5D67B).
- **Enemies**: Teal (#00C9A7) and magenta (#C850C0) accents.
- **Projectiles**: Bright amber with white-hot core.
- **UI text**: Warm white (#F0E6D3), art deco font style if possible.

### Shader Effects (go overboard)
- **Bloom/Glow**: Everything luminous gets a bloom pass. Gates, projectiles, ship outlines, explosions.
- **Thruster exhaust**: Particle system with additive blending. Flame-like gradient from white core to amber to transparent.
- **Projectile trails**: Short glowing tails using motion blur or trail rendering.
- **Gate portals**: Swirling/pulsing energy effect. Intensity increases before a wave spawns.
- **Explosions**: Radial burst with light flash, expanding ring, fading particles.
- **Heat distortion**: When gun is near overheat, subtle screen warp near the turret.
- **Screen shake**: On death, on hard wall impacts, on big explosions.

### Resolution & Rendering
- Game renders at 1920x1080.
- All game objects are drawn procedurally (no sprite sheets) — vector shapes, lines, and shader effects.
- This keeps the art style consistent and resolution-independent.

---

## UI (Minimal)

- **Heat gauge**: Thin bar near the bottom center. Fills amber → red.
- **Score**: Top-right corner. Simple number.
- **Wave indicator**: Top-center. "WAVE 3" appears briefly at wave start, then fades.
- **Game Over screen**: Final score, "PRESS ENTER TO RESTART".
- No health bar (one-hit-kill). No minimap. No ammo counter.

---

## Scoring

- Points per enemy kill. Base: 100 points.
- Multiplier for consecutive kills without pause (combo system): kills within 2 seconds of each other increment a multiplier (x2, x3, ...). Resets after 2 seconds of no kills.
- Wave completion bonus: 500 * wave number.
- Score displayed on game-over screen.

---

## Sound (minimal, placeholder-quality is fine)

- Thruster hum (looping, per-thruster).
- Gun firing (rapid staccato).
- Enemy hit / enemy death.
- Player death explosion.
- Overheat warning beep.
- Wave announcement chime.

---

## Architecture (suggested)

```
dfender/
  main.go              — entry point, ebiten.RunGame
  game/
    game.go            — Game struct, Update/Draw/Layout
    arena.go           — arena boundaries, wall collision
    player.go          — ship physics, input handling
    turret.go          — turret rotation, firing, heat
    projectile.go      — bullet pool, movement, collision
    enemy.go           — enemy types, behavior, spawning
    wave.go            — wave management, scheduling
    particle.go        — particle system for effects
    collision.go       — collision detection (circles/AABB)
    score.go           — scoring, combo tracking
  shader/
    bloom.kage         — bloom/glow post-process
    trail.kage         — projectile trail effect
    gate.kage          — gate portal swirl
    heat.kage          — heat distortion
    background.kage    — animated background (subtle grid/stars)
  assets/
    sounds/            — .wav or .ogg files
  SPEC.md
```

---

## Event Architecture

Central state on `Game` struct, with a frame-scoped event queue (`[]Event`). No dependency injection, no pub/sub registration — just functions taking `*Game`.

### Frame Loop Order

1. Clear event queue from previous frame.
2. Run systems (input → physics → collisions → wave → turret). Each system reads/writes `Game` state directly and appends events to `g.Events`.
3. Drain event queue in a single pass — a switch dispatches reactions (particles, score, screen shake, state transitions).
4. Update particles and effects.

### Event Types

| Event             | Emitted by   | Reactions                                  |
|-------------------|--------------|--------------------------------------------|
| `EnemyKilled`     | collision    | Add score, combo update, death particles   |
| `EnemyHit`        | collision    | Hit flash, damage number                   |
| `PlayerDied`      | collision    | Screen shake, game over state              |
| `WallBounce`      | physics      | Particle burst, small screen shake         |
| `WallDeath`       | physics      | Explosion, game over state                 |
| `WaveComplete`    | wave         | Start next wave timer                      |
| `WaveStart`       | wave         | "WAVE N" announcement, gate activation     |
| `Overheat`        | turret       | Gun lock, steam effect                     |
| `CooldownReady`   | turret       | Gun unlock visual                          |

### Design Rules

- Events are value types (struct with `Type`, `X`, `Y`, `Value` fields). No interfaces.
- Event slice is reused each frame — zero allocation in steady state.
- All event reactions live in one place (the drain loop in `Game.Update`). This is the only place that crosses system boundaries.
- Systems never call each other. They communicate exclusively through shared state and events.

---

## Open Questions / Future

- Enemy variants (shooters, fast flankers) — defer to v2.
- Power-ups (shield, weapon upgrades, speed boost) — defer to v2.
- Leaderboard / high score persistence — defer.
- Fullscreen toggle — easy to add later.
- Gamepad support — defer.
