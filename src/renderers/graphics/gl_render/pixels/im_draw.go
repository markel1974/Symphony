package pixels

import (
	"image/color"
	"math"
)

// EndShape represents the type for specifying the shape style at the end of a graphical path or line.
type EndShape int

// NoEndShape represents an absence of an end shape.
// SharpEndShape represents a sharp end shape.
// RoundEndShape represents a rounded end shape.
const (
	NoEndShape EndShape = iota
	SharpEndShape
	RoundEndShape
)

// IMDraw is a structure used for immediate-mode 2D drawing, such as shapes or points, with customizable attributes.
// The struct handles parameters like color, intensity, precision, and transformation matrices for rendering.
type IMDraw struct {
	Color     color.Color
	Picture   Vector
	Intensity float64
	Precision int
	EndShape  EndShape
	points    []Point
	pool      [][]Point
	matrix    Matrix
	mask      RGBA
	tri       *TrianglesData
	batch     *Batch
}

// _ is used to enforce that IMDraw implements the IBasicTarget interface at compile time.
var _ IBasicTarget = (*IMDraw)(nil)

// NewIMDraw initializes and returns a new instance of IMDraw using the provided IPicture. It sets up default values and state.
func NewIMDraw(pic IPicture) *IMDraw {
	tri := &TrianglesData{}
	im := &IMDraw{
		tri:   tri,
		batch: NewBatch(tri, pic),
	}
	im.SetMatrix(IM)
	im.SetColorMask(Alpha(1))
	im.Reset()
	return im
}

// Clear removes all points, shapes, and drawing data from the current IMDraw instance, resetting it for new drawing operations.
func (imd *IMDraw) Clear() {
	imd.tri.SetLen(0)
	imd.batch.Dirty()
}

// Reset resets the state of the IMDraw object to its default values, clearing points and restoring default properties.
func (imd *IMDraw) Reset() {
	imd.points = imd.points[:0]
	imd.Color = Alpha(1)
	imd.Picture = ZeroVector
	imd.Intensity = 0
	imd.Precision = 64
	imd.EndShape = NoEndShape
}

// Draw executes the rendering operations encapsulated in the IMDraw instance onto the specified target.
func (imd *IMDraw) Draw(t ITarget) {
	imd.batch.Draw(t)
}

// Push adds one or more points to the current drawing state with associated options like color, precision, and shape end.
func (imd *IMDraw) Push(pts ...Vector) {
	// Assert that Color is of type pixel.RGBA,
	if _, ok := imd.Color.(RGBA); !ok {
		// otherwise cast it
		imd.Color = ToRGBA(imd.Color)
	}
	opts := Point{
		col:       imd.Color.(RGBA),
		pic:       imd.Picture,
		in:        imd.Intensity,
		precision: imd.Precision,
		endShape:  imd.EndShape,
	}
	for _, pt := range pts {
		imd.pushPt(pt, opts)
	}
}

// pushPt appends a Point with the given position and attributes to the IMDraw points slice.
func (imd *IMDraw) pushPt(pos Vector, pt Point) {
	pt.pos = pos
	imd.points = append(imd.points, pt)
}

// SetMatrix sets the transformation matrix for drawing operations in the IMDraw instance.
func (imd *IMDraw) SetMatrix(m Matrix) {
	imd.matrix = m
	imd.batch.SetMatrix(imd.matrix)
}

// SetColorMask sets the color mask for rendering, converting the input color to RGBA and applying it to the batch.
func (imd *IMDraw) SetColorMask(color color.Color) {
	imd.mask = ToRGBA(color)
	imd.batch.SetColorMask(imd.mask)
}

// MakeTriangles creates a new ITargetTriangles by passing the given ITriangles to the batch's MakeTriangles method.
func (imd *IMDraw) MakeTriangles(t ITriangles) ITargetTriangles {
	return imd.batch.MakeTriangles(t)
}

// MakePicture creates an ITargetPicture from the provided IPicture. The resulting ITargetPicture can be drawn.
func (imd *IMDraw) MakePicture(p IPicture) ITargetPicture {
	return imd.batch.MakePicture(p)
}

// Line draws a line connecting the points added via Push with the specified thickness.
// If less than two points exist, no line will be drawn.
func (imd *IMDraw) Line(thickness float64) {
	imd.polyline(thickness, false)
}

