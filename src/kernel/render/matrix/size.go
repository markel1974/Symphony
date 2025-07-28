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

package matrix

// Size represents a 2D dimension with width (w) and height (h) as float64 values.
type Size struct {
	w float64
	h float64
}

// NewSize creates and returns a new Size instance with the specified width (w) and height (h).
func NewSize(w float64, h float64) Size {
	return Size{w: w, h: h}
}

// GetWidth returns the width value of the Size instance as a float64.
func (s *Size) GetWidth() float64 {
	return s.w
}

// GetHeight returns the height value of the Size instance.
func (s *Size) GetHeight() float64 {
	return s.h
}
