package pixels

import (
	"fmt"
	"math"
)

// Line is a 2D line segment, between points A and B.
type Line struct {
	A Vector
	B Vector
}

// NewLine creates and returns a new Line.
func NewLine(from, to Vector) Line {
	return Line{
		A: from,
		B: to,
	}
}

// Bounds returns the lines bounding box.  This is in the form of a normalized Rect.
func (l Line) Bounds() Rect {
	return NewRect(l.A.X, l.A.Y, l.B.X, l.B.Y).Norm()
}

// Center will return the Point at the center of the line; that is, the Point equidistant from either end.
func (l Line) Center() Vector {
	return l.A.Add(l.A.To(l.B).Scaled(0.5))
}

// Closest will return the Point on the line which is closest to the Vector provided.
func (l Line) Closest(v Vector) Vector {
	// between is a helper function which determines whether x is greater than min(a, b) and less than max(a, b)
	between := func(a, b, x float64) bool {
		minV := math.Min(a, b)
		maxV := math.Max(a, b)
		return minV < x && x < maxV
	}

	// The closest Point will be on a line which perpendicular to this line.
	// If and only if the infinite perpendicular line intersects the segment.
	m, b := l.Formula()

	// Account for horizontal lines
	if m == 0 {
		x := v.X
		y := l.A.Y

		// check if the X coordinate of v is on the line
		if between(l.A.X, l.B.X, v.X) {
			return NewVec(x, y)
		}

		// Otherwise, get the closest endpoint
		if l.A.To(v).Len() < l.B.To(v).Len() {
			return l.A
		}
		return l.B
	}

	// Account for vertical lines
	if math.IsInf(math.Abs(m), 1) {
		x := l.A.X
		y := v.Y

		// check if the Y coordinate of v is on the line
		if between(l.A.Y, l.B.Y, v.Y) {
			return NewVec(x, y)
		}

		// Otherwise, get the closest endpoint
		if l.A.To(v).Len() < l.B.To(v).Len() {
			return l.A
		}
		return l.B
	}

	perpendicularM := -1 / m
	perpendicularB := v.Y - (perpendicularM * v.X)

	// Coordinates intersected (of infinite lines)
	x := (perpendicularB - b) / (m - perpendicularM)
	y := m*x + b

	// Check if the Point lies between the x and y bounds of the segment
	if !between(l.A.X, l.B.X, x) && !between(l.A.Y, l.B.Y, y) {
		// Not within bounding box
		toStart := v.To(l.A)
		toEnd := v.To(l.B)

		if toStart.Len() < toEnd.Len() {
			return l.A
		}
		return l.B
	}

	return NewVec(x, y)
}

// Contains returns whether the provided Vector lies on the line.
func (l Line) Contains(v Vector) bool {
	return l.Closest(v).Eq(v)
}

// Formula will return the values that represent the line in the formula: y = mx + b
// This function will return math.Inf+, math.Inf- for a vertical line.
func (l Line) Formula() (m, b float64) {
	// Account for horizontal lines
	if l.B.Y == l.A.Y {
		return 0, l.A.Y
	}

	m = (l.B.Y - l.A.Y) / (l.B.X - l.A.X)
	b = l.A.Y - (m * l.A.X)

	return m, b
}

