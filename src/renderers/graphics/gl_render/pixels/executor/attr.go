package executor

// AttrFormat represents a collection of attributes defining the format of vertex or uniform data.
type AttrFormat []Attr

// Size calculates the total size (in bytes) of all attributes in the AttrFormat by summing their individual sizes.
func (af AttrFormat) Size() int {
	total := 0
	for _, attr := range af {
		total += attr.Type.Size()
	}
	return total
}

// Attr represents a vertex attribute with a name and type used in shader programs.
type Attr struct {
	Name string
	Type AttrType
}

// AttrType represents the type of an attribute, such as a scalar, vector, or matrix, in a graphics or shader context.
type AttrType int

// Int defines a single integer attribute type.
// Float defines a single floating-point attribute type.
// Vec2 defines a 2D vector attribute type.
// Vec3 defines a 3D vector attribute type.
// Vec4 defines a 4D vector attribute type.
// Mat2 defines a 2x2 matrix attribute type.
// Mat23 defines a 2x3 matrix attribute type.
// Mat24 defines a 2x4 matrix attribute type.
// Mat3 defines a 3x3 matrix attribute type.
// Mat32 defines a 3x2 matrix attribute type.
// Mat34 defines a 3x4 matrix attribute type.
// Mat4 defines a 4x4 matrix attribute type.
// Mat42 defines a 4x2 matrix attribute type.
// Mat43 defines a 4x3 matrix attribute type.
const (
	Int AttrType = iota
	Float
	Vec2
	Vec3
	Vec4
	Mat2
	Mat23
	Mat24
	Mat3
	Mat32
	Mat34
	Mat4
	Mat42
	Mat43
)

// Size returns the memory size in bytes for the given AttrType. Panics if the type is invalid.
func (at AttrType) Size() int {
	switch at {
	case Int:
		return 4
	case Float:
		return 4
	case Vec2:
		return 8 //2 * 4
	case Vec3:
		return 12 //3 * 4
	case Vec4:
		return 16 //4 * 4
	case Mat2:
		return 16 //2 * 2 * 4
	case Mat23:
		return 24 //2 * 3 * 4
	case Mat24:
		return 32 //2 * 4 * 4
	case Mat3:
		return 36 //3 * 3 * 4
	case Mat32:
		return 24 //3 * 2 * 4
	case Mat34:
		return 32 //3 * 4 * 4
	case Mat4:
		return 64 //4 * 4 * 4
	case Mat42:
		return 32 //4 * 2 * 4
	case Mat43:
		return 48 //4 * 3 * 4
	default:
		panic("size of vertex attribute type: invalid type")
	}
}
