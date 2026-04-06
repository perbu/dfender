package game

import (
	"math"
	"sync/atomic"
)

// EngineSound is a streaming audio source that generates a composite sine wave
// for the thruster. Always outputs at full amplitude — volume control is done
// via Player.SetVolume() at the mixer to avoid buffer latency.
type EngineSound struct {
	pos  int64
	freq atomic.Int64 // fundamental freq * 100 (fixed-point)
}

func NewEngineSound() *EngineSound {
	s := &EngineSound{}
	s.freq.Store(11000) // 110.00 Hz
	return s
}

// SetFreq updates the fundamental frequency. Called from the game loop.
func (s *EngineSound) SetFreq(f float64) {
	s.freq.Store(int64(f * 100))
}

// Read fills buf with stereo signed-16-bit PCM at constant amplitude.
func (s *EngineSound) Read(buf []byte) (int, error) {
	freq := float64(s.freq.Load()) / 100.0
	pos := s.pos

	n := len(buf) / 4 // 2 bytes per channel, 2 channels
	for i := 0; i < n; i++ {
		t := float64(pos) / float64(sampleRate)

		v := math.Sin(2*math.Pi*freq*t) * 0.4
		v += math.Sin(2*math.Pi*freq*2.0*t) * 0.2
		v += math.Sin(2*math.Pi*freq*3.5*t) * 0.1

		sample := int16(v * 32767)

		buf[4*i] = byte(sample)
		buf[4*i+1] = byte(sample >> 8)
		buf[4*i+2] = byte(sample)
		buf[4*i+3] = byte(sample >> 8)
		pos++
	}

	s.pos = pos
	return len(buf), nil
}
