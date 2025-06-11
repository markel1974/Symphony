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

package adaptiveticker

// eventType represents the type of event in the system as an integer enumeration.
type eventType int

// eventTypeCreate represents the event type for creating an entity.
// eventTypeRemove represents the event type for removing an entity.
// eventTypeExpire represents the event type for an entity expiration event.
// eventTypeQuit represents the event type for quitting or terminating an entity.
const (
	eventTypeCreate eventType = iota
	eventTypeRemove eventType = iota
	eventTypeExpire eventType = iota
	eventTypeQuit   eventType = iota
)

// iEvent defines an interface for event types with methods to get the event type and post the event onto a channel.
type iEvent interface {
	GetType() eventType
	PostEvent(ch chan iEvent)
}

// createEvent represents an event used to create and schedule a new timer via an associated TimerHandler.
type createEvent struct {
	handler *TimerHandler
}

// newCreateEvent creates and returns a new createEvent instance initialized with the provided TimerHandler.
func newCreateEvent(handler *TimerHandler) *createEvent {
	return &createEvent{handler: handler}
}

// GetType returns the type of the event as eventTypeCreate, indicating the event represents a "create" operation.
func (b *createEvent) GetType() eventType {
	return eventTypeCreate
}

// PostEvent sends the current createEvent instance asynchronously to the provided channel implementing iEvent.
func (b *createEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}

// removeEvent represents an event that encapsulates a list of target IDs for removal operations.
type removeEvent struct {
	tids []int
}

// newRemoveEvent creates a new removeEvent instance with the specified list of timer IDs.
func newRemoveEvent(tids []int) *removeEvent {
	return &removeEvent{tids: tids}
}

// GetType returns the type of the current event, which is eventTypeRemove.
func (b *removeEvent) GetType() eventType {
	return eventTypeRemove
}

// PostEvent sends the removeEvent instance to the provided channel asynchronously in a separate goroutine.
func (b *removeEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}

// expireEvent represents an event type used to signal expiration in the event system.
// This type implements the iEvent interface, providing methods for identifying the event type and posting the event.
type expireEvent struct {
}

// newExpireEvent creates and returns a new instance of expireEvent.
func newExpireEvent() *expireEvent {
	return &expireEvent{}
}

// GetType returns the type of the event as eventTypeExpire. Used to identify the event as an expiration event.
func (b *expireEvent) GetType() eventType {
	return eventTypeExpire
}

// PostEvent sends the event to the provided channel in a separate goroutine to avoid blocking.
func (b *expireEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}

// quitEvent represents an event type used to signal a quit operation within the event handling system.
type quitEvent struct {
}

// newQuitEvent creates and returns a new instance of quitEvent.
func newQuitEvent() *quitEvent {
	return &quitEvent{}
}

// GetType returns the type of the quitEvent as eventTypeQuit.
func (b *quitEvent) GetType() eventType {
	return eventTypeQuit
}

// PostEvent sends the current event instance to the provided channel asynchronously using a new goroutine.
func (b *quitEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}
