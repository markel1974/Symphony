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

package kernel

// MessageType represents an enumerated type used to classify different types of messages in the system.
type MessageType int

// MessageTypeRead represents a message type for read operations.
// MessageTypeTimer represents a message type for timer events.
// MessageTypePaint represents a message type for paint or render events.
// MessageTypeQuit represents a message type for quit or shutdown events.
const (
	MessageTypeRead  MessageType = iota
	MessageTypeTimer MessageType = iota
	MessageTypePaint MessageType = iota
	MessageTypeQuit  MessageType = iota
)

// iMessage is an interface representing a message with a type and the ability to post itself to a channel.
// getType retrieves the MessageType of the implementing message.
// postEvent posts the message to the provided channel of type iMessage.
type iMessage interface {
	getType() MessageType
	postEvent(chan iMessage)
}

// MessageRead represents a message containing a slice of byte data, typically used in read operations.
// Implements the iMessage interface to provide message type identification and event posting functionality.
type MessageRead struct {
	data []byte
}

// newMessageRead creates and returns a new MessageRead instance with a specified slice of data up to n bytes.
func newMessageRead(data []byte, n int) iMessage {
	if n > len(data) {
		n = len(data) - 1
	}
	x := data[:n]
	return &MessageRead{data: x}
}

// getType returns the message type associated with the MessageRead instance, which is MessageTypeRead.
func (m *MessageRead) getType() MessageType {
	return MessageTypeRead
}

// postEvent sends the MessageRead instance to the provided channel asynchronously using a goroutine.
func (m *MessageRead) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}

// MessageTimer represents a timer-based message with associated process ID, timer ID, and interval duration.
type MessageTimer struct {
	pid      int
	tid      int
	interval int
}

// newMessageTimer creates and initializes a new MessageTimer with the specified process ID and interval.
func newMessageTimer(pid int, interval int) *MessageTimer {
	return &MessageTimer{pid: pid, interval: interval}
}

// getType returns the type of the message, specifically MessageTypeTimer, for MessageTimer instances.
func (m *MessageTimer) getType() MessageType {
	return MessageTypeTimer
}

// postEvent sends the current MessageTimer instance to the provided channel implementing the iMessage interface asynchronously.
func (m *MessageTimer) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}

// MessageQuit represents a message used to signal the termination of an event loop or process.
// It implements the iMessage interface to support message handling and processing.
type MessageQuit struct {
}

// newMessageQuit creates and returns a pointer to a new instance of MessageQuit.
func newMessageQuit() *MessageQuit {
	return &MessageQuit{}
}

// getType returns the MessageType associated with the MessageQuit struct, which is MessageTypeQuit.
func (m *MessageQuit) getType() MessageType {
	return MessageTypeQuit
}

// postEvent sends the MessageQuit instance to the provided iMessage channel asynchronously.
func (m *MessageQuit) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}

// MessagePaint represents a message triggering a paint event, commonly used for rendering or visual updates in a system.
type MessagePaint struct {
}

// newMessagePaint creates and returns a new instance of MessagePaint.
func newMessagePaint() *MessagePaint {
	return &MessagePaint{}
}

// getType returns the type of the message, which is MessageTypePaint.
func (m *MessagePaint) getType() MessageType {
	return MessageTypePaint
}

// postEvent sends the current MessagePaint instance to the provided channel using a goroutine.
func (m *MessagePaint) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}
