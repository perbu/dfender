# dFender

Top-down 2D arena shooter written in Go with [Ebitengine](https://ebitengine.org/). Defend against waves of enemies pouring through four gates. Newtonian physics, a freely rotating turret, and art deco visuals with bloom and shader effects.

Concept copied from Stardew Valley arcade.

![screenshot](assets/screenshot.png)

No external assets — all graphics are procedural (vector shapes + Kage shaders).

## Controls

- **WASD** — thrust (Newtonian — you drift, brake yourself)
- **Arrow keys** — rotate turret
- **Space** — fire (overheats)
- **E** — fire missile (if carrying any)

Hit a wall too fast and you die. Touch an enemy and you die.

## Install

```
go install github.com/perbu/dfender@latest
```
It might take a minute or two to install.
Requires Go 1.26+ and the Ebitengine dependencies for your platform (see [Ebitengine install guide](https://ebitengine.org/en/documents/install.html)).

## Gameplay Mechanics

### Ship & Movement

- Newtonian physics — four thrusters (WASD) apply force impulses, no friction, no speed cap.
- Thruster exhaust produces particles and glow.
- **Soft wall impact** (below speed threshold): bounce with velocity dampening and particle burst.
- **Hard wall impact** (above speed threshold): instant death. You can reach lethal speed after ~3 seconds of single-thruster acceleration.
- One-hit-kill from enemy contact.

### Turret & Weapon

- Turret rotates independently at ~180°/sec via arrow keys.
- Machine gun fires while Space is held (~10 rounds/sec, unlimited ammo).
- Heat gauge fills over ~4 seconds of continuous fire. Overheat locks the gun for ~2 seconds. Heat dissipates in ~3 seconds when not firing.
- Turret glows from amber to red as heat builds.

### Enemies & Waves

- Enemies spawn from four gate portals (one per wall) and home toward the player.
- **Wave 1**: 1 HP enemies, small count. Each subsequent wave adds +1 HP and more enemies.
- Enemies trickle in over several seconds per wave, with a ~3 second pause between waves.
- Enemy speed increases slightly per wave.

### Power-Ups

Enemies have a 20% chance to drop a power-up on death, starting from wave 2. Power-ups float in the arena for 10 seconds (blinking before they despawn). All three types coexist independently.

- **Shield** (gold hexagon): Absorbs one enemy hit and kills the attacker. Carry one at a time; pick up another to refill.
- **Guns** (amber pentagon): 20-second buff — double barrel with faster fire rate. Picking up another resets the timer.
- **Missile** (red diamond): Adds one homing missile to your inventory (max 9). Fire with **E** — the missile homes toward the nearest enemy in the turret's direction and one-hit kills regardless of HP. Explodes on wall contact.

### Scoring

- 100 base points per kill.
- Combo multiplier: kills within 2 seconds of each other increment the multiplier (x2, x3, ...). Resets after 2 seconds of no kills.
- Wave completion bonus: 500 × wave number.

### Visual Style

Art deco palette — deep navy background, gold/amber player, teal enemies. Heavy shader work: bloom on everything luminous, gate portal swirls, heat distortion near overheat, screen shake on impacts and deaths. All graphics are procedural vector shapes rendered at 1920×1080.

## License

Apache Public License 2.0
