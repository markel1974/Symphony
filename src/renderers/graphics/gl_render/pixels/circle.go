package pixels

import (
	"fmt"
	"math"
)

// Circle represents a geometric circle defined by a center point and a radius.
type Circle struct {
	Center Vector
	Radius float64
}

// NewCircle creates a new Circle with the specified center point and radius.
func NewCircle(center Vector, radius float64) Circle {
	return Circle{
		Center: center,
		Radius: radius,
	}
}

// String returns a string representation of the Circle in the format "Circle(Center, Radius)".
func (c Circle) String() string {
	return fmt.Sprintf("Circle(%s, %.2f)", c.Center, c.Radius)
}

// Norm returns a new Circle with the same center but ensures the radius is always a non-negative value.
func (c Circle) Norm() Circle {
	return Circle{
		Center: c.Center,
		Radius: math.Abs(c.Radius),
	}
}

// Area calculates and returns the area of the circle using the formula π * r^2, where r is the radius of the circle.
func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

// Moved returns a new Circle translated by the specified delta Vector without modifying the original Circle.
func (c Circle) Moved(delta Vector) Circle {
	return Circle{
		Center: c.Center.Add(delta),
		Radius: c.Radius,
	}
}

// Resized returns a new Circle with the same center and a radius incremented by the specified radiusDelta.
func (c Circle) Resized(radiusDelta float64) Circle {
	return Circle{
		Center: c.Center,
		Radius: c.Radius + radiusDelta,
	}
}

// Contains checks whether the given Vector u is located within or on the boundary of the Circle.
func (c Circle) Contains(u Vector) bool {
	toCenter := c.Center.To(u)
	return c.Radius >= toCenter.Len()
}

// Formula returns the X and Y coordinates of the Circle's center as two float64 values.
func (c Circle) Formula() (h, k float64) {
	return c.Center.X, c.Center.Y
}

// Union computes the smallest circle that can fully encompass two given circles. Returns the resulting enclosing circle.
func (c Circle) Union(d Circle) Circle {
	biggerC := maxCircle(c.Norm(), d.Norm())
	smallerC := minCircle(c.Norm(), d.Norm())

	dist := c.Center.To(d.Center).Len()

	if dist+smallerC.Radius <= biggerC.Radius {
		return biggerC
	}

	r := (dist + biggerC.Radius + smallerC.Radius) / 2

	theta := .5 + (biggerC.Radius-smallerC.Radius)/(2*dist)
	center := LinearInterpolation(smallerC.Center, biggerC.Center, theta)

	return Circle{
		Center: center,
		Radius: r,
	}
}

// Intersect calculates the intersection of two circles and returns the resulting Circle.
// If one circle fully encompasses the other, the larger circle is returned.
// If the circles do not overlap, a Circle with zero radius is returned at the calculated midpoint.
// Otherwise, it computes a Circle based on the overlapping region.
func (c Circle) Intersect(d Circle) Circle {
	// Check if one of the circles encompasses the other; if so, return that one
	biggerC := maxCircle(c.Norm(), d.Norm())
	smallerC := minCircle(c.Norm(), d.Norm())

	if biggerC.Radius >= biggerC.Center.To(smallerC.Center).Len()+smallerC.Radius {
		return biggerC
	}

	// Calculate the midpoint between the two radii
	// Distance between centers
	dist := c.Center.To(d.Center).Len()
	// Difference between radii
	diff := dist - (c.Radius + d.Radius)
	// Distance from c.Center to the weighted midpoint
	distToMidpoint := c.Radius + 0.5*diff
	// Weighted midpoint
	center := LinearInterpolation(c.Center, d.Center, distToMidpoint/dist)

	// No need to calculate radius if the circles do not overlap
	if c.Center.To(d.Center).Len() >= c.Radius+d.Radius {
		return NewCircle(center, 0)
	}

	radius := c.Center.To(d.Center).Len() - (c.Radius + d.Radius)

	return Circle{
		Center: center,
		Radius: math.Abs(radius),
	}
}

// IntersectLine calculates a vector adjustment needed to separate a line and a circle when they intersect.
func (c Circle) IntersectLine(l Line) Vector {
	return l.IntersectCircle(c).Scaled(-1)
}

// IntersectRect computes the intersection vector between a circle and a rectangle, indicating overlap or displacement.
func (c Circle) IntersectRect(r Rect) Vector {
	// Checks if the c.Center is not in the diagonal quadrants of the rectangle
	if (r.Min.X <= c.Center.X && c.Center.X <= r.Max.X) || (r.Min.Y <= c.Center.Y && c.Center.Y <= r.Max.Y) {
		// 'grow' the Rect by c.Radius in each orthogonal
		grown := Rect{Min: r.Min.Sub(NewVec(c.Radius, c.Radius)), Max: r.Max.Add(NewVec(c.Radius, c.Radius))}
		if !grown.Contains(c.Center) {
			// c.Center doesn't close enough to overlap, return zero-vector
			return ZeroVector
		}

		// Get minimum distance to travel out of Rect
		rToC := r.Center().To(c.Center)
		h := c.Radius - math.Abs(rToC.X) + (r.W() / 2)
		v := c.Radius - math.Abs(rToC.Y) + (r.H() / 2)

		if rToC.X < 0 {
			h = -h
		}
		if rToC.Y < 0 {
			v = -v
		}

		// No intersect
		if h == 0 && v == 0 {
			return ZeroVector
		}

		if math.Abs(h) > math.Abs(v) {
			// Vertical distance shorter
			return NewVec(0, v)
		}

		return NewVec(h, 0)
	} else {
		// The center is in the diagonal quadrants

		// Helper points to make code below easy to read.
		rectTopLeft := NewVec(r.Min.X, r.Max.Y)
		rectBottomRight := NewVec(r.Max.X, r.Min.Y)

		// Check for overlap.
		if !(c.Contains(r.Min) || c.Contains(r.Max) || c.Contains(rectTopLeft) || c.Contains(rectBottomRight)) {
			// No overlap.
			return ZeroVector
		}

		var centerToCorner Vector
		if c.Center.To(r.Min).Len() <= c.Radius {
			// Closest to bottom-left
			centerToCorner = c.Center.To(r.Min)
		}
		if c.Center.To(r.Max).Len() <= c.Radius {
			// Closest to top-right
			centerToCorner = c.Center.To(r.Max)
		}
		if c.Center.To(rectTopLeft).Len() <= c.Radius {
			// Closest to top-left
			centerToCorner = c.Center.To(rectTopLeft)
		}
		if c.Center.To(rectBottomRight).Len() <= c.Radius {
			// Closest to bottom-right
			centerToCorner = c.Center.To(rectBottomRight)
		}

		cornerToCircumferenceLen := c.Radius - centerToCorner.Len()

		return centerToCorner.Unit().Scaled(cornerToCircumferenceLen)
	}
}

