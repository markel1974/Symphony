/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"fmt"
	"strings"
	"unicode"
)

// commandSorterByName is a type that represents a slice of *Command used for sorting commands by their names.

// EnablePrefixMatching allows setting automatic prefix matching.
// Automatic prefix matching can be a dangerous thing to automatically enable in CLI tools.
// Set this to true to enable it.
//var EnablePrefixMatching = false

// trimRightSpace removes all trailing whitespace characters from the input string and returns the modified string.
func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// rPad right-pads the input string 's' with spaces to match the specified 'padding' length.
func rPad(s string, padding int) string {
	t := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(t, s)
}

// levenshteinDistance computes the Levenshtein distance between two strings s and s2.
// If ignoreCase is true, the comparison is case-insensitive.
func levenshteinDistance(s string, s2 string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		s2 = strings.ToLower(s2)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(s2)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(s2); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == s2[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(s2)]
}

// argsMinusFirstX removes the first occurrence of the string x from the args slice and returns the modified slice.
func argsMinusFirstX(args []string, x string) []string {
	for i, y := range args {
		if x == y {
			var ret []string
			ret = append(ret, args[:i]...)
			ret = append(ret, args[i+1:]...)
			return ret
		}
	}
	return args
}
