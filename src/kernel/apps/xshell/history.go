package xshell

import (
	"encoding/json"
	"os"
	"sync"
)

// historySaveLock is a mutex used to synchronize access to shared resources during history save and restore operations.
var historySaveLock sync.Mutex

// HistoryHandler is a structure to manage a history queue with optional autosave functionality and position tracking.
type HistoryHandler struct {
	Queue    []string
	queuePos int
	max      uint
	def      string
	enabled  bool
	autosave bool
}

// NewHistoryHandler initializes and returns a HistoryHandler with a maximum history size and optional autosave functionality.
func NewHistoryHandler(max uint, autosave bool) *HistoryHandler {
	h := &HistoryHandler{
		max:      max,
		enabled:  true,
		autosave: autosave,
	}
	h.Clear()

	if h.autosave {
		h.restore()
	}
	return h
}

// save writes the current state of the HistoryHandler instance to a file in JSON format using a mutex for concurrency safety.
func (h *HistoryHandler) save() {
	if out, err := json.Marshal(h); err == nil {
		historySaveLock.Lock()
		_ = os.WriteFile("history.json", out, 0644)
		historySaveLock.Unlock()
	}
}

// restore loads the history data from a JSON file and updates the handler's state.
func (h *HistoryHandler) restore() {
	historySaveLock.Lock()
	body, err := os.ReadFile("history.json")
	historySaveLock.Unlock()

	if err == nil {
		h.Clear()
		_ = json.Unmarshal(body, h)
		h.queuePos = len(h.Queue) - 1
		h.def = ""
	}
}

// Clear resets the history by clearing the queue, reinitializing it, resetting the queue position, and clearing the default value.
func (h *HistoryHandler) Clear() {
	h.Queue = nil
	h.Queue = append(h.Queue, "")
	h.queuePos = 0
	h.def = ""
}

// SetEnabled configures whether the history tracking is enabled or disabled for the HistoryHandler.
func (h *HistoryHandler) SetEnabled(enabled bool) {
	h.enabled = enabled
}

// AddToHistory appends a new entry to the history queue if history tracking is enabled, and handles queue size and autosave.
func (h *HistoryHandler) AddToHistory(data string) {
	if h.enabled {
		h.Queue[len(h.Queue)-1] = data

		h.Queue = append(h.Queue, "")
		if len(h.Queue) > int(h.max) {
			h.Queue = h.Queue[1:]
		}

		h.queuePos = len(h.Queue) - 1

		if h.autosave {
			h.save()
		}
	}
}

// getHistoryAtIndex returns the history item at the specified index if within valid range and history is enabled.
func (h *HistoryHandler) getHistoryAtIndex(idx int) string {
	var out string
	if h.enabled {
		if idx >= 0 && idx < len(h.Queue)-1 {
			out = h.Queue[idx]
		}
	}
	return out
}

// GetHistoryPrev returns the previous entry in the history if available and moves the queue position back by one.
// Returns an empty string and false if history is disabled or at the beginning of the queue.
func (h *HistoryHandler) GetHistoryPrev() (string, bool) {
	if !h.enabled {
		return "", false
	}

	if h.queuePos == 0 {
		return "", false
	}

	h.queuePos--
	data := h.getHistoryAtIndex(h.queuePos)

	return data, true
}

// GetHistory returns a slice of strings containing all historical entries excluding the most recent placeholder entry.
func (h *HistoryHandler) GetHistory() []string {
	var out []string
	l := len(h.Queue) - 1
	if l > 0 {
		for x := 0; x < l; x++ {
			out = append(out, h.Queue[x])
		}
	}
	return out
}

// GetHistoryAtPos retrieves the history item at the specified position if valid, returning the item and a boolean flag.
func (h *HistoryHandler) GetHistoryAtPos(pos int) (string, bool) {
	var out string
	found := false
	l := len(h.Queue) - 1
	if l > 0 && pos < l {
		out = h.Queue[pos]
		if len(out) > 0 {
			found = true
		}
	}
	return out, found
}

// GetHistoryNext retrieves the next item in the history queue if available, or returns the default value if at the end.
func (h *HistoryHandler) GetHistoryNext() (string, bool) {
	if !h.enabled {
		return "", false
	}

	var data string
	maxPos := len(h.Queue) - 1

	if h.queuePos == maxPos {
		return "", false
	}

	if h.queuePos < maxPos-1 {
		h.queuePos++
		data = h.getHistoryAtIndex(h.queuePos)
	} else {
		h.queuePos++
		data = h.def
	}

	return data, true
}

// SetDefault sets the default string value used by the HistoryHandler when navigating beyond the current history range.
func (h *HistoryHandler) SetDefault(def string) {
	h.def = def
}
