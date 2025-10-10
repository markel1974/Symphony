package pixels

// Point represents a graphical Point with position, color, picture coordinates, intensity, precision, and end shape.
type Point struct {
	pos       Vector
	col       RGBA
	pic       Vector
	in        float64
	precision int
	endShape  EndShape
}
