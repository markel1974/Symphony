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

// TimerGroupHandler represents a collection of TimerHandler instances grouped by a shared deadline for efficient processing.
type TimerGroupHandler struct {
	deadline  int64
	container []*TimerHandler
}

// NewGroupEvent creates a new TimerGroupHandler initialized with the given TimerHandler and sets its deadline.
func NewGroupEvent(event *TimerHandler) *TimerGroupHandler {
	g := &TimerGroupHandler{
		deadline: event.deadline,
	}
	g.Add(event)
	return g
}

// Add appends a TimerHandler event to the container within the TimerGroupHandler.
func (g *TimerGroupHandler) Add(event *TimerHandler) {
	g.container = append(g.container, event)
}

func (g *TimerGroupHandler) Len() int {
	return len(g.container)
}

func (g *TimerGroupHandler) GetDeadline() int64 {
	return g.deadline
}

func (g *TimerGroupHandler) Remove(tid int) bool {
	for i, e := range g.container {
		if e.sourceId == tid {
			g.container = append(g.container[:i], g.container[i+1:]...)
			return true
		}
	}
	return false
}
