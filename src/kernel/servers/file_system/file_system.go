package file_system

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
)

const (
	fsQueueLen = 1024
	fsQueueMax = fsQueueLen - 1
)

// pathSeparator is a constant string representing the character used to separate components in a path, typically "/".
const pathSeparator = "/"

// FileSystem provides utilities for command hierarchy traversal and completion suggestions.
type FileSystem struct {
	process     interfaces.IUserProcess
	pid         int
	root        interfaces.ICommand
	searchPaths []interfaces.ICommand
	cwd         interfaces.ICommand
	parser      *Parser
	messageChan chan interfaces.IMessage
	kRouter     interfaces.IKernelResponseRouter
	handlers    map[interfaces.MessageType]func(interfaces.IMessage)
}

// NewFileSystem initializes and returns a new FileSystem instance with the given root command and an empty search path list.
func NewFileSystem(root interfaces.ICommand, sp []interfaces.ICommand) *FileSystem {
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
		root:        root,
		cwd:         root,
		searchPaths: searchPath,
		parser:      NewParser(false, false, ""),
		messageChan: make(chan interfaces.IMessage, fsQueueLen),
		handlers:    make(map[interfaces.MessageType]func(interfaces.IMessage)),
	}

	fs.handlers[interfaces.MessageTypeFileSystemCWDGetRequest] = fs.handleCWDGet
	fs.handlers[interfaces.MessageTypeFileSystemSuggestionRequest] = fs.handleSuggestion
	fs.handlers[interfaces.MessageTypeFileSystemCWDSetRequest] = fs.handleCWDSet
	fs.handlers[interfaces.MessageTypeFileSystemCWDPathRequest] = fs.handleCWDPath
	fs.handlers[interfaces.MessageTypeFileSystemCWDDirectoryListingRequest] = fs.handleCWDDirectoryListing
	fs.handlers[interfaces.MessageTypeFileSystemFindRequest] = fs.handleFindRequest
	fs.handlers[interfaces.MessageTypeFileSystemHelpRequest] = fs.handleHelp
	fs.handlers[interfaces.MessageTypeNotifyProcessCreate] = fs.handleProcessCreate
	fs.handlers[interfaces.MessageTypeNotifyProcessForeground] = fs.handleProcessForeground
	fs.handlers[interfaces.MessageTypeNotifyProcessTerminate] = fs.handleProcessTerminate
	return fs
}

// Name returns the name of the Render object as a string.
func (c *FileSystem) Name() string {
	return "fs"
}

func (c *FileSystem) PID() int {
	return c.pid
}

// Process returns the process implementation adhering to the interfaces.IUserProcess interface.
func (c *FileSystem) Process() interfaces.IUserProcess {
	return c.process
}

// PID returns an identifier for the file system process. It always returns a fixed value of -2.
//func (c *FileSystem) PID() int {
//	return c.process.PID()
//}

// User returns the username associated with the file system.
//func (c *FileSystem) User() string {
//	return c.process.User()
//}

// Register registers the file system with the provided kRouter and returns a slice of message types that it handles.
func (c *FileSystem) Register() []interfaces.MessageType {
	var out []interfaces.MessageType
	for id := range c.handlers {
		out = append(out, id)
	}
	return out
}

// Setup initializes the FileSystem with the provided process and starts the event loop to handle messages.
func (c *FileSystem) Setup(router interfaces.IKernelResponseRouter, pid int, process interfaces.IUserProcess) error {
	c.kRouter = router
	c.pid = pid
	c.process = process
	b := make(chan bool)
	c.eventLoop(b)
	_ = <-b
	return nil
}

// PostMessage sends a message of type IMessage to the file system's message channel for further processing.
func (c *FileSystem) PostMessage(m interfaces.IMessage) {
	if len(c.messageChan) >= fsQueueMax {
		log.Printf("FS: message queue full, dropping message: %d", m.GetType())
		return
	}
	//m.SetDestination(c.PID())
	c.messageChan <- m
}

// CallCWDSet updates the current working directory to the specified path and returns true if the operation is successful.
func (c *FileSystem) handleCWDSet(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDSetRequest)
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
			mt.CreateResponse(c.PID(), false)
			c.kRouter.PostKernelResponse(mt.Source(), mt)
			return
		}
		c.cwd = cmd
		mt.CreateResponse(c.PID(), true)
		c.kRouter.PostKernelResponse(mt.Source(), mt)
		return
	}
	mt.CreateResponse(c.PID(), false)
	c.kRouter.PostKernelResponse(mt.Source(), mt)
}

