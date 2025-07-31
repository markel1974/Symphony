package tetris

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Ranking represents a structure to manage a score leaderboard with a fixed-size slice of scores.
type Ranking struct {
	scores []int
}

// NewRanking initializes a new Ranking instance, loading scores from a config file if it exists or defaulting to zeros.
func NewRanking() *Ranking {
	ranking := new(Ranking)
	ranking.scores = make([]int, 10)

	if fileExists(configFilePath()) {
		buf, err := ioutil.ReadFile(configFilePath())
		if err != nil {
			log.Fatal(err)
		}

		scoreTexts := strings.Split(string(buf), ",")
		for idx, text := range scoreTexts {
			num, err := strconv.Atoi(text)
			if err != nil {
				log.Fatal(err)
			}
			ranking.scores[idx] = num
		}
	} else {
		for i := 0; i < 10; i++ {
			ranking.scores[i] = 0
		}
	}
	return ranking
}

// save writes the current scores of the Ranking to a configuration file as a comma-separated string.
func (r *Ranking) save() {
	var texts []string
	for _, sc := range r.scores {
		texts = append(texts, strconv.Itoa(sc))
	}
	config := strings.Join(texts, ",")
	_ = ioutil.WriteFile(configFilePath(), []byte(config), 0644)
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

// fileExists checks if a file exists at the given filename path. Returns true if the file exists, otherwise false.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// configFilePath returns the file path to the user's configuration file located in their home directory.
func configFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, ".tetris")
}
