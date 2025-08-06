package invaders

import (
	"sort"
)

// highScoreFilename defines the filename used to store high scores.
// highScoreSeparator is the delimiter used between high scores in the file.
// maxHighScores specifies the maximum number of high scores to keep.
const (
	//highScoreFilename  = "hs"
	//highScoreSeparator = ":"
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
		//g.saveHighScores()
	}
}

//func (g *Invaders) saveHighScores() {
//data := ""
//for i, score := range g.highScores {
//	data += fmt.Sprintf("%s%s%d", score.name, highScoreSeparator, score.score)
//	if i != len(g.highScores)-1 {
//		data += "\n"
//	}
//}
//_ = os.WriteFile(highScoreFilename, []byte(data), 0666)
//}

/*
// _ parses high score data from a byte slice, validates the format, sorts the scores, and limits the list to the top 5.
func func (g *Invaders) parseHighScore (data []byte) error {
	highScores := make([]*HighScore, 0)
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no highscore data in file")
	}
	for _, l := range lines {
		parts := strings.Split(l, highScoreSeparator)
		if len(parts) != 2 {
			return fmt.Errorf("corrupted highscore data")
		}
		if i, err := strconv.Atoi(parts[1]); err == nil {
			if i < 0 {
				return fmt.Errorf("negative highscore - data corrupted")
			}
			highScores = append(highScores, &HighScore{i, parts[0]})
		} else {
			return err
		}
	}

	sort.Sort(sort.Reverse(ByScore(highScores)))
	if len(highScores) > 5 {
		highScores = highScores[:5]
	}
	return nil
}
*/

/*
// loadHighScores reads and parses the high scores file, then loads the scores into memory, maintaining a sorted top 5 list.
// Any data corruption issues in the file are logged, and invalid entries are ignored to ensure data integrity.
func (g *Invaders) loadHighScores() {
	data, err := os.ReadFile(highScoreFilename)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return
	}
	for _, l := range lines {
		parts := strings.Split(l, highScoreSeparator)
		if len(parts) != 2 {
			log.Println("high scores file has been corrupted - please correct/delete it")
			continue
		} else if len(parts[0]) < 3 || len(parts[0]) > 10 {
			log.Println("high score file has been corrupted (name too long/short) - please correct/delete it")
			continue
		}
		if i, err := strconv.Atoi(parts[1]); err == nil {
			if i < 0 {
				log.Println("negative high score found - data corrupted - please correct/delete the hs file")
				continue
			}
			g.highScores = append(g.highScores, &HighScore{i, parts[0]})
		} else {
			log.Fatalln(err)
		}
	}

	sort.Sort(sort.Reverse(ByScore(g.highScores)))
	if len(g.highScores) > 5 {
		g.highScores = g.highScores[:5]
	}
}
*/
