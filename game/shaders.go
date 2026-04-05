package game

import (
	_ "embed"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shader/bloom_bright.kage
var bloomBrightSrc []byte

//go:embed shader/blur_h.kage
var blurHSrc []byte

//go:embed shader/blur_v.kage
var blurVSrc []byte

//go:embed shader/background.kage
var backgroundSrc []byte

//go:embed shader/gate.kage
var gateSrc []byte

//go:embed shader/heat_distort.kage
var heatDistortSrc []byte

//go:embed shader/starfield.kage
var starfieldSrc []byte

// Shaders holds all compiled shaders and offscreen images for the render pipeline.
type Shaders struct {
	BloomBright  *ebiten.Shader
	BlurH        *ebiten.Shader
	BlurV        *ebiten.Shader
	Background   *ebiten.Shader
	Gate         *ebiten.Shader
	HeatDistort  *ebiten.Shader
	Starfield    *ebiten.Shader

	// Offscreen images for bloom pipeline.
	SceneImage  *ebiten.Image // full-res scene
	BrightImage *ebiten.Image // half-res bright extraction
	BlurTemp    *ebiten.Image // half-res blur intermediate
	BloomImage  *ebiten.Image // half-res final bloom

	// Gate portal rendering.
	GateImage  *ebiten.Image // render target
	GateDummy  *ebiten.Image // dummy source (different from target)

	// Temp image for heat distortion pass.
	HeatTemp *ebiten.Image
}

func NewShaders() *Shaders {
	s := &Shaders{}

	var err error
	s.BloomBright, err = ebiten.NewShader(bloomBrightSrc)
	if err != nil {
		log.Fatalf("bloom_bright shader: %v", err)
	}
	s.BlurH, err = ebiten.NewShader(blurHSrc)
	if err != nil {
		log.Fatalf("blur_h shader: %v", err)
	}
	s.BlurV, err = ebiten.NewShader(blurVSrc)
	if err != nil {
		log.Fatalf("blur_v shader: %v", err)
	}
	s.Background, err = ebiten.NewShader(backgroundSrc)
	if err != nil {
		log.Fatalf("background shader: %v", err)
	}
	s.Gate, err = ebiten.NewShader(gateSrc)
	if err != nil {
		log.Fatalf("gate shader: %v", err)
	}
	s.HeatDistort, err = ebiten.NewShader(heatDistortSrc)
	if err != nil {
		log.Fatalf("heat_distort shader: %v", err)
	}
	s.Starfield, err = ebiten.NewShader(starfieldSrc)
	if err != nil {
		log.Fatalf("starfield shader: %v", err)
	}

	// Offscreen images.
	s.SceneImage = ebiten.NewImage(ScreenWidth, ScreenHeight)
	s.BrightImage = ebiten.NewImage(ScreenWidth, ScreenHeight)
	s.BlurTemp = ebiten.NewImage(ScreenWidth, ScreenHeight)
	s.BloomImage = ebiten.NewImage(ScreenWidth, ScreenHeight)
	s.GateImage = ebiten.NewImage(200, 200)
	s.GateDummy = ebiten.NewImage(200, 200)
	s.HeatTemp = ebiten.NewImage(ScreenWidth, ScreenHeight)

	return s
}

// DrawStarfield renders the animated starfield for the menu screen.
func (s *Shaders) DrawStarfield(dst *ebiten.Image, tick uint64) {
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"Time":       float32(tick) / 60.0,
		"Resolution": []float32{float32(ScreenWidth), float32(ScreenHeight)},
	}
	dst.DrawRectShader(ScreenWidth, ScreenHeight, s.Starfield, opts)
}

// DrawBackground renders the animated art deco background.
func (s *Shaders) DrawBackground(dst *ebiten.Image, tick uint64) {
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"Time":       float32(tick) / 60.0,
		"Resolution": []float32{float32(ScreenWidth), float32(ScreenHeight)},
	}
	dst.DrawRectShader(ScreenWidth, ScreenHeight, s.Background, opts)
}

