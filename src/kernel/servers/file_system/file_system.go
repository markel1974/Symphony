package file_system

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
)

// pathSeparator is a constant string representing the character used to separate components in a path, typically "/".
const pathSeparator = "/"

// FileSystem provides utilities for command hierarchy traversal and completion suggestions.
type FileSystem struct {
	user        string
	root        interfaces.ICommand
	searchPaths []interfaces.ICommand
	cwd         interfaces.ICommand
	parser      *Parser
	messageChan chan interfaces.IMessage
	router      interfaces.IRouter
	handlers    map[interfaces.MessageType]func(interfaces.IMessage)
}

// NewFileSystem initializes and returns a new FileSystem instance with the given root command and an empty search path list.
func NewFileSystem(user string, root interfaces.ICommand, sp []interfaces.ICommand) *FileSystem {
	var searchPath []interfaces.ICommand
	for _, k := range sp {
		if k != nil {
			searchPath = append(searchPath, k)
		}
	}
	if searchPath == nil {
		searchPath = []interfaces.ICommand{}
	}
	fs := &FileSystem{
		user:        user,
		root:        root,
		cwd:         root,
		searchPaths: searchPath,
		parser:      NewParser(false, false, ""),
		messageChan: make(chan interfaces.IMessage, 128),
		handlers:    make(map[interfaces.MessageType]func(interfaces.IMessage)),
	}

	fs.handlers[interfaces.MessageTypeFileSystemCWDGet] = fs.handleCWDGetRequest
	fs.handlers[interfaces.MessageTypeFileSystemSuggestion] = fs.handleSuggestion
	fs.handlers[interfaces.MessageTypeFileSystemCWDSet] = fs.handleCWDSet
	fs.handlers[interfaces.MessageTypeFileSystemCWDPath] = fs.handleCWDPath
	fs.handlers[interfaces.MessageTypeFileSystemCWDDirectoryListing] = fs.handleCWDDirectoryListing
	fs.handlers[interfaces.MessageTypeFileSystemFindRequest] = fs.handleFindRequest
	fs.handlers[interfaces.MessageTypeFileSystemHelp] = fs.handleHelp
	return fs
}

// Process returns nil, as the file system does not have a process.
func (c *FileSystem) Process() interfaces.IProcess {
	return nil
}

// PID returns an identifier for the file system process. It always returns a fixed value of -2.
func (c *FileSystem) PID() int {
	return -2
}

func (c *FileSystem) User() string {
	return c.user
}

func (c *FileSystem) Register(router interfaces.IRouter) []interfaces.MessageType {
	c.router = router
	var out []interfaces.MessageType
	for id := range c.handlers {
		out = append(out, id)
	}
	return out
}

// Start begins the process by setting its state to running and initiating its event loop asynchronously.
func (c *FileSystem) Start() {
	b := make(chan bool)
	c.eventLoop(b)
	_ = <-b
}

// PostMessage sends a message of type IMessage to the file system's message channel for further processing.
func (c *FileSystem) PostMessage(m interfaces.IMessage) {
	c.messageChan <- m
}

// CallCWDSet updates the current working directory to the specified path and returns true if the operation is successful.
func (c *FileSystem) handleCWDSet(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDSet)
	if !ok {
		return
	}
	var path []string
	for _, part := range strings.Split(mt.Path(), interfaces.PathSeparator) {
		if len(part) > 0 {
			path = append(path, part)
		}
	}
	if cmd := c.cwd.Traverse(path); cmd != nil {
		if cmd.Type() != interfaces.CommandTypeDirectory {
			mt.SetResult(false)
			c.router.PostMessage(mt)
			return
		}
		c.cwd = cmd
		mt.SetResult(true)
		c.router.PostMessage(mt)
		return
	}
	mt.SetResult(false)
	c.router.PostMessage(mt)
}

// CallAddSearchPath adds a new ICommand instance to the searchPaths slice for fs resolution.
//func (c *FileSystem) CallAddSearchPath(router interfaces.IRouter, sp interfaces.ICommand) {
//	if sp == nil {
//		return
//	}
//	c.searchPaths = append(c.searchPaths, sp)
//}

// CallCWDPath returns the command path of the current working directory command.
func (c *FileSystem) handleCWDPath(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDPath)
	if !ok {
		return
	}
	mt.SetResult(c.cwd.Path())
	c.router.PostMessage(mt)
}

// CallCWDDirectoryListing retrieves the directory listing of the current working directory as a slice of strings.
func (c *FileSystem) handleCWDDirectoryListing(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDDirectoryListing)
	if !ok {
		return
	}
	var out []string
	for _, z := range c.cwd.DirectoryListing() {
		out = append(out, z) // z.Name())
	}
	mt.SetResult(out)
	c.router.PostMessage(mt)
}

