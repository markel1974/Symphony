package pixels

import "image/color"

// Sprite represents a drawable object consisting of texture data, transformation, masking, and frame information.
type Sprite struct {
	tri    *TrianglesData
	frame  Rect
	drawer *Drawer
	matrix Matrix
	mask   RGBA
}

// NewSprite creates and returns a new Sprite instance with default settings and an associated TrianglesData and Drawer.
func NewSprite() *Sprite {
	tri := NewTrianglesData(6)
	s := &Sprite{
		tri:    tri,
		drawer: NewDrawer(tri, nil, CacheModePicture),
		matrix: IM,
		mask:   Alpha(1),
	}
	return s
}

// NewSpriteFromPicture creates and initializes a new Sprite from the specified picture and frame rectangle.
func NewSpriteFromPicture(pic IPicture, frame Rect) *Sprite {
	s := NewSprite()
	s.Set(pic, frame)
	return s
}

// Set assigns a new picture and frame to the Sprite, recalculating data if the frame has changed.
func (s *Sprite) Set(pic IPicture, frame Rect) {
	s.drawer.SetPicture(pic)
	if frame != s.frame {
		s.frame = frame
		s.calcData()
	}
}

// SetCachedMode sets the cache mode for the Sprite, which influences how rendering resources are managed internally.
func (s *Sprite) SetCachedMode(cacheMode CacheMode) {
	s.drawer.SetCacheMode(cacheMode)
}

// Picture returns the IPicture currently associated with the Sprite through its Drawer.
func (s *Sprite) Picture() IPicture {
	return s.drawer.Picture()
}

// Frame returns the currently set rectangular frame of the Sprite as a Rect.
func (s *Sprite) Frame() Rect {
	return s.frame
}

// Draw renders the Sprite onto the provided ITarget, applying the specified transformation Matrix.
func (s *Sprite) Draw(t ITarget, matrix Matrix) {
	s.DrawColorMask(t, matrix, nil)
}

// DrawColorMask draws the Sprite onto the given ITarget with the specified transformation matrix and color mask.
func (s *Sprite) DrawColorMask(t ITarget, matrix Matrix, mask color.Color) {
	dirty := false
	if matrix != s.matrix {
		s.matrix = matrix
		dirty = true
	}
	if mask == nil {
		mask = Alpha(1)
	}
	rgba := ToRGBA(mask)
	if rgba != s.mask {
		s.mask = rgba
		dirty = true
	}
	if dirty {
		s.calcData()
	}
	s.drawer.Draw(t)
}

// calcData recalculates the positional and visual data for the sprite's triangles based on its frame, matrix, and mask.
func (s *Sprite) calcData() {
	center := s.frame.Center()
	horizontal := NewVec(s.frame.W()/2, 0)
	vertical := NewVec(0, s.frame.H()/2)

	(*s.tri)[0].Position = Vector{}.Sub(horizontal).Sub(vertical)
	(*s.tri)[1].Position = Vector{}.Add(horizontal).Sub(vertical)
	(*s.tri)[2].Position = Vector{}.Add(horizontal).Add(vertical)
	(*s.tri)[3].Position = Vector{}.Sub(horizontal).Sub(vertical)
	(*s.tri)[4].Position = Vector{}.Add(horizontal).Add(vertical)
	(*s.tri)[5].Position = Vector{}.Sub(horizontal).Add(vertical)

	for i := range *s.tri {
		(*s.tri)[i].Color = s.mask
		(*s.tri)[i].Picture = center.Add((*s.tri)[i].Position)
		(*s.tri)[i].Intensity = 1
		(*s.tri)[i].Position = s.matrix.Project((*s.tri)[i].Position)
	}

	s.drawer.Dirty()
}