// DrawGatePortal renders a swirling gate effect.
func (s *Shaders) DrawGatePortal(dst *ebiten.Image, gate Gate, tick uint64, spawning bool) {
	size := 200
	s.GateImage.Clear()

	intensity := float32(0.4)
	if spawning {
		intensity = 0.8 + 0.2*float32(math.Sin(float64(tick)*0.1))
	}

	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"Time":      float32(tick) / 60.0,
		"Intensity": intensity,
		"GateColor": []float32{0.0, 0.788, 0.655}, // teal
	}
	opts.Images[0] = s.GateDummy // source must differ from target
	s.GateImage.DrawRectShader(size, size, s.Gate, opts)

	// Draw gate image centered on gate position.
	drawOpts := &ebiten.DrawImageOptions{}
	drawOpts.GeoM.Translate(-float64(size)/2, -float64(size)/2)
	drawOpts.GeoM.Translate(gate.X, gate.Y)
	drawOpts.Blend = ebiten.BlendSourceOver
	dst.DrawImage(s.GateImage, drawOpts)
}

// blurPass runs one horizontal+vertical blur at the given spread.
// Reads from src, writes to BloomImage (via BlurTemp intermediate).
func (s *Shaders) blurPass(src *ebiten.Image, spread float32) {
	w, h := ScreenWidth, ScreenHeight

	s.BlurTemp.Clear()
	hOpts := &ebiten.DrawRectShaderOptions{}
	hOpts.Uniforms = map[string]any{
		"Spread": spread,
	}
	hOpts.Images[0] = src
	s.BlurTemp.DrawRectShader(w, h, s.BlurH, hOpts)

	s.BloomImage.Clear()
	vOpts := &ebiten.DrawRectShaderOptions{}
	vOpts.Uniforms = map[string]any{
		"Spread": spread,
	}
	vOpts.Images[0] = s.BlurTemp
	s.BloomImage.DrawRectShader(w, h, s.BlurV, vOpts)
}

// ApplyBloom uses cascaded blur passes for a smooth, wide glow.
// Each pass blurs the previous result, doubling the effective radius
// without the ghost artifacts of single-pass large offsets.
func (s *Shaders) ApplyBloom(dst *ebiten.Image) {
	w, h := ScreenWidth, ScreenHeight

	// 1. Extract bright pixels.
	s.BrightImage.Clear()
	brightOpts := &ebiten.DrawRectShaderOptions{}
	brightOpts.Uniforms = map[string]any{
		"Threshold": float32(0.08),
	}
	brightOpts.Images[0] = s.SceneImage
	s.BrightImage.DrawRectShader(w, h, s.BloomBright, brightOpts)

	// 2. Composite scene first.
	dst.DrawImage(s.SceneImage, nil)

	// 3. Cascaded blur: blur the blur result each iteration.
	// Pass 1: BrightImage → BloomImage (small blur, tight glow).
	s.blurPass(s.BrightImage, 2.0)
	// Copy result to BrightImage for next pass input.
	s.BrightImage.Clear()
	s.BrightImage.DrawImage(s.BloomImage, nil)
	// Composite tight glow.
	bloomOpts := &ebiten.DrawImageOptions{}
	bloomOpts.Blend = ebiten.BlendLighter
	dst.DrawImage(s.BloomImage, bloomOpts)

	// Pass 2: blur the already-blurred image (medium glow).
	s.blurPass(s.BrightImage, 3.0)
	s.BrightImage.Clear()
	s.BrightImage.DrawImage(s.BloomImage, nil)
	bloomOpts2 := &ebiten.DrawImageOptions{}
	bloomOpts2.Blend = ebiten.BlendLighter
	bloomOpts2.ColorScale.Scale(0.6, 0.6, 0.6, 0.6)
	dst.DrawImage(s.BloomImage, bloomOpts2)

	// Pass 3: blur again (wide, soft outer glow).
	s.blurPass(s.BrightImage, 4.0)
	bloomOpts3 := &ebiten.DrawImageOptions{}
	bloomOpts3.Blend = ebiten.BlendLighter
	bloomOpts3.ColorScale.Scale(0.3, 0.3, 0.3, 0.3)
	dst.DrawImage(s.BloomImage, bloomOpts3)
}

// ApplyHeatDistortion applies screen-space heat distortion.
func (s *Shaders) ApplyHeatDistortion(dst, src *ebiten.Image, heat float64, centerX, centerY float64, tick uint64) {
	if heat < 0.3 {
		// Below threshold, just copy.
		dst.DrawImage(src, nil)
		return
	}
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"HeatAmount": float32(heat),
		"HeatCenter": []float32{float32(centerX), float32(centerY)},
		"Resolution": []float32{float32(ScreenWidth), float32(ScreenHeight)},
		"Time":       float32(tick) / 60.0,
	}
	opts.Images[0] = src
	dst.DrawRectShader(ScreenWidth, ScreenHeight, s.HeatDistort, opts)
}
