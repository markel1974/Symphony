package context

import (
	"fmt"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"sort"
	"strings"
)

// pathSeparator is a constant string representing the character used to separate components in a path, typically "/".
const pathSeparator = "/"

// CommandInteractor provides utilities for command hierarchy traversal and completion suggestions.
type CommandInteractor struct {
	root        interfaces.ICommand
	searchPaths []interfaces.ICommand
	cwd         interfaces.ICommand
	parser      *Parser
}

// NewCommandInteractor initializes and returns a new CommandInteractor instance with the given root command and an empty search path list.
func NewCommandInteractor(root interfaces.ICommand, sp []interfaces.ICommand) *CommandInteractor {
	var searchPath []interfaces.ICommand
	for _, k := range sp {
		if k != nil {
			searchPath = append(searchPath, k)
		}
	}
	if searchPath == nil {
		searchPath = []interfaces.ICommand{}
	}
	return &CommandInteractor{
		root:        root,
		cwd:         root,
		searchPaths: searchPath,
		parser:      NewParser(false, false, ""),
	}
}

// AddSearchPath adds a new ICommand instance to the searchPaths slice for interactor resolution.
func (c *CommandInteractor) AddSearchPath(sp interfaces.ICommand) {
	if sp == nil {
		return
	}
	c.searchPaths = append(c.searchPaths, sp)
}

// CWD retrieves the current working directory command interface.
func (c *CommandInteractor) CWD() interfaces.ICommand {
	return c.cwd
}

// CWDSet updates the current working directory to the specified path and returns true if the operation is successful.
func (c *CommandInteractor) CWDSet(arg string) bool {
	var path []string
	for _, part := range strings.Split(arg, interfaces.PathSeparator) {
		if len(part) > 0 {
			path = append(path, part)
		}
	}
	if cmd := c.cwd.Traverse(path); cmd != nil {
		if !cmd.HasSubCommands() {
			return false
		}
		c.cwd = cmd
		return true
	}
	return false
}

// Find parses and executes a given command line string, associating it with a task, and manages its lifecycle.
func (c *CommandInteractor) Find(line string) (interfaces.ICommand, []string, error) {
	el, err := c.parser.Parse(line)
	if err != nil {
		return nil, nil, err
	}
	if len(el) == 0 {
		return nil, nil, fmt.Errorf("invalid command: '%s'", line)
	}
	name := el[0]
	args := el[1:]
	cmd := c.cwd
	if interfaces.IsPathAbsolute(name) {
		cmd = c.root
	}
	dirPath := interfaces.PathToSegments(name)
	sel := cmd.Traverse(dirPath)
	if sel == nil {
		for _, n := range c.searchPaths {
			sel = n.Traverse(dirPath)
			if sel != nil {
				break
			}
		}
	}
	if sel == nil {
		return nil, nil, fmt.Errorf("unknown command: '%s'", name)
	}
	//if len(args) > 0 {
	//	if sel, args, err = sel.Find(args); err != nil || sel == nil {
	//		return false, fmt.Errorf("invalid command: '%s %s'", name, strings.Join(args, " "))
	//	}
	//}
	return sel, args, nil
}

// Help retrieves help information for a given command line string,
// resolving it within the current context and search paths.
func (c *CommandInteractor) Help(path string) (string, error) {
	var pathSegments []string
	absolute := interfaces.IsPathAbsolute(path)
	if !absolute {
		pathSegments = append(pathSegments, c.cwd.Path()...)
	}
	segments := interfaces.PathToSegments(path)
	pathSegments = append(pathSegments, segments...)

	sel := c.root.Traverse(pathSegments)
	if sel == nil {
		for _, searchRoot := range c.searchPaths {
			node := searchRoot.Traverse(segments)
			if node != nil {
				sel = node
				break
			}
		}
	}
	if sel == nil {
		return "", fmt.Errorf("invalid command %s", path)
	}
	return sel.Help(), nil
}

// Suggestion generates command suggestions based on the provided input and current directory context.
// It returns the input prefix, a list of suggestions, and a boolean indicating if suggestions exist.
func (c *CommandInteractor) Suggestion(in string, cursor int) (string, []string, bool) {
	textBeforeSegment, nodeToQuery, prefixToComplete, basePath, isCompletingCommand, err := c.parseInput(c.cwd, in, cursor)
	if err != nil || nodeToQuery == nil {
		return "", nil, false
	}
	prefix := prefixToComplete

	rawSuggestions := nodeToQuery.SuggestionsFor_NEW(prefixToComplete)
	isFromSearchPath := false

	if len(rawSuggestions) == 0 && isCompletingCommand {
		for _, s := range c.searchPaths {
			if pathSuggestions := s.SuggestionsFor_NEW(prefixToComplete); len(pathSuggestions) > 0 {
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
func (c *CommandInteractor) parseInput(cwd interfaces.ICommand, in string, cursor int) (string, interfaces.ICommand, string, string, bool, error) {
	var input string
	if cursor < len(in) {
		input = in[:cursor]
	} else {
		input = in
	}

	isCompletingCommand := false
	pathPart := ""
	textBeforeSegment := ""

	//el, err := c.parse.Parse(input)
	//if err != nil {
	//	return "", nil, "", "", false, err
	//}
	//if len(el) == 0 {
	//	isCompletingCommand = true
	//} else if len(el) == 1 {
	//	pathPart = el[0]
	//	isCompletingCommand = true
	//} else if len(el) > 1 {
	//	textBeforeSegment = el[0]
	//	pathPart = el[1]
	//	isCompletingCommand = len(strings.TrimSpace(textBeforeSegment)) == 0
	//}

	//TODO BETTER IMPLEMENTATION.....
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
	isAbsolute := interfaces.IsPathAbsolute(pathPart)
	if isAbsolute {
		baseNode = c.root
		pathPart = strings.TrimPrefix(pathPart, pathSeparator)
	}

	pathSegments := interfaces.PathToSegments(pathPart)

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
func (c *CommandInteractor) mergeSuggestions(s1 []string, s2 []string) []string {
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
func (c *CommandInteractor) deduplicateSuggestions(s []string) []string {
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
