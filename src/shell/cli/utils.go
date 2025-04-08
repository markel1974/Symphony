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
	"text/template"
	"unicode"
)

var _templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"rPad":                    rPad,
}

// commandSorterByName is a type that represents a slice of *Command used for sorting commands by their names.

var _initializers []func()

// EnablePrefixMatching allows setting automatic prefix matching.
// Automatic prefix matching can be a dangerous thing to automatically enable in CLI tools.
// Set this to true to enable it.
//var EnablePrefixMatching = false

// AddTemplateFunc adds a template function that's available to Usage and Help template generation.
func _(name string, tmplFunc interface{}) {
	_templateFuncs[name] = tmplFunc
}

// AddTemplateFuncs adds multiple template functions that are available to Usage and Help template generation.
func _(tmplFuncs template.FuncMap) {
	for k, v := range tmplFuncs {
		_templateFuncs[k] = v
	}
}

// OnInitialize sets the passed functions to be run when each command's Execute method is called.
func _(y ...func()) {
	_initializers = append(_initializers, y...)
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func rPad(s string, padding int) string {
	t := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(t, s)
}

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// argsMinusFirstX removes the first occurrence of the string x from the slice args and returns the resulting slice.
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

// isFlagArg checks if the given string represents a flag argument, either in long format (--flag) or short format (-f).
func isFlagArg(arg string) bool {
	return (len(arg) >= 3 && arg[1] == '-') ||
		(len(arg) >= 2 && arg[0] == '-' && arg[1] != '-')
}
