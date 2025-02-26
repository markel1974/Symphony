package virtualdrive

// Matcher represents a structure for handling wildcard-based string pattern matching.
// S defines the single-character wildcard symbol.
// M defines the multi-character wildcard symbol.
type Matcher struct {
	S byte
	M byte
}

// NewMatcher creates and returns a new instance of Matcher with default wildcard characters '?' and '*'.
func NewMatcher() *Matcher {
	return &Matcher{
		S: '?',
		M: '*',
	}
}

// isWildPattern determines if the given pattern contains wildcards defined by the Matcher instance.
func (m *Matcher) isWildPattern(pattern string) bool {
	for i := range pattern {
		c := pattern[i]
		if c == m.M {
			return true
		}
		if c == m.S {
			return true
		}
	}
	return false
}

// Match checks if the given string `s` matches the `pattern` using wildcard characters defined in the Matcher instance.
// Returns a boolean indicating the match result and an error if applicable.
func (m *Matcher) Match(pattern string, s string) (bool, error) {
	if pattern == string(m.M) {
		return true, nil
	}
	if pattern == "" {
		if s == "" {
			return true, nil
		}
		return false, nil
	}
	if !m.isWildPattern(pattern) {
		return pattern == s, nil
	}
	lp := len(pattern)
	ls := len(s)
	dp := make([][]bool, lp+1)
	for i := 0; i < lp+1; i++ {
		dp[i] = make([]bool, ls+1)
	}
	dp[0][0] = true
	for i := 0; i < lp; i++ {
		if pattern[i] == m.M {
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
			case m.M:
				dp[i+1][j+1] = dp[i][j] || dp[i][j+1] || dp[i+1][j]
			case m.S:
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
	return dp[lp][ls], nil
}