// IntersectionPoints calculates the points where a line intersects with the circle.
// Returns an empty slice if no intersection exists, a slice with one point if tangent, and two points if it fully intersects.
func (c Circle) IntersectionPoints(l Line) []Vector {
	cContainsA := c.Contains(l.A)
	cContainsB := c.Contains(l.B)

	// Special case for both endpoints being contained within the circle
	if cContainsA && cContainsB {
		return []Vector{}
	}

	// Get the closest Point on the line to this circles' center
	closestToCenter := l.Closest(c.Center)

	// If the distance to the closest Point is greater than the radius, there are no points of intersection
	if closestToCenter.To(c.Center).Len() > c.Radius {
		return []Vector{}
	}

	// If the distance to the closest Point is equal to the radius, the line is tangent and the closest Point is the
	// Point at which it touches the circle.
	if closestToCenter.To(c.Center).Len() == c.Radius {
		return []Vector{closestToCenter}
	}

	// Special case for endpoint being on the circles' center
	if c.Center == l.A || c.Center == l.B {
		otherEnd := l.B
		if c.Center == l.B {
			otherEnd = l.A
		}
		intersect := c.Center.Add(c.Center.To(otherEnd).Unit().Scaled(c.Radius))
		return []Vector{intersect}
	}

	// This means the distance to the closest Point is less than the radius, so there is at least one intersection,
	// possibly two.

	// If one of the end points exists within the circle, there is only one intersection
	if cContainsA || cContainsB {
		containedPoint := l.A
		otherEnd := l.B
		if cContainsB {
			containedPoint = l.B
			otherEnd = l.A
		}

		// Use trigonometry to get the length of the line between the contained Point and the intersection Point.
		// The following is used to describe the triangle formed:
		//  A) Is the side between contained Point and circle center.
		//  B) is the side between the center and the intersection Point (radius).
		//  C) Is the side between the contained Point and the intersection Point.
		// The capitals of these letters are used as the angles opposite the respective sides.
		// A and b are known
		a := containedPoint.To(c.Center).Len()
		b := c.Radius
		// B can be calculated by subtracting the angle of b (to the x-axis) from the angle of c (to the x-axis)
		B := containedPoint.To(c.Center).Angle() - containedPoint.To(otherEnd).Angle()
		// Using the Sin rule, we can get A
		A := math.Asin((a * math.Sin(B)) / b)
		// Using the rule that there are 180 degrees (or Pi radians) in a triangle, we can now get C
		C := math.Pi - A + B
		// If C is zero, the line segment is in-line with the center-intersect line.
		var cz float64
		if C == 0 {
			cz = b - a
		} else {
			// Using the Sine rule again, we can now get cz
			cz = (a * math.Sin(C)) / math.Sin(A)
		}
		// Traveling from the contained Point to the other end by length of a will provide the intersection Point.
		return []Vector{
			containedPoint.Add(containedPoint.To(otherEnd).Unit().Scaled(cz)),
		}
	}

	// Otherwise, the endpoints exist outside the circle, and the line segment intersects in two locations.
	// The vector formed by going from the closest Point to the center of the circle will be perpendicular to the line;
	// this forms a right-angled triangle with the intersection points, with the radius as the hypotenuse.
	// Calculate the other triangles' sides' length.
	a := math.Sqrt(math.Pow(c.Radius, 2) - math.Pow(closestToCenter.To(c.Center).Len(), 2))

	// Traveling in both directions from the closest Point by length of a will provide the two intersection points.
	first := closestToCenter.Add(closestToCenter.To(l.A).Unit().Scaled(a))
	second := closestToCenter.Add(closestToCenter.To(l.B).Unit().Scaled(a))

	if first.To(l.A).Len() < second.To(l.A).Len() {
		return []Vector{first, second}
	}
	return []Vector{second, first}
}

// maxCircle returns the Circle with the larger radius between two given Circles. If radii are equal, the first Circle is returned.
func maxCircle(c Circle, d Circle) Circle {
	if c.Radius < d.Radius {
		return d
	}
	return c
}

// minCircle returns the smallest Circle by radius between two given Circle instances.
// If both Circles have the same radius, the first Circle is returned.
func minCircle(c Circle, d Circle) Circle {
	if c.Radius < d.Radius {
		return c
	}
	return d
}
