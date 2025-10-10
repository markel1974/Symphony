package pixels

import (
	"fmt"
	"math"
)

// Vector represents a 2D vector with float64 components X and Y.
type Vector struct {
	X float64
	Y float64
}

// ZeroVector represents the zero vector (0, 0) of type Vector.
var ZeroVector = Vector{0, 0}

// NewVec creates a new Vector with the specified x and y components.
func NewVec(x float64, y float64) Vector {
	return Vector{x, y}
}

// Eq checks whether two vectors are approximately equal by comparing their components with a tolerance for rounding errors.
func (u Vector) Eq(v Vector) bool {
	return nearlyEqual(u.X, v.X) && nearlyEqual(u.Y, v.Y)
}

// Unit returns a unit vector rotated by the specified angle in radians.
func Unit(angle float64) Vector {
	return Vector{1, 0}.Rotated(angle)
}

// String returns a string representation of the Vector in the format "Vector(X, Y)".
func (u Vector) String() string {
	return fmt.Sprintf("Vector(%v, %v)", u.X, u.Y)
}

// XY returns the X and Y components of the Vector as two float64 values.
func (u Vector) XY() (float64, float64) {
	return u.X, u.Y
}

// Add adds the components of two vectors and returns the resulting vector.
func (u Vector) Add(v Vector) Vector {
	return Vector{
		u.X + v.X,
		u.Y + v.Y,
	}
}

// Sub subtracts the components of vector v from the components of vector u and returns the resulting vector.
func (u Vector) Sub(v Vector) Vector {
	return Vector{
		u.X - v.X,
		u.Y - v.Y,
	}
}

// Floor returns a new Vector with both components rounded down to the nearest integer using math.Floor.
func (u Vector) Floor() Vector {
	return Vector{
		math.Floor(u.X),
		math.Floor(u.Y),
	}
}

// To returns a new Vector that represents the difference between the given Vector v and the current Vector u.
func (u Vector) To(v Vector) Vector {
	return Vector{
		v.X - u.X,
		v.Y - u.Y,
	}
}

// Scaled returns a new vector by scaling the current vector by a given constant multiplier.
func (u Vector) Scaled(c float64) Vector {
	return Vector{u.X * c, u.Y * c}
}

// ScaledXY returns a new Vector resulting from element-wise multiplication of the current vector with another vector.
func (u Vector) ScaledXY(v Vector) Vector {
	return Vector{u.X * v.X, u.Y * v.Y}
}

// Len calculates and returns the Euclidean length (magnitude) of the vector.
func (u Vector) Len() float64 {
	return math.Hypot(u.X, u.Y)
}

// Angle computes and returns the angle of the vector in radians, measured counterclockwise from the positive X-axis.
func (u Vector) Angle() float64 {
	return math.Atan2(u.Y, u.X)
}

// Unit returns the unit (normalized) vector of the current vector. Defaults to Vector{1, 0} if the vector's length is zero.
func (u Vector) Unit() Vector {
	if u.X == 0 && u.Y == 0 {
		return Vector{1, 0}
	}
	return u.Scaled(1 / u.Len())
}

// Rotated returns a new Vector that is the result of rotating the original vector by the specified angle in radians.
func (u Vector) Rotated(angle float64) Vector {
	sin, cos := math.Sincos(angle)
	return Vector{
		u.X*cos - u.Y*sin,
		u.X*sin + u.Y*cos,
	}
}

// Normal returns a new Vector that is the 90-degree clockwise rotation (normal vector) of the current Vector.
func (u Vector) Normal() Vector {
	return Vector{-u.Y, u.X}
}

// Dot computes the dot product of the current vector and another vector, returning a scalar value.
func (u Vector) Dot(v Vector) float64 {
	return u.X*v.X + u.Y*v.Y
}

// Cross calculates the 2D cross product (determinant) of two vectors, returning a scalar result.
func (u Vector) Cross(v Vector) float64 {
	return u.X*v.Y - v.X*u.Y
}

// Project returns the projection of the vector u onto the vector v.
func (u Vector) Project(v Vector) Vector {
	length := u.Dot(v) / v.Len()
	return v.Unit().Scaled(length)
}

// Map applies the given function to each component of the vector and returns a new Vector with the results.
func (u Vector) Map(f func(float64) float64) Vector {
	return Vector{
		f(u.X),
		f(u.Y),
	}
}

// LinearInterpolation calculates a Point at a given ratio `t` between two vectors `a` and `b`.
// `t` is a value between 0 and 1, where 0 returns `a` and 1 returns `b`.
func LinearInterpolation(a Vector, b Vector, t float64) Vector {
	return a.Scaled(1 - t).Add(b.Scaled(t))
}

// nearlyEqual compares two float64 values for near equality within a small epsilon threshold.
// Returns true if the difference between the values is within the acceptable range.
func nearlyEqual(a float64, b float64) bool {
	epsilon := 0.000001
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	if a == 0.0 || b == 0.0 || diff < math.SmallestNonzeroFloat64 {
		return diff < (epsilon * math.SmallestNonzeroFloat64)
	}
	absA := math.Abs(a)
	absB := math.Abs(b)
	return diff/math.Min(absA+absB, math.MaxFloat64) < epsilon
}
