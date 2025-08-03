package interfaces

import "strings"

// PathSeparator is the string used to separate elements in a hierarchical path format, typically a forward slash "/".
const PathSeparator = "/"

// PathToSegments splits a hierarchical path into its individual non-empty segments based on the defined PathSeparator.
func PathToSegments(path string) []string {
	var segments []string
	for _, part := range strings.Split(path, PathSeparator) {
		if len(part) > 0 {
			segments = append(segments, part)
		}
	}
	return segments
}

// IsPathAbsolute determines if the given path string starts with the defined PathSeparator, indicating an absolute path.
func IsPathAbsolute(path string) bool {
	isAbsolute := strings.HasPrefix(path, PathSeparator)
	return isAbsolute
}
