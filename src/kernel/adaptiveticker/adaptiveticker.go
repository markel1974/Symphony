package adaptiveticker

import (
	"container/list"
	"log"
	"sync"
	"time"
)

// maxEventsPerSecond defines the maximum number of events allowed per second for rate limiting.
// adaptiveTickerMaxQueueLen specifies the maximum length of the queue for adaptive ticking mechanisms.
const (
	maxEventsPerSecond        = 10
	adaptiveTickerMaxQueueLen = 1024
)

// AdaptiveTicker manages scheduled events, handling their creation, execution, removal, and expiration efficiently.
type AdaptiveTicker struct {
	container   list.List
	timer       *time.Timer
	runDeadline int64
	minInterval int64
	ids         *Ids
	index       map[int64]*TimerGroupHandler

	lock     sync.Mutex
	messages chan iEvent
	quit     bool

	statsEpsSecond    int64
	statsExpiredTotal int64
	statsExpired      int64
}

// NewAdaptiveTicker initializes and returns a pointer to an AdaptiveTicker with default configurations and event loop setup.
func NewAdaptiveTicker() *AdaptiveTicker {
	a := &AdaptiveTicker{
		timer:       nil,
		runDeadline: 0,
		ids:         NewIds(8192),
		minInterval: 1000 / maxEventsPerSecond,
		index:       make(map[int64]*TimerGroupHandler),
		messages:    make(chan iEvent, adaptiveTickerMaxQueueLen),
		quit:        false,
	}
	a.container.Init()

	a.eventLoop()

	return a
}

// Create initializes a new TimerHandler, assigns it an ID, and posts a creation event to the message channel.
func (a *AdaptiveTicker) Create(target chan *TimerHandler, event interface{}, first int64, interval int64, count int64) int {
	var current = NewTimerHandler(target, event, first, interval, count, a.minInterval)

	a.lock.Lock()
	if !a.quit {
		a.ids.Set(current)
	}
	a.lock.Unlock()

	if current.id != UnknownId {
		//fmt.Println("create - trying go send message …", current.id, len(a.messages))
		//debug.PrintStack()
		e := newCreateEvent(current)
		e.PostEvent(a.messages)
		//fmt.Println("create - message sent", current.id, len(a.messages))
	}

	return current.id
}

// RemoveEntries deletes the specified IDs from the AdaptiveTicker. Returns true if all IDs were successfully removed, false otherwise.
func (a *AdaptiveTicker) RemoveEntries(tids []int) bool {
	var removed []int
	a.lock.Lock()
	if !a.quit {
		for _, tid := range tids {
			if a.ids.Unset(tid) {
				removed = append(removed, tid)
			}
		}
	}
	a.lock.Unlock()
	if len(removed) > 0 {
		var event = newRemoveEvent(removed)
		event.PostEvent(a.messages)
	}
	return len(tids) == len(removed)
}

// Quit gracefully stops the AdaptiveTicker by setting the quit flag and emitting a quit event to the message channel.
func (a *AdaptiveTicker) Quit() {
	var emit = false
	a.lock.Lock()
	if !a.quit {
		a.quit = true
		emit = true
	}
	a.lock.Unlock()

	if emit {
		var event = newQuitEvent()
		event.PostEvent(a.messages)
	}
}

// doExpire checks for expired timers and handles their processing, returning a slice of expired TimerHandler instances.
func (a *AdaptiveTicker) doExpire(now int64) []*TimerHandler {
	var expired []*TimerHandler

	a.runDeadline = 0

	var next *list.Element
	for e := a.container.Front(); e != nil; e = next {
		next = e.Next()
		group := e.Value.(*TimerGroupHandler)
		if group.deadline <= now {
			for _, current := range group.container {
				if !current.removed {
					expired = append(expired, current)
				}
			}

			delete(a.index, group.deadline)
			a.container.Remove(e)
		} else {
			break
		}
	}

	for _, current := range expired {
		if current.IsUsable() {
			a.doAdd(now, current)
		} else {
			a.lock.Lock()
			a.ids.Unset(current.id)
			a.lock.Unlock()
		}
	}

	return expired
}

// doHeadRun schedules the next timer event execution based on the current time and the nearest deadline in the container.
// It resets or creates a new timer as needed to trigger the expiration event at the appropriate time.
// If the container is empty or a run is already scheduled (runDeadline > 0), the method returns immediately.
func (a *AdaptiveTicker) doHeadRun(now int64) {
	if a.container.Len() == 0 {
		return
	}

	if a.runDeadline > 0 {
		return
	}

	var target = a.container.Front().Value.(*TimerGroupHandler)
	a.runDeadline = target.deadline

	interval := a.runDeadline - now
	if interval < a.minInterval {
		interval = a.minInterval
	}

	d := time.Duration(interval) * time.Millisecond

	if a.timer != nil {
		if started := a.timer.Reset(d); !started {
			a.timer = nil
		}
	}

	if a.timer == nil {
		a.timer = time.NewTimer(d)
		go func(t *time.Timer) {
			select {
			case <-t.C:
				var event = newExpireEvent()
				event.PostEvent(a.messages)
			}
		}(a.timer)
	}
}

