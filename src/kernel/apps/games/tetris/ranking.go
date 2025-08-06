package tetris

// Ranking represents a structure to manage a score leaderboard with a fixed-size slice of scores.
type Ranking struct {
	scores []int
}

// NewRanking initializes a new Ranking instance, loading scores from a config file if it exists or defaulting to zeros.
func NewRanking() *Ranking {
	ranking := new(Ranking)
	ranking.scores = make([]int, 10)
	for _, i := range ranking.scores {
		ranking.scores[i] = 0
	}
	return ranking
}

// insertScore inserts a score into the ranking's score list while maintaining a sorted order.
// If the new score surpasses an existing one, it shifts scores and inserts it into the correct position.
func (r *Ranking) insertScore(sc int) {
	for idx, rsc := range r.scores {
		if rsc < sc {
			r.slideScores(idx)
			r.scores[idx] = sc
			return
		}
	}
}

// slideScores shifts all scores in the ranking to the right starting from the given index to make space for a new score.
func (r *Ranking) slideScores(index int) {
	for i := len(r.scores) - 1; i > index; i-- {
		r.scores[i] = r.scores[i-1]
	}
}