// Rectangle draws a rectangle with the given thickness. A thickness of 0 fills the rectangle; otherwise, it outlines it.
func (imd *IMDraw) Rectangle(thickness float64) {
	if thickness == 0 {
		imd.fillRectangle()
	} else {
		imd.outlineRectangle(thickness)
	}
}

// Polygon draws a polygon using the specified thickness. A thickness of 0 fills the polygon; otherwise, it outlines it.
func (imd *IMDraw) Polygon(thickness float64) {
	if thickness == 0 {
		imd.fillPolygon()
	} else {
		imd.polyline(thickness, true)
	}
}

// Circle draws a circle with the specified radius and thickness.
// If thickness is zero, the circle is filled; otherwise, it is outlined.
func (imd *IMDraw) Circle(radius, thickness float64) {
	if thickness == 0 {
		imd.fillEllipseArc(NewVec(radius, radius), 0, 2*math.Pi)
	} else {
		imd.outlineEllipseArc(NewVec(radius, radius), 0, 2*math.Pi, thickness, false)
	}
}

// CircleArc draws a circular arc with the specified radius, angle range (low to high), and optional thickness.
// A thickness of 0 fills the arc, while a positive thickness outlines it.
func (imd *IMDraw) CircleArc(radius, low, high, thickness float64) {
	if thickness == 0 {
		imd.fillEllipseArc(NewVec(radius, radius), low, high)
	} else {
		imd.outlineEllipseArc(NewVec(radius, radius), low, high, thickness, true)
	}
}

// Ellipse draws an ellipse with the specified radius and thickness centered at the current position.
// If thickness is 0, the ellipse is filled; otherwise, it is outlined.
func (imd *IMDraw) Ellipse(radius Vector, thickness float64) {
	if thickness == 0 {
		imd.fillEllipseArc(radius, 0, 2*math.Pi)
	} else {
		imd.outlineEllipseArc(radius, 0, 2*math.Pi, thickness, false)
	}
}

// EllipseArc draws an elliptical arc with the specified radius, angular range (low to high), and thickness.
// If thickness is zero, the arc is filled; otherwise, it is outlined.
func (imd *IMDraw) EllipseArc(radius Vector, low, high, thickness float64) {
	if thickness == 0 {
		imd.fillEllipseArc(radius, low, high)
	} else {
		imd.outlineEllipseArc(radius, low, high, thickness, true)
	}
}

// getAndClearPoints retrieves and clears the current list of points, reusing allocated memory where possible to reduce reallocation.
func (imd *IMDraw) getAndClearPoints() []Point {
	points := imd.points
	// use one of the existing pools so we don't reallocate as often
	if len(imd.pool) > 0 {
		pos := len(imd.pool) - 1
		imd.points = imd.pool[pos][:0]
		imd.pool = imd.pool[:pos]
	} else {
		imd.points = nil
	}
	return points
}

// restorePoints moves the current set of points into the pool and sets the points slice to the provided empty slice.
func (imd *IMDraw) restorePoints(points []Point) {
	imd.pool = append(imd.pool, imd.points)
	imd.points = points[:0]
}

// applyMatrixAndMask applies the transformation matrix and color mask to the TrianglesData starting at the given offset.
func (imd *IMDraw) applyMatrixAndMask(off int) {
	for i := range (*imd.tri)[off:] {
		(*imd.tri)[off+i].Position = imd.matrix.Project((*imd.tri)[off+i].Position)
		(*imd.tri)[off+i].Color = imd.mask.Mul((*imd.tri)[off+i].Color)
	}
}