// CallAddSearchPath adds a new ICommand instance to the searchPaths slice for fs resolution.
//func (c *FileSystem) CallAddSearchPath(kRouter interfaces.IRouter, sp interfaces.ICommand) {
//	if sp == nil {
//		return
//	}
//	c.searchPaths = append(c.searchPaths, sp)
//}

// CallCWDPath returns the command path of the current working directory command.
func (c *FileSystem) handleCWDPath(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDPathRequest)
	if !ok {
		return
	}
	mt.CreateResponse(c.PID(), c.cwd.Path())
	c.kRouter.PostKernelResponse(mt.Source(), mt)
}

// CallCWDDirectoryListing retrieves the directory listing of the current working directory as a slice of strings.
func (c *FileSystem) handleCWDDirectoryListing(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDDirectoryListingRequest)
	if !ok {
		return
	}
	var out []string
	for _, z := range c.cwd.DirectoryListing() {
		out = append(out, z) // z.Name())
	}
	mt.CreateResponse(c.PID(), out)
	c.kRouter.PostKernelResponse(mt.Source(), mt)
}

// CallFind parses and executes a given command line string, associating it with a process, and manages its lifecycle.
func (c *FileSystem) handleFindRequest(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemFindRequest)
	if !ok {
		return
	}
	el, err := c.parser.Parse(mt.Line())
	if err != nil {
		mt.CreateResponse(c.PID(), nil, nil, err)
		c.kRouter.PostKernelResponse(mt.Source(), mt)
		return
	}
	if len(el) == 0 {
		mt.CreateResponse(c.PID(), nil, nil, fmt.Errorf("invalid command: '%s'", mt.Line()))
		c.kRouter.PostKernelResponse(mt.Source(), mt)
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
		mt.CreateResponse(c.PID(), nil, nil, fmt.Errorf("unknown command: '%s'", name))
		c.kRouter.PostKernelResponse(mt.Source(), mt)
		return
	}
	//if len(args) > 0 {
	//	if sel, args, err = sel.Find(args); err != nil || sel == nil {
	//		return false, fmt.Errorf("invalid command: '%s %s'", name, strings.Join(args, " "))
	//	}
	//}
	mt.CreateResponse(c.PID(), sel, args, nil)
	c.kRouter.PostKernelResponse(mt.Source(), mt)
}

// CallHelp retrieves help information for a given command line string,
// resolving it within the current context and search paths.
func (c *FileSystem) handleHelp(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemHelpRequest)
	if !ok {
		return
	}

	//kRouter interfaces.IRouter, path string
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
		mt.CreateResponse(c.PID(), "", fmt.Errorf("invalid command %s", mt.Path()))
		c.kRouter.PostKernelResponse(mt.Source(), mt)
		return
	}
	mt.CreateResponse(c.PID(), sel.Help(), nil)
	c.kRouter.PostKernelResponse(mt.Source(), mt)
	return
}

// CallSuggestion generates command suggestions based on the provided input and current directory context.
// It returns the input prefix, a list of suggestions, and a boolean indicating if suggestions exist.
func (c *FileSystem) handleSuggestion(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemSuggestionRequest)
	if !ok {
		return
	}
	textBeforeSegment, nodeToQuery, prefixToComplete, basePath, isCompletingCommand, err := c.parseInput(c.cwd, mt.In(), mt.Cursor())
	if err != nil || nodeToQuery == nil {
		mt.CreateResponse(c.PID(), "", nil, false)
		c.kRouter.PostKernelResponse(mt.Source(), mt)
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
		mt.CreateResponse(c.PID(), prefix, nil, false)
		c.kRouter.PostKernelResponse(mt.Source(), mt)
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
	mt.CreateResponse(c.PID(), prefix, suggestions, len(suggestions) > 0)
	c.kRouter.PostKernelResponse(mt.Source(), mt)
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
func (c *FileSystem) handleProcessCreate(_ interfaces.IMessage) {
	//todo notify cwd (cwd non è del filesystem, ma è del processo stesso!)
}

// NotifyProcessTermination handles notifications related to the termination of a process within the file system context.
func (c *FileSystem) handleProcessTerminate(_ interfaces.IMessage) {
}

// NotifyProcessForeground handles notifications related to the foregrounding of a process within the file system context.
func (c *FileSystem) handleProcessForeground(_ interfaces.IMessage) {
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
func (c *FileSystem) handleCWDGet(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemCWDGetRequest)
	if !ok {
		return
	}
	mt.CreateResponse(c.PID(), c.cwd.Name())
	c.kRouter.PostKernelResponse(mt.Source(), mt)
}
