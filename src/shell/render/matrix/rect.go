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

// Rect represents a 2D rectangle defined by its origin point, size, center point, z-depth, and an optional AABB.
type Rect struct {
	point  Point
	size   Size
	center Point
	z      float64
	aabb   *AABB
}

// NewRect creates and returns a new Rect instance with the specified position (x, y), dimensions (w, h) and depth (z).
func NewRect(x float64, y float64, w float64, h float64, z float64) Rect {
	r := Rect{
		aabb:  &AABB{},
		point: NewPointFloat(x, y),
		size:  NewSize(w, h),
		z:     z,
	}
	r.rebuild()

	return r
}

// rebuild recalculates the center, AABB boundaries, and surface area of the Rect based on its current position and size.
func (r *Rect) rebuild() {
	r.center.x = r.point.x + (r.size.w / 2)
	r.center.y = r.point.y + (r.size.h / 2)

	r.aabb.minX = r.point.x
	r.aabb.maxX = r.point.x + r.size.w
	r.aabb.minY = r.point.y
	r.aabb.maxY = r.point.y + r.size.h
	r.aabb.minZ = 0
	r.aabb.maxZ = r.z
	r.aabb.surfaceArea = r.aabb.calculateSurfaceArea()
}

// SetSize adjusts the width and height of the rectangle by the specified values and updates its internal properties.
func (r *Rect) SetSize(w float64, h float64) {
	r.size.w += w
	r.size.h += h
	r.rebuild()
}

// AddTo adjusts the position of the rectangle by adding the specified x and y offsets to its top-left corner coordinates.
func (r *Rect) AddTo(x float64, y float64) {
	r.point.x += x
	r.point.y += y
	r.rebuild()
}

// AddToX increments the x-coordinate of the Rect's point by the specified value and updates related properties.
func (r *Rect) AddToX(x float64) {
	r.point.x += x
	r.rebuild()
}

// AddToY increments the y-coordinate of the Rect's point by the specified value and updates its associated properties.
func (r *Rect) AddToY(y float64) {
	r.point.y += y
	r.rebuild()
}

// MoveTo repositions the Rect to the specified x and y coordinates and updates its properties.
func (r *Rect) MoveTo(x float64, y float64) {
	r.point.x = x
	r.point.y = y
	r.rebuild()
}

// MoveToX updates the x-coordinate of the Rect's top-left corner and recalculates its derived properties.
func (r *Rect) MoveToX(x float64) {
	r.point.x = x
	r.rebuild()
}

// MoveToY sets the y-coordinate of the Rect's position to the specified value and updates its internal state.
func (r *Rect) MoveToY(y float64) {
	r.point.y = y
	r.rebuild()
}

// GetX returns the x-coordinate of the rectangle's top-left corner.
func (r *Rect) GetX() float64 {
	return r.point.x
}

// GetY retrieves the y-coordinate of the Rect's top-left corner.
func (r *Rect) GetY() float64 {
	return r.point.y
}

// GetWidth returns the width of the Rect as a float64.
func (r *Rect) GetWidth() float64 {
	return r.size.w
}

// GetHeight returns the height of the Rect.
func (r *Rect) GetHeight() float64 {
	return r.size.h
}

// GetAABB returns the axis-aligned bounding box (AABB) associated with the Rect object.
func (r *Rect) GetAABB() *AABB {
	return r.aabb
}