// fillRectangle creates and fills rectangles based on input points and color blending. Handles matrix and mask transformations.
func (imd *IMDraw) fillRectangle() {
	points := imd.getAndClearPoints()

	if len(points) < 2 {
		imd.restorePoints(points)
		return
	}

	off := imd.tri.Len()
	imd.tri.SetLen(imd.tri.Len() + 6*(len(points)-1))

	for i, j := 0, off; i+1 < len(points); i, j = i+1, j+6 {
		a, b := points[i], points[i+1]
		c := Point{
			pos: NewVec(a.pos.X, b.pos.Y),
			col: a.col.Add(b.col).Mul(Alpha(0.5)),
			pic: NewVec(a.pic.X, b.pic.Y),
			in:  (a.in + b.in) / 2,
		}
		d := Point{
			pos: NewVec(b.pos.X, a.pos.Y),
			col: a.col.Add(b.col).Mul(Alpha(0.5)),
			pic: NewVec(b.pic.X, a.pic.Y),
			in:  (a.in + b.in) / 2,
		}
		for k, p := range [...]Point{a, b, c, a, b, d} {
			(*imd.tri)[j+k].Position = p.pos
			(*imd.tri)[j+k].Color = p.col
			(*imd.tri)[j+k].Picture = p.pic
			(*imd.tri)[j+k].Intensity = p.in
		}
	}

	imd.applyMatrixAndMask(off)
	imd.batch.Dirty()

	imd.restorePoints(points)
}

// outlineRectangle draws an outlined rectangle using the given thickness.
// It processes the stored points, determines rectangle boundaries, and creates a polyline around them.
func (imd *IMDraw) outlineRectangle(thickness float64) {
	points := imd.getAndClearPoints()

	if len(points) < 2 {
		imd.restorePoints(points)
		return
	}

	for i := 0; i+1 < len(points); i++ {
		a, b := points[i], points[i+1]
		mid := a
		mid.col = a.col.Add(b.col).Mul(Alpha(0.5))
		mid.in = (a.in + b.in) / 2

		imd.pushPt(a.pos, a)
		imd.pushPt(NewVec(a.pos.X, b.pos.Y), mid)
		imd.pushPt(b.pos, b)
		imd.pushPt(NewVec(b.pos.X, a.pos.Y), mid)
		imd.polyline(thickness, true)
	}

	imd.restorePoints(points)
}

// fillPolygon fills a polygon from a set of points, assuming at least three points are provided to define the shape.
func (imd *IMDraw) fillPolygon() {
	points := imd.getAndClearPoints()

	if len(points) < 3 {
		imd.restorePoints(points)
		return
	}

	off := imd.tri.Len()
	imd.tri.SetLen(imd.tri.Len() + 3*(len(points)-2))

	for i, j := 1, off; i+1 < len(points); i, j = i+1, j+3 {
		for k, p := range [...]int{0, i, i + 1} {
			tri := &(*imd.tri)[j+k]
			tri.Position = points[p].pos
			tri.Color = points[p].col
			tri.Picture = points[p].pic
			tri.Intensity = points[p].in
		}
	}

	imd.applyMatrixAndMask(off)
	imd.batch.Dirty()

	imd.restorePoints(points)
}

// fillEllipseArc fills an elliptical arc defined by the specified radius, start angle (low), and end angle (high).
func (imd *IMDraw) fillEllipseArc(radius Vector, low, high float64) {
	points := imd.getAndClearPoints()

	for _, pt := range points {
		num := math.Ceil(math.Abs(high-low) / (2 * math.Pi) * float64(pt.precision))
		delta := (high - low) / num

		off := imd.tri.Len()
		imd.tri.SetLen(imd.tri.Len() + 3*int(num))

		for i := range (*imd.tri)[off:] {
			(*imd.tri)[off+i].Color = pt.col
			(*imd.tri)[off+i].Picture = ZeroVector
			(*imd.tri)[off+i].Intensity = 0
		}

		for i, j := 0.0, off; i < num; i, j = i+1, j+3 {
			angle := low + i*delta
			sin, cos := math.Sincos(angle)
			a := pt.pos.Add(NewVec(
				radius.X*cos,
				radius.Y*sin,
			))

			angle = low + (i+1)*delta
			sin, cos = math.Sincos(angle)
			b := pt.pos.Add(NewVec(
				radius.X*cos,
				radius.Y*sin,
			))

			(*imd.tri)[j+0].Position = pt.pos
			(*imd.tri)[j+1].Position = a
			(*imd.tri)[j+2].Position = b
		}

		imd.applyMatrixAndMask(off)
		imd.batch.Dirty()
	}

	imd.restorePoints(points)
}

