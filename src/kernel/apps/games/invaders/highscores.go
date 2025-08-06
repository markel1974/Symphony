package invaders

import (
	"sort"
)

// maxHighScores specifies the maximum number of high scores to keep.
const (
	maxHighScores = 5
)

// HighScore represents a player's score and associated name in the high scores list.
type HighScore struct {
	score int
	name  string
}

// ByScore is a slice of pointers to HighScore used for sorting based on the HighScore values.
type ByScore []*HighScore

// Len returns the number of elements in the ByScore slice.
func (a ByScore) Len() int { return len(a) }

// Swap exchanges the elements at indexes i and j in the ByScore slice.
func (a ByScore) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less reports whether the element with index i should sort before the element with index j based on their scores.
func (a ByScore) Less(i, j int) bool { return a[i].score < a[j].score }

// checkHighScores checks if the player's score qualifies for the high scores and updates the high scores list accordingly.
func (g *Invaders) checkHighScores() {
	if len(g.highScores) < maxHighScores || g.player.score > g.highScores[len(g.highScores)-1].score {
		name := "Invader"
		g.highScores = append(g.highScores, &HighScore{g.player.score, name})
		sort.Sort(sort.Reverse(ByScore(g.highScores)))
		if len(g.highScores) > maxHighScores {
			g.highScores = append([]*HighScore(nil), g.highScores[:maxHighScores]...)
		}
	}
}
