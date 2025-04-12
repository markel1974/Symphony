package context

import (
	"fmt"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"sort"
	"strings"
)

// pathSeparator is a constant string representing the character used to separate components in a path, typically "/".
const pathSeparator = "/"

// Suggestion provides utilities for command hierarchy traversal and completion suggestions.
type Suggestion struct {
	root        interfaces.ICommand
	searchPaths []interfaces.ICommand
}

// NewSuggestion initializes and returns a new Suggestion instance with the given root command and an empty search path list.
func NewSuggestion(root interfaces.ICommand, searchPaths []interfaces.ICommand) *Suggestion {
	if searchPaths == nil {
		searchPaths = []interfaces.ICommand{}
	}
	return &Suggestion{
		root:        root,
		searchPaths: searchPaths,
	}
}

// AddSearchPath adds a new ICommand instance to the searchPaths slice for suggestion resolution.
func (c *Suggestion) AddSearchPath(sp interfaces.ICommand) {
	c.searchPaths = append(c.searchPaths, sp)
}

// Get generates command suggestions based on the provided input and current directory context.
// It returns the input prefix, a list of suggestions, and a boolean indicating if suggestions exist.
func (c *Suggestion) Get(cwd interfaces.ICommand, in string) (string, []string, bool) {
	textBeforeSegment, nodeToQuery, prefixToComplete, basePath, isCompletingCommand, err := c.parseInput(in, cwd)
	if err != nil || nodeToQuery == nil {
		return "", nil, false
	}
	prefix := prefixToComplete

	rawSuggestions := nodeToQuery.SuggestionsFor_NEW(prefixToComplete)
	isFromSearchPath := false

	if len(rawSuggestions) == 0 && isCompletingCommand {
		for _, searchRoot := range c.searchPaths {
			if searchRoot == nil {
				continue
			}
			pathSuggestions := searchRoot.SuggestionsFor_NEW(prefixToComplete)
			if len(pathSuggestions) > 0 {
				rawSuggestions = pathSuggestions
				isFromSearchPath = true
				break
			}
		}
	}

	if len(rawSuggestions) == 0 {
		return prefix, nil, false
	}

	suggestions := make([]string, 0, len(rawSuggestions))
	for _, rawSuggestion := range rawSuggestions {
		fullSuggestion := ""
		if isFromSearchPath {
			fullSuggestion = rawSuggestion
		} else {
			if basePath == interfaces.PathSeparator {
				fullSuggestion = interfaces.PathSeparator + rawSuggestion
			} else if basePath != "" {
				bp := basePath
				if !strings.HasSuffix(bp, interfaces.PathSeparator) {
					bp += interfaces.PathSeparator
				}
				fullSuggestion = bp + rawSuggestion
			} else {
				fullSuggestion = rawSuggestion
			}
			if sNode := nodeToQuery.FindChildren(rawSuggestion); sNode != nil && sNode.HasSubCommands() {
				if !strings.HasSuffix(fullSuggestion, interfaces.PathSeparator) {
					fullSuggestion += interfaces.PathSeparator
				}
			}
		}
		suggestions = append(suggestions, fullSuggestion)
	}

	if len(suggestions) > 1 {
		sort.Strings(suggestions)
		suggestions = c.deduplicateSuggestions(suggestions)
	}

	if len(textBeforeSegment) > 0 {
		for i, suggestion := range suggestions {
			suggestions[i] = textBeforeSegment + " " + suggestion
		}
	}
	return prefix, suggestions, len(suggestions) > 0
}

// parseInput parses the input string to determine the relevant command context, path, and completion prefix details.
// It returns the text before the path segment, the current command node, the prefix for completion, the base path,
// a boolean indicating if the input addresses a command name, and any error encountered during processing.
func (c *Suggestion) parseInput(input string, cwd interfaces.ICommand) (string, interfaces.ICommand, string, string, bool, error) {
	isCompletingCommand := false
	pathPart := ""
	textBeforeSegment := ""

	if pos := strings.LastIndex(input, " "); pos < 0 {
		pathPart = input
		isCompletingCommand = true
	} else {
		pathPart = input[pos+1:]
		textBeforeSegment = input[:pos]
		isCompletingCommand = strings.TrimSpace(textBeforeSegment) == ""
	}

	baseNode := cwd
	currentNode := baseNode
	isAbsolute := strings.HasPrefix(pathPart, pathSeparator)
	if isAbsolute {
		baseNode = c.root
		pathPart = strings.TrimPrefix(pathPart, pathSeparator)
	}

	var pathSegments []string
	for _, part := range strings.Split(pathPart, pathSeparator) {
		if len(part) > 0 {
			pathSegments = append(pathSegments, part)
		}
	}

	prefixToComplete := ""
	var dirParts []string

	if len(pathSegments) > 0 && !strings.HasSuffix(pathPart, pathSeparator) {
		prefixToComplete = pathSegments[len(pathSegments)-1]
		dirParts = pathSegments[:len(pathSegments)-1]
	} else {
		prefixToComplete = ""
		dirParts = pathSegments
	}

	var basePath string

	if len(dirParts) == 0 {
		if isAbsolute {
			basePath = interfaces.PathSeparator
		}
	} else if len(dirParts) > 0 {
		if isAbsolute {
			basePath = interfaces.PathSeparator + strings.Join(dirParts, interfaces.PathSeparator)
		} else {
			basePath = strings.Join(dirParts, interfaces.PathSeparator)
		}
		if !strings.HasSuffix(basePath, interfaces.PathSeparator) {
			basePath += interfaces.PathSeparator
		}
	}

	var traversedPathParts []string
	for _, part := range dirParts {
		if part == "" {
			if isAbsolute && len(traversedPathParts) == 0 {
				continue
			} else if !isAbsolute && len(traversedPathParts) == 0 {
				continue
			} else {
				continue
			}
		}
		foundNode := currentNode.FindChildren(part)
		if foundNode == nil {
			return "", nil, "", "", isCompletingCommand, fmt.Errorf("path not found: %s", part)
		}
		if !foundNode.HasSubCommands() && len(dirParts) > len(traversedPathParts)+1 {
			return "", nil, "", "", isCompletingCommand, fmt.Errorf("cannot traverse into non-directory: %s", part)
		}
		currentNode = foundNode
		traversedPathParts = append(traversedPathParts, part)
	}

	return textBeforeSegment, currentNode, prefixToComplete, basePath, isCompletingCommand, nil
}

// mergeSuggestions merges two slices of suggestions, removes duplicates, and returns the combined slice.
func (c *Suggestion) mergeSuggestions(s1 []string, s2 []string) []string {
	if s1 == nil {
		return c.deduplicateSuggestions(s2)
	}
	if s2 == nil {
		return c.deduplicateSuggestions(s1)
	}
	m := make(map[string]bool)
	for _, item := range s1 {
		m[item] = true
	}
	for _, item := range s2 {
		m[item] = true
	}
	merged := make([]string, 0, len(m))
	for item := range m {
		merged = append(merged, item)
	}
	return merged
}

// deduplicateSuggestions removes duplicate strings from the input slice and maintains the original order of unique elements.
func (c *Suggestion) deduplicateSuggestions(s []string) []string {
	if len(s) < 2 {
		return s
	}
	keys := make(map[string]bool)
	var list []string
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
