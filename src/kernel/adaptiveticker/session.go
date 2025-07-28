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

import "sync"

// Session represents a controlled lifecycle with acquire and release callbacks.
// It ensures thread-safe operations using a mutex and tracks a timestamp.
type Session struct {
	mutex     sync.Mutex
	now       int64
	acquireFn []func()
	releaseFn []func()
}

// NewSession creates and returns a new Session instance with provided acquire and release function slices.
func NewSession(acquireFn []func(), releaseFn []func()) *Session {
	return &Session{
		acquireFn: acquireFn,
		releaseFn: releaseFn,
	}
}

// Acquire locks the session, updates the current timestamp, and executes all registered acquisition functions, if any.
func (s *Session) Acquire() {
	s.mutex.Lock()
	s.now = getEpochMs()
	if s.acquireFn != nil {
		for _, fn := range s.acquireFn {
			fn()
		}
	}
}

// Release unlocks the Session's mutex and executes all registered release functions in the releaseFn slice, if any exist.
func (s *Session) Release() {
	if s.releaseFn != nil {
		for _, fn := range s.releaseFn {
			fn()
		}
	}
	s.mutex.Unlock()
}

// Now returns the current epoch time in milliseconds stored in the session.
func (s *Session) Now() int64 {
	return s.now
}
