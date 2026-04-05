package game

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"time"
)

const maxHighScores = 10

type HighScoreEntry struct {
	Name  string    `json:"name"`
	Score int       `json:"score"`
	Wave  int       `json:"wave"`
	Date  time.Time `json:"date"`
}

type HighScoreTable struct {
	Entries []HighScoreEntry `json:"entries"`
}

func highScorePath() string {
	dir := filepath.Join(os.Getenv("HOME"), ".cache", "dfender")
	if dir == filepath.Join("", ".cache", "dfender") {
		// Fallback if HOME is not set.
		if u, err := user.Current(); err == nil {
			dir = filepath.Join(u.HomeDir, ".cache", "dfender")
		}
	}
	return filepath.Join(dir, "highscores.json")
}

func LoadHighScores() HighScoreTable {
	data, err := os.ReadFile(highScorePath())
	if err != nil {
		return HighScoreTable{}
	}
	var table HighScoreTable
	if err := json.Unmarshal(data, &table); err != nil {
		return HighScoreTable{}
	}
	return table
}

func (t *HighScoreTable) Save() error {
	path := highScorePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (t *HighScoreTable) Qualifies(score int) bool {
	if len(t.Entries) < maxHighScores {
		return score > 0
	}
	return score > t.Entries[len(t.Entries)-1].Score
}

func (t *HighScoreTable) Add(entry HighScoreEntry) {
	t.Entries = append(t.Entries, entry)
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Score > t.Entries[j].Score
	})
	if len(t.Entries) > maxHighScores {
		t.Entries = t.Entries[:maxHighScores]
	}
}

func getUsername() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	return "PLAYER"
}
