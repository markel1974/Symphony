package pixels

import "image/color"

type Sprite struct {
	tri    *TrianglesData
	frame  Rect
	d      *Drawer
	matrix Matrix
	mask   RGBA
}

func NewSprite() *Sprite {
	tri := NewTrianglesData(6)
	s := &Sprite{tri: tri, d: NewDrawer(tri, nil, CacheModePicture)}
	s.matrix = IM
	s.mask = Alpha(1)
	return s
}

func NewSpriteFromPicture(pic IPicture, frame Rect) *Sprite {
	s := NewSprite()
	s.Set(pic, frame)
	return s
}

func (s *Sprite) Set(pic IPicture, frame Rect) {
	s.d.SetPicture(pic)
	if frame != s.frame {
		s.frame = frame
		s.calcData()
	}
}

func (s *Sprite) SetCachedMode(cacheMode CacheMode) {
	s.d.SetCacheMode(cacheMode)
}

func (s *Sprite) Picture() IPicture {
	return s.d.Picture()
}

func (s *Sprite) Frame() Rect {
	return s.frame
}

func (s *Sprite) Draw(t ITarget, matrix Matrix) {
	s.DrawColorMask(t, matrix, nil)
}

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
	s.d.Draw(t)
}

func (s *Sprite) calcData() {
	center := s.frame.Center()
	horizontal := NewVec(s.frame.W()/2, 0)
	vertical := NewVec(0, s.frame.H()/2)

	(*s.tri)[0].Position = Vec{}.Sub(horizontal).Sub(vertical)
	(*s.tri)[1].Position = Vec{}.Add(horizontal).Sub(vertical)
	(*s.tri)[2].Position = Vec{}.Add(horizontal).Add(vertical)
	(*s.tri)[3].Position = Vec{}.Sub(horizontal).Sub(vertical)
	(*s.tri)[4].Position = Vec{}.Add(horizontal).Add(vertical)
	(*s.tri)[5].Position = Vec{}.Sub(horizontal).Add(vertical)

	for i := range *s.tri {
		(*s.tri)[i].Color = s.mask
		(*s.tri)[i].Picture = center.Add((*s.tri)[i].Position)
		(*s.tri)[i].Intensity = 1
		(*s.tri)[i].Position = s.matrix.Project((*s.tri)[i].Position)
	}

	s.d.Dirty()
}
