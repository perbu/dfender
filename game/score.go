package game

const (
	ComboWindow = 120 // frames (2 seconds)
)

type ScoreTracker struct {
	Score      int
	Combo      int
	ComboTimer int // frames since last kill
}

func (s *ScoreTracker) AddKill(points int) {
	if s.ComboTimer > 0 && s.ComboTimer < ComboWindow {
		s.Combo++
	} else {
		s.Combo = 1
	}
	s.ComboTimer = 1 // reset timer — will count up each frame
	s.Score += points * s.Combo
}

func (s *ScoreTracker) AddWaveBonus(wave int) {
	s.Score += 500 * wave
}

func (s *ScoreTracker) Update() {
	if s.ComboTimer > 0 {
		s.ComboTimer++
		if s.ComboTimer > ComboWindow {
			s.Combo = 0
			s.ComboTimer = 0
		}
	}
}
