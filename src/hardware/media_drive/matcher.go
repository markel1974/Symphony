package media_drive

import "strings"

// Matcher represents a structure for handling wildcard-based string pattern matching.
// S defines the single-character wildcard symbol.
// M defines the multi-character wildcard symbol.
type Matcher struct {
	s byte
	m byte
}

// NewMatcher creates and returns a new instance of Matcher with default wildcard characters '?' and '*'.
func NewMatcher() *Matcher {
	return &Matcher{
		s: '?',
		m: '*',
	}
}

// isWildPattern determines if the given pattern contains wildcards defined by the Matcher instance.
func (m *Matcher) isWildPattern(pattern string) bool {
	for i := range pattern {
		c := pattern[i]
		if c == m.m {
			return true
		}
		if c == m.s {
			return true
		}
	}
	return false
}

func (m *Matcher) Contains(plainName string) bool {
	return strings.Contains(plainName, string(m.m)) || strings.Contains(plainName, string(m.s))
}

// Match checks if the given string `s` matches the `pattern` using wildcard characters defined in the Matcher instance.
// Returns a boolean indicating the match result and an error if applicable.
func (m *Matcher) Match(pattern string, s string) bool {
	if pattern == string(m.m) {
		return true
	}
	if pattern == "" {
		if s == "" {
			return true
		}
		return false
	}
	if !m.isWildPattern(pattern) {
		return pattern == s
	}
	lp := len(pattern)
	ls := len(s)
	dp := make([][]bool, lp+1)
	for i := 0; i < lp+1; i++ {
		dp[i] = make([]bool, ls+1)
	}
	dp[0][0] = true
	for i := 0; i < lp; i++ {
		if pattern[i] == m.m {
			dp[i+1][0] = dp[i][0]
		} else {
			dp[i+1][0] = false
		}
	}
	for j := 0; j < ls; j++ {
		dp[0][j+1] = false
	}
	for i := 0; i < lp; i++ {
		for j := 0; j < ls; j++ {
			pc := pattern[i]
			sc := s[j]
			switch pattern[i] {
			case m.m:
				dp[i+1][j+1] = dp[i][j] || dp[i][j+1] || dp[i+1][j]
			case m.s:
				dp[i+1][j+1] = dp[i][j]
			default:
				if pc == sc {
					dp[i+1][j+1] = dp[i][j]
				} else {
					dp[i+1][j+1] = false
				}
			}
		}
	}
	return dp[lp][ls]
}