// outlineEllipseArc draws an outlined elliptical arc with the given radius, angle range, and thickness.
// radius specifies the semi-axes' lengths of the ellipse.
// low and high define the starting and ending angles (in radians) of the arc.
// thickness determines the width of the outlined arc.
// doEndShape specifies whether to draw end caps for the arc.
func (imd *IMDraw) outlineEllipseArc(radius Vector, low, high, thickness float64, doEndShape bool) {
	points := imd.getAndClearPoints()

	for _, pt := range points {
		num := math.Ceil(math.Abs(high-low) / (2 * math.Pi) * float64(pt.precision))
		delta := (high - low) / num

		off := imd.tri.Len()
		imd.tri.SetLen(imd.tri.Len() + 6*int(num))

		for i := range (*imd.tri)[off:] {
			(*imd.tri)[off+i].Color = pt.col
			(*imd.tri)[off+i].Picture = ZeroVector
			(*imd.tri)[off+i].Intensity = 0
		}

		for i, j := 0.0, off; i < num; i, j = i+1, j+6 {
			angle := low + i*delta
			sin, cos := math.Sincos(angle)
			normalSin, normalCos := NewVec(sin, cos).ScaledXY(radius).Unit().XY()
			a := pt.pos.Add(NewVec(
				radius.X*cos-thickness/2*normalCos,
				radius.Y*sin-thickness/2*normalSin,
			))
			b := pt.pos.Add(NewVec(
				radius.X*cos+thickness/2*normalCos,
				radius.Y*sin+thickness/2*normalSin,
			))

			angle = low + (i+1)*delta
			sin, cos = math.Sincos(angle)
			normalSin, normalCos = NewVec(sin, cos).ScaledXY(radius).Unit().XY()
			c := pt.pos.Add(NewVec(
				radius.X*cos-thickness/2*normalCos,
				radius.Y*sin-thickness/2*normalSin,
			))
			d := pt.pos.Add(NewVec(
				radius.X*cos+thickness/2*normalCos,
				radius.Y*sin+thickness/2*normalSin,
			))

			(*imd.tri)[j+0].Position = a
			(*imd.tri)[j+1].Position = b
			(*imd.tri)[j+2].Position = c
			(*imd.tri)[j+3].Position = c
			(*imd.tri)[j+4].Position = b
			(*imd.tri)[j+5].Position = d
		}

		imd.applyMatrixAndMask(off)
		imd.batch.Dirty()

		if doEndShape {
			lowSin, lowCos := math.Sincos(low)
			lowCenter := pt.pos.Add(NewVec(
				radius.X*lowCos,
				radius.Y*lowSin,
			))
			normalLowSin, normalLowCos := NewVec(lowSin, lowCos).ScaledXY(radius).Unit().XY()
			normalLow := NewVec(normalLowCos, normalLowSin).Angle()

			highSin, highCos := math.Sincos(high)
			highCenter := pt.pos.Add(NewVec(
				radius.X*highCos,
				radius.Y*highSin,
			))
			normalHighSin, normalHighCos := NewVec(highSin, highCos).ScaledXY(radius).Unit().XY()
			normalHigh := NewVec(normalHighCos, normalHighSin).Angle()

			orientation := 1.0
			if low > high {
				orientation = -1.0
			}

			switch pt.endShape {
			case NoEndShape:
				// nothing
			case SharpEndShape:
				thick := NewVec(thickness/2, 0).Rotated(normalLow)
				imd.pushPt(lowCenter.Add(thick), pt)
				imd.pushPt(lowCenter.Sub(thick), pt)
				imd.pushPt(lowCenter.Sub(thick.Normal().Scaled(orientation)), pt)
				imd.fillPolygon()
				thick = NewVec(thickness/2, 0).Rotated(normalHigh)
				imd.pushPt(highCenter.Add(thick), pt)
				imd.pushPt(highCenter.Sub(thick), pt)
				imd.pushPt(highCenter.Add(thick.Normal().Scaled(orientation)), pt)
				imd.fillPolygon()
			case RoundEndShape:
				imd.pushPt(lowCenter, pt)
				imd.fillEllipseArc(NewVec(thickness/2, thickness/2), normalLow, normalLow-math.Pi*orientation)
				imd.pushPt(highCenter, pt)
				imd.fillEllipseArc(NewVec(thickness/2, thickness/2), normalHigh, normalHigh+math.Pi*orientation)
			}
		}
	}

	imd.restorePoints(points)
}