// doAdd adds a TimerHandler to the appropriate TimerGroupHandler based on its deadline and the provided current time.
// It ensures the TimerHandler is organized into a sorted container for efficient processing of expiration events.
func (a *AdaptiveTicker) doAdd(now int64, event *TimerHandler) {
	event.Prepare(now)

	if group, ok := a.index[event.deadline]; ok {
		group.Add(event)
		return
	}

	var newGroup = NewGroupEvent(event)

	a.index[newGroup.deadline] = newGroup

	if a.container.Len() == 0 {
		a.container.PushFront(newGroup)
		a.runDeadline = 0
		return
	}

	var head = a.container.Front()
	if head == nil {
		return
	}
	var g = head.Value.(*TimerGroupHandler)
	if newGroup.deadline < g.deadline {
		if a.runDeadline > 0 {
			if newGroup.deadline < a.runDeadline {
				if diff := a.runDeadline - newGroup.deadline; diff > a.minInterval {
					a.runDeadline = 0
				}
			}
		}

		a.container.InsertBefore(newGroup, head)
		//fmt.Println("ADDING", event.id, event.interval, event.deadline, ". ITERATION", 0)
		return
	}
	//var count = 0
	for elem := head.Next(); elem != nil; elem = elem.Next() {
		//count ++
		g = elem.Value.(*TimerGroupHandler)
		if newGroup.deadline < g.deadline {
			a.container.InsertBefore(newGroup, elem)
			//fmt.Println("ADDING", event.id, event.interval, event.deadline, "+ ITERATION", count)
			return
		}
	}

	a.container.PushBack(newGroup)
	//fmt.Println("ADDING", event.id, event.interval, event.deadline, "$ ITERATION", count)
}

// doQuit stops the timer, resets the run deadline to 0, and clears the container.
func (a *AdaptiveTicker) doQuit() {
	if a.timer != nil {
		a.timer.Stop()
		a.timer = nil
	}
	a.runDeadline = 0
	a.container.Init()
}

// eventLoop continuously handles incoming events from the message channel and processes them based on their event type.
// It spawns a goroutine to manage event handling and terminates when a quit event is received.
func (a *AdaptiveTicker) eventLoop() {
	go func() {
		for {
			select {
			case msg := <-a.messages:
				switch msg.GetType() {
				case eventTypeCreate:
					a.createEventHandler(msg.(*createEvent))
				case eventTypeRemove:
					a.removeEventHandler()
				case eventTypeExpire:
					a.expireEventHandler()
				case eventTypeQuit:
					a.quitEventHandler()
					return
				}
			}
		}
	}()
}

// createEventHandler processes a createEvent by adding its handler and scheduling the next timer execution.
func (a *AdaptiveTicker) createEventHandler(event *createEvent) {
	var now = getEpochMs()
	a.doAdd(now, event.handler)
	a.doHeadRun(now)
}

// removeEventHandler is a method used to manage the removal of event handlers and update the ticker's state accordingly.
func (a *AdaptiveTicker) removeEventHandler() {
	var now = getEpochMs()
	a.doHeadRun(now)
}

// expireEventHandler processes expired timer events by determining expired handlers, rescheduling, and posting events.
func (a *AdaptiveTicker) expireEventHandler() {
	var now = getEpochMs()
	var expired = a.doExpire(now)
	a.doHeadRun(now)

	//a.computeEps(now, int64(len(expired)))

	for _, current := range expired {
		current.PostEvent()
	}
}

// quitEventHandler handles the quit event by performing cleanup operations and ensuring the ticker stops gracefully.
func (a *AdaptiveTicker) quitEventHandler() {
	//var now = getEpochMs()
	a.doQuit()
}

// computeEps updates internal statistics based on expired events per second and prints aggregated data when the second changes.
func (a *AdaptiveTicker) computeEps(now int64, expired int64) {
	var epsSecond = now / 1000
	if epsSecond == a.statsEpsSecond {
		a.statsExpiredTotal += expired
		a.statsExpired++
	} else {
		log.Println("expired", a.statsExpired, "total", a.statsExpiredTotal)
		a.statsEpsSecond = epsSecond
		a.statsExpired = 1
		a.statsExpiredTotal = expired
	}
}
