package context

import (
	"fmt"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"sort"
	"strings"
)

const pathSeparator = "/"

// Suggestion provides mechanisms to handle and generate command suggestion operations in a hierarchical command system.
type Suggestion struct {
	root   interfaces.ICommand
	system interfaces.ICommand
}

// NewSuggestion creates and returns a new instance of Suggestion with the specified root and system ICommand instances.
func NewSuggestion(root interfaces.ICommand, system interfaces.ICommand) *Suggestion {
	if root == nil {
		panic("root command is nil")
	}
	if system == nil {
		panic("system command is nil")
	}
	return &Suggestion{
		root:   root,
		system: system,
	}
}

// Get generates suggestions based on the input string and cursor position within the command hierarchy.
// It returns the prefix to complete, a list of suggestions, and a boolean indicating if suggestions were found.
func (c *Suggestion) Get(cwd interfaces.ICommand, in string) (string, []string, bool) {
	textBeforeSegment, cmd, prefix, basePath, isCompletingCommand, err := c.parseInput(in, cwd)
	if err != nil || cmd == nil {
		return "", nil, false
	}

	var rawSuggestions []string
	if isCompletingCommand {
		rawSuggestions = cmd.SuggestionsFor_NEW(prefix)
		if !strings.Contains(prefix, interfaces.PathSeparator) && c.system != nil {
			coreSuggestions := c.system.SuggestionsFor_NEW(prefix)
			rawSuggestions = c.mergeSuggestions(rawSuggestions, coreSuggestions)
		}
	} else {
		rawSuggestions = cmd.SuggestionsFor_NEW(prefix)
	}

	if rawSuggestions == nil || len(rawSuggestions) == 0 {
		return prefix, nil, false
	}

	suggestions := make([]string, 0, len(rawSuggestions))
	for _, rawSuggestion := range rawSuggestions {
		fullSuggestion := ""
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

		if suggestionNode := cmd.FindChildren(rawSuggestion); suggestionNode != nil {
			if suggestionNode.HasSubCommands() {
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
	found := len(suggestions) > 0
	return prefix, suggestions, found
}

// parseInput identifies the command node, base path, and completion prefix from user input and cursor position.
// It returns the node to query, the prefix to complete, the base path, whether completing a command, and an error if any.
// Validation errors arise for invalid cursor positions, nil root, or current working directory commands.
// Traverses input paths to locate and validate nodes, supporting both absolute and relative paths.
// Handles completion for commands and subcommands by analyzing the input and directory structure.
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

// mergeSuggestions combines two slices of suggestion strings into a single deduplicated slice.
// If either input slice is nil, the other is deduplicated and returned.
// Deduplication is achieved by using a map to track unique entries.
// Returns the merged and deduplicated slice of suggestions.
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

// deduplicateSuggestions removes duplicate strings from the provided slice and returns a slice with unique elements in order.
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
