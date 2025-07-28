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

// Point represents a 2D coordinate with x and y float64 values.
type Point struct {
	x float64
	y float64
}

// NewPointFloat creates a new Point instance with specified x and y float64 coordinates.
func NewPointFloat(x float64, y float64) Point {
	return Point{
		x: x,
		y: y,
	}
}

// AddTo modifies the current point by adding the specified x and y values to its coordinates.
func (p *Point) AddTo(x float64, y float64) {
	p.x += x
	p.y += y
}

// AddToX adds the given value to the x coordinate of the Point.
func (p *Point) AddToX(x float64) {
	p.x += x
}

// AddToY adds the specified value to the y-coordinate of the Point.
func (p *Point) AddToY(y float64) {
	p.y += y
}

// MoveTo repositions the Point to the specified x and y coordinates.
func (p *Point) MoveTo(x float64, y float64) {
	p.x = x
	p.y = y
}

// MoveToX updates the x-coordinate of the Point to the specified value.
func (p *Point) MoveToX(x float64) {
	p.x = x
}

// MoveToY sets the y-coordinate of the Point to the specified value.
func (p *Point) MoveToY(y float64) {
	p.y = y
}

// GetX returns the x-coordinate of the Point.
func (p *Point) GetX() float64 {
	return p.x
}

// GetY returns the y-coordinate of the Point as a float64.
func (p *Point) GetY() float64 {
	return p.y
}