// polyline draws a series of connected line segments with a specified thickness.
// closed determines whether the start and end points are connected to form a loop.
func (imd *IMDraw) polyline(thickness float64, closed bool) {
	points := imd.getAndClearPoints()

	if len(points) == 0 {
		imd.restorePoints(points)
		return
	}
	if len(points) == 1 {
		// one Point special case
		points = append(points, points[0])
	}

	// first Point
	j1, i1 := 0, 1
	ijNormal := points[0].pos.To(points[1].pos).Normal().Unit().Scaled(thickness / 2)

	if !closed {
		switch points[j1].endShape {
		case NoEndShape:
			// nothing
		case SharpEndShape:
			imd.pushPt(points[j1].pos.Add(ijNormal), points[j1])
			imd.pushPt(points[j1].pos.Sub(ijNormal), points[j1])
			imd.pushPt(points[j1].pos.Add(ijNormal.Normal()), points[j1])
			imd.fillPolygon()
		case RoundEndShape:
			imd.pushPt(points[j1].pos, points[j1])
			imd.fillEllipseArc(NewVec(thickness/2, thickness/2), ijNormal.Angle(), ijNormal.Angle()+math.Pi)
		}
	}

	imd.pushPt(points[j1].pos.Add(ijNormal), points[j1])
	imd.pushPt(points[j1].pos.Sub(ijNormal), points[j1])

	// middle points
	for i := 0; i < len(points); i++ {
		j, k := i+1, i+2

		closing := false
		if j >= len(points) {
			j %= len(points)
			closing = true
		}
		if k >= len(points) {
			if !closed {
				break
			}
			k %= len(points)
		}

		jkNormal := points[j].pos.To(points[k].pos).Normal().Unit().Scaled(thickness / 2)

		orientation := 1.0
		if ijNormal.Cross(jkNormal) > 0 {
			orientation = -1.0
		}

		imd.pushPt(points[j].pos.Sub(ijNormal), points[j])
		imd.pushPt(points[j].pos.Add(ijNormal), points[j])
		imd.fillPolygon()

		switch points[j].endShape {
		case NoEndShape:
			// nothing
		case SharpEndShape:
			imd.pushPt(points[j].pos, points[j])
			imd.pushPt(points[j].pos.Add(ijNormal.Scaled(orientation)), points[j])
			imd.pushPt(points[j].pos.Add(jkNormal.Scaled(orientation)), points[j])
			imd.fillPolygon()
		case RoundEndShape:
			imd.pushPt(points[j].pos, points[j])
			imd.fillEllipseArc(NewVec(thickness/2, thickness/2), ijNormal.Angle(), ijNormal.Angle()-math.Pi)
			imd.pushPt(points[j].pos, points[j])
			imd.fillEllipseArc(NewVec(thickness/2, thickness/2), jkNormal.Angle(), jkNormal.Angle()+math.Pi)
		}

		if !closing {
			imd.pushPt(points[j].pos.Add(jkNormal), points[j])
			imd.pushPt(points[j].pos.Sub(jkNormal), points[j])
		}
		// "next" normal becomes previous normal
		ijNormal = jkNormal
	}

	// last Point
	i1, j1 = len(points)-2, len(points)-1
	ijNormal = points[i1].pos.To(points[j1].pos).Normal().Unit().Scaled(thickness / 2)

	imd.pushPt(points[j1].pos.Sub(ijNormal), points[j1])
	imd.pushPt(points[j1].pos.Add(ijNormal), points[j1])
	imd.fillPolygon()

	if !closed {
		switch points[j1].endShape {
		case NoEndShape:
			// nothing
		case SharpEndShape:
			imd.pushPt(points[j1].pos.Add(ijNormal), points[j1])
			imd.pushPt(points[j1].pos.Sub(ijNormal), points[j1])
			imd.pushPt(points[j1].pos.Add(ijNormal.Normal().Scaled(-1)), points[j1])
			imd.fillPolygon()
		case RoundEndShape:
			imd.pushPt(points[j1].pos, points[j1])
			imd.fillEllipseArc(NewVec(thickness/2, thickness/2), ijNormal.Angle(), ijNormal.Angle()-math.Pi)
		}
	}

	imd.restorePoints(points)
}
