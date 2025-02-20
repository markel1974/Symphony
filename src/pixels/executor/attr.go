package executor

// AttrFormat represents a slice of Attr, defining the format of attributes for a vertex or uniform in a shader.
type AttrFormat []Attr

// Size calculates and returns the total size in bytes required by the AttrFormat based on the sizes of its attributes.
func (af AttrFormat) Size() int {
	total := 0
	for _, attr := range af {
		total += attr.Type.Size()
	}
	return total
}

// Attr represents a vertex attribute in a shader, including its name and type.
type Attr struct {
	Name string
	Type AttrType
}

// AttrType is an enumerated type used to represent various attribute types in a graphics or shader context.
type AttrType int

// Int represents a single integer attribute type.
// Float represents a single floating-point attribute type.
// Vec2 represents a 2D vector attribute type.
// Vec3 represents a 3D vector attribute type.
// Vec4 represents a 4D vector attribute type.
// Mat2 represents a 2x2 matrix attribute type.
// Mat23 represents a 2x3 matrix attribute type.
// Mat24 represents a 2x4 matrix attribute type.
// Mat3 represents a 3x3 matrix attribute type.
// Mat32 represents a 3x2 matrix attribute type.
// Mat34 represents a 3x4 matrix attribute type.
// Mat4 represents a 4x4 matrix attribute type.
// Mat42 represents a 4x2 matrix attribute type.
// Mat43 represents a 4x3 matrix attribute type.
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

// Size returns the size in bytes of the AttrType based on its definition, panicking if the type is invalid.
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