// CallFind parses and executes a given command line string, associating it with a process, and manages its lifecycle.
func (c *FileSystem) handleFindRequest(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemFindRequest)
	if !ok {
		return
	}
	el, err := c.parser.Parse(mt.Line())
	if err != nil {
		c.router.PostMessage(mt.CreateResponse(nil, nil, err))
		return
	}
	if len(el) == 0 {
		err = fmt.Errorf("invalid command: '%s'", mt.Line())
		c.router.PostMessage(mt.CreateResponse(nil, nil, err))
		return
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
			if sel = n.Traverse(dirPath); sel != nil {
				break
			}
		}
	}
	if sel == nil {
		err = fmt.Errorf("unknown command: '%s'", name)
		c.router.PostMessage(mt.CreateResponse(nil, nil, err))
		return
	}
	//if len(args) > 0 {
	//	if sel, args, err = sel.Find(args); err != nil || sel == nil {
	//		return false, fmt.Errorf("invalid command: '%s %s'", name, strings.Join(args, " "))
	//	}
	//}
	c.router.PostMessage(mt.CreateResponse(sel, args, nil))
}

// CallHelp retrieves help information for a given command line string,
// resolving it within the current context and search paths.
func (c *FileSystem) handleHelp(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemHelp)
	if !ok {
		return
	}

	//router interfaces.IRouter, path string
	//(string, error)
	var pathSegments []string
	absolute := interfaces.IsPathAbsolute(mt.Path())
	if !absolute {
		pathSegments = append(pathSegments, c.cwd.PathEntries()...)
	}
	segments := interfaces.PathToSegments(mt.Path())
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
		mt.SetResponse("", fmt.Errorf("invalid command %s", mt.Path()))
		mt.Router().PostMessage(mt)
		return
	}
	mt.SetResponse(sel.Help(), nil)
	mt.Router().PostMessage(mt)
	return
}

// CallSuggestion generates command suggestions based on the provided input and current directory context.
// It returns the input prefix, a list of suggestions, and a boolean indicating if suggestions exist.
func (c *FileSystem) handleSuggestion(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemSuggestion)
	if !ok {
		return
	}
	textBeforeSegment, nodeToQuery, prefixToComplete, basePath, isCompletingCommand, err := c.parseInput(c.cwd, mt.In(), mt.Cursor())
	if err != nil || nodeToQuery == nil {
		mt.SetResponse("", nil, false)
		c.router.PostMessage(mt)
		return
	}
	prefix := prefixToComplete

	rawSuggestions := nodeToQuery.SuggestionsFor(prefixToComplete)
	isFromSearchPath := false

	if len(rawSuggestions) == 0 && isCompletingCommand {
		for _, s := range c.searchPaths {
			if pathSuggestions := s.SuggestionsFor(prefixToComplete); len(pathSuggestions) > 0 {
				rawSuggestions = pathSuggestions
				isFromSearchPath = true
				break
			}
		}
	}

	if len(rawSuggestions) == 0 {
		mt.SetResponse(prefix, nil, false)
		c.router.PostMessage(mt)
		return
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
			if sNode := nodeToQuery.FindChildren(rawSuggestion); sNode != nil && sNode.Type() == interfaces.CommandTypeDirectory {
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
	mt.SetResponse(prefix, suggestions, len(suggestions) > 0)
	c.router.PostMessage(mt)
}

// parseInput parses the input string to determine the relevant command context, path, and completion prefix details.
// It returns the text before the path segment, the current command node, the prefix for completion, the base path,
// a boolean indicating if the input addresses a command name, and any error encountered during processing.
func (c *FileSystem) parseInput(cwd interfaces.ICommand, in string, cursor int) (string, interfaces.ICommand, string, string, bool, error) {
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
		if foundNode.Type() != interfaces.CommandTypeDirectory && len(dirParts) > len(traversedPathParts)+1 {
			return "", nil, "", "", isCompletingCommand, fmt.Errorf("cannot traverse into non-directory: %s", part)
		}
		currentNode = foundNode
		traversedPathParts = append(traversedPathParts, part)
	}

	return textBeforeSegment, currentNode, prefixToComplete, basePath, isCompletingCommand, nil
}

// mergeSuggestions merges two slices of suggestions, removes duplicates, and returns the combined slice.
func (c *FileSystem) mergeSuggestions(s1 []string, s2 []string) []string {
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
func (c *FileSystem) deduplicateSuggestions(s []string) []string {
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

// NotifyProcessCreation handles notifications related to the creation of a process within the file system context.
func (c *FileSystem) NotifyProcessCreation(pid int, name string) {
	//todo notify cwd (cwd non è del filesystem, ma è del processo stesso!)
}

// NotifyProcessTermination handles notifications related to the termination of a process within the file system context.
func (c *FileSystem) NotifyProcessTermination(pid int) {
}

func (c *FileSystem) NotifyProcessForeground(pid int) {
}

// evenLoop continuously listens on the message channel and processes incoming messages until a quit message is received.
func (c *FileSystem) eventLoop(r chan bool) {
	go func() {
		r <- true
		for {
			select {
			case m, ok := <-c.messageChan:
				if !ok {
					return
				}
				id := m.GetType()
				if id == interfaces.MessageTypeQuit {
					close(c.messageChan)
					return
				}
				//fmt.Println("Executing id", id)
				if handler, _ := c.handlers[id]; handler != nil {
					handler(m)
				} else {
					log.Printf("FileSystem: unknown message type: %d", id)
				}
			}
		}
	}()
}

// handleCWDName returns the name of the current working directory command.
func (c *FileSystem) handleCWDGetRequest(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDGet)
	if !ok {
		return
	}
	mt.SetResult(c.cwd.Name())
	c.router.PostMessage(mt)
}
