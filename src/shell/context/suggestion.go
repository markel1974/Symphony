package context

import (
	"errors"
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
	return &Suggestion{root: root, system: system}
}

// Get generates suggestions based on the input string and cursor position within the command hierarchy.
// It returns the prefix to complete, a list of suggestions, and a boolean indicating if suggestions were found.
func (c *Suggestion) Get(cwd interfaces.ICommand, in string, cursorPos int) (string, []string, bool) {
	nodeToQuery, prefixToComplete, basePath, isCompletingCommand, err := c.parseInputForCompletion(in, cursorPos, c.root, cwd)
	if err != nil {
		return "", nil, false
	}
	prefix := prefixToComplete
	if nodeToQuery == nil {
		return prefix, nil, false
	}

	var rawSuggestions []string
	if isCompletingCommand {
		rawSuggestions = nodeToQuery.SuggestionsFor_NEW(prefixToComplete)
		if !strings.Contains(prefixToComplete, interfaces.PathSeparator) && c.system != nil {
			coreSuggestions := c.system.SuggestionsFor_NEW(prefixToComplete)
			rawSuggestions = c.mergeSuggestions(rawSuggestions, coreSuggestions)
		}
	} else {
		rawSuggestions = nodeToQuery.SuggestionsFor_NEW(prefixToComplete)
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

		if suggestionNode := nodeToQuery.FindChildren(rawSuggestion); suggestionNode != nil {
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
		suggestions = c.deduplicateSuggestions(suggestions) // Rimuovi duplicati esatti
	}
	found := len(suggestions) > 0
	return prefix, suggestions, found
}

// parseInputForCompletion identifies the command node, base path, and completion prefix from user input and cursor position.
// It returns the node to query, the prefix to complete, the base path, whether completing a command, and an error if any.
// Validation errors arise for invalid cursor positions, nil root, or current working directory commands.
// Traverses input paths to locate and validate nodes, supporting both absolute and relative paths.
// Handles completion for commands and subcommands by analyzing the input and directory structure.
func (c *Suggestion) parseInputForCompletion(in string, cursorPos int, root interfaces.ICommand, cwd interfaces.ICommand) (interfaces.ICommand, string, string, bool, error) {
	if cursorPos < 0 || cursorPos > len(in) {
		return nil, "", "", false, errors.New("invalid cursor position")
	}
	isCompletingCommand := false
	relevantInput := in[:cursorPos]
	lastSpacePos := strings.LastIndex(relevantInput, " ")
	var currentSegment string
	var textBeforeSegment string
	if lastSpacePos == -1 {
		currentSegment = relevantInput
		textBeforeSegment = ""
		isCompletingCommand = true
	} else {
		currentSegment = relevantInput[lastSpacePos+1:]
		textBeforeSegment = strings.TrimSpace(relevantInput[:lastSpacePos])
		isCompletingCommand = textBeforeSegment == ""
	}
	prefixToComplete := currentSegment

	var baseNode interfaces.ICommand
	pathPart := currentSegment
	isAbsolute := strings.HasPrefix(pathPart, pathSeparator)
	if isAbsolute {
		if root == nil {
			return nil, "", "", isCompletingCommand, errors.New("root command is nil")
		}
		baseNode = root
		pathPart = strings.TrimPrefix(pathPart, pathSeparator)
	} else {
		if cwd == nil {
			return nil, "", "", isCompletingCommand, errors.New("cwd command is nil")
		}
		baseNode = cwd
	}

	var dirParts []string
	if parts := strings.Split(pathPart, pathSeparator); len(parts) > 0 {
		prefixToComplete = parts[len(parts)-1]
		dirParts = parts[:len(parts)-1]
	} else {
		prefixToComplete = ""
		dirParts = []string{}
	}

	nodeToQuery := baseNode
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
		foundNode := nodeToQuery.FindChildren(part)
		if foundNode == nil {
			return nil, "", "", isCompletingCommand, fmt.Errorf("path not found: %s", part)
		}
		if !foundNode.HasSubCommands() && len(dirParts) > len(traversedPathParts)+1 {
			return nil, "", "", isCompletingCommand, fmt.Errorf("cannot traverse into non-directory: %s", part)
		}
		nodeToQuery = foundNode
		traversedPathParts = append(traversedPathParts, part)
	}
	basePath := nodeToQuery.CommandPath()

	//	   if isAbsolute {
	//	       basePath = pathSeparator + strings.Join(traversedPathParts,pathSeparator)
	//	   } else {
	//	       // Need cwd path + traversed parts... complex. Let's rely on nodeToQuery.CommandPath()
	//	       basePath = nodeToQuery.CommandPath()
	//	   }

	return nodeToQuery, prefixToComplete, basePath, isCompletingCommand, nil
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