// Intersect will return the Point of intersection for the two-line segments.
// If the line segments do not intersect, this function will return the zero-vector and false.
func (l Line) Intersect(k Line) (Vector, bool) {
	// Check if the lines are parallel
	lDir := l.A.To(l.B)
	kDir := k.A.To(k.B)
	if lDir.X == kDir.X && lDir.Y == kDir.Y {
		return ZeroVector, false
	}

	// The lines intersect - but potentially not within the line segments.
	// Get the intersection Point for the lines if they were infinitely long, check if the Point exists on both of the
	// segments
	lm, lb := l.Formula()
	km, kb := k.Formula()

	// Account for vertical lines
	if math.IsInf(math.Abs(lm), 1) && math.IsInf(math.Abs(km), 1) {
		// Both vertical, therefore, parallel
		return ZeroVector, false
	}

	var x, y float64

	if math.IsInf(math.Abs(lm), 1) || math.IsInf(math.Abs(km), 1) {
		// One line is vertical
		intersectM := lm
		intersectB := lb
		verticalLine := k

		if math.IsInf(math.Abs(lm), 1) {
			intersectM = km
			intersectB = kb
			verticalLine = l
		}

		y = intersectM*verticalLine.A.X + intersectB
		x = verticalLine.A.X
	} else {
		// Coordinates intersected
		x = (kb - lb) / (lm - km)
		y = lm*x + lb
	}

	if l.Contains(NewVec(x, y)) && k.Contains(NewVec(x, y)) {
		// The intersected Point is on both line segments, they intersect.
		return NewVec(x, y), true
	}

	return ZeroVector, false
}

// IntersectCircle will return the shortest Vector such that moving the Line by that Vector will cause the Line and Circle
// to no longer intersect.
// If they do not intersect at all, this function will return a zero-vector.
func (l Line) IntersectCircle(c Circle) Vector {
	// Get the Point on the line closest to the center of the circle.
	closest := l.Closest(c.Center)
	cirToClosest := c.Center.To(closest)

	if cirToClosest.Len() >= c.Radius {
		return ZeroVector
	}

	return cirToClosest.Scaled(cirToClosest.Len() - c.Radius)
}

// IntersectRect will return the shortest Vector such that moving the Line by that Vector will cause the Line and Rect to
// no longer intersect.
// If they do not intersect at all, this function will return a zero-vector.
func (l Line) IntersectRect(r Rect) Vector {
	// Check if either end of the line segment are within the rectangle
	if r.Contains(l.A) || r.Contains(l.B) {
		// Use the Rect.Intersect to get minimal return value
		rIntersect := l.Bounds().Intersect(r)
		if rIntersect.H() > rIntersect.W() {
			// Go vertical
			return NewVec(0, rIntersect.H())
		}
		return NewVec(rIntersect.W(), 0)
	}

	// Check if any of the rectangles' edges intersect with this line.
	for _, edge := range r.Edges() {
		if _, ok := l.Intersect(edge); ok {
			// Get the closest points on the line to each corner, where:
			//  - the Point is contained by the rectangle
			//  - the Point is not the corner itself
			corners := r.Vertices()
			var closest *Vector
			closestCorner := corners[0]
			for _, c := range corners {
				cc := l.Closest(c)
				if closest == nil || (closest.Len() > cc.Len() && r.Contains(cc)) {
					closest = &cc
					closestCorner = c
				}
			}
			if closest != nil {
				return closest.To(closestCorner)
			}
		}
	}

	// No intersect
	return ZeroVector
}

// Len returns the length of the line segment.
func (l Line) Len() float64 {
	return l.A.To(l.B).Len()
}

// Moved will return a line moved by the delta Vector provided.
func (l Line) Moved(delta Vector) Line {
	return Line{
		A: l.A.Add(delta),
		B: l.B.Add(delta),
	}
}

// Rotated will rotate the line around the provided Vector.
func (l Line) Rotated(around Vector, angle float64) Line {
	// Move the line so we can use `Vector.Rotated`
	lineShifted := l.Moved(around.Scaled(-1))

	lineRotated := Line{
		A: lineShifted.A.Rotated(angle),
		B: lineShifted.B.Rotated(angle),
	}

	return lineRotated.Moved(around)
}

// Scaled will return the line scaled around the center Point.
func (l Line) Scaled(scale float64) Line {
	return l.ScaledXY(l.Center(), scale)
}

// ScaledXY will return the line scaled around the Vector provided.
func (l Line) ScaledXY(around Vector, scale float64) Line {
	toA := around.To(l.A).Scaled(scale)
	toB := around.To(l.B).Scaled(scale)

	return Line{
		A: around.Add(toA),
		B: around.Add(toB),
	}
}

func (l Line) String() string {
	return fmt.Sprintf("Line(%v, %v)", l.A, l.B)
}
