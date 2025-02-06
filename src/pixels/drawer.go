package pixels

type CacheMode int

const (
	CacheModeUpdate        CacheMode = 1
	CacheModePicture       CacheMode = 2
	CacheModePictureUpdate CacheMode = 3
)

// Drawer glues all the fundamental interfaces (ITarget, ITriangles, IPicture) into a coherent and the
// only intended usage pattern.
//
// Drawer makes it possible to draw any combination of ITriangles and IPicture onto any ITarget
// efficiently.
//
// To create a Drawer, just assign its ITriangles and IPicture fields:
//
//	d := pixel.Drawer{ITriangles: t, IPicture: p}
//
// If ITriangles is nil, nothing will be drawn.
// If IPicture is nil, ITriangles will be drawn without an IPicture.
//
// Whenever you change the ITriangles, call Dirty to notify Drawer that ITriangles changed. You don't
// need to notify Drawer about a change of the IPicture.
//
// Note that Drawer caches the results of MakePicture from Targets it's drawn to for each IPicture
// it's set to. What it means is that using a Drawer with an unbounded number of Pictures leads to a
// memory leak, since the Drawer caches them and never forgets. In such a situation, create a new Drawer
// for each IPicture.
type Drawer struct {
	triangles ITriangles
	picture   IPicture
	cacheMode CacheMode
	targets   map[ITarget]*drawerTarget
}

type drawerTarget struct {
	tris  ITargetTriangles
	pic   ITargetPicture
	pics  map[IPicture]ITargetPicture
	clean bool
}

// NewDrawer creates and initializes a new Drawer instance with the provided ITriangles, IPicture, and CacheMode.
func NewDrawer(triangles ITriangles, picture IPicture, cacheMode CacheMode) *Drawer {
	return &Drawer{
		triangles: triangles,
		picture:   picture,
		cacheMode: cacheMode,
		targets:   make(map[ITarget]*drawerTarget),
	}
}

// Triangles retrieves the ITriangles instance currently set in the Drawer.
func (d *Drawer) Triangles() ITriangles {
	return d.triangles
}

// SetTriangles sets the ITriangles instance to be used by the Drawer for later drawing operations.
func (d *Drawer) SetTriangles(triangles ITriangles) {
	d.triangles = triangles
}

// Picture retrieves the current IPicture instance associated with the Drawer.
func (d *Drawer) Picture() IPicture {
	return d.picture
}

// SetPicture assigns the provided IPicture instance to the Drawer for use in later drawing operations.
func (d *Drawer) SetPicture(picture IPicture) {
	d.picture = picture
}

// CacheMode retrieves the current caching strategy used by the Drawer for managing pictures during drawing operations.
func (d *Drawer) CacheMode() CacheMode {
	return d.cacheMode
}

// SetCacheMode sets the caching strategy for the Drawer, determining how pictures are managed during drawing operations.
func (d *Drawer) SetCacheMode(cacheMode CacheMode) {
	d.cacheMode = cacheMode
}

// Dirty marks the ITriangles of this Drawer as changed. If not called, changes will not be visible when drawing.
func (d *Drawer) Dirty() {
	for _, dt := range d.targets {
		dt.clean = false
	}
}

// Draw efficiently draws ITriangles with IPicture onto the provided ITarget.
// If ITriangles is nil, nothing will be drawn.
// If IPicture is nil, ITriangles will be drawn without an IPicture.
func (d *Drawer) Draw(t ITarget) {
	if d.triangles == nil {
		return
	}

	dt := d.targets[t]
	if dt == nil {
		dt = &drawerTarget{pics: make(map[IPicture]ITargetPicture)}
		d.targets[t] = dt
	}

	if dt.tris == nil {
		dt.tris = t.MakeTriangles(d.triangles)
		dt.clean = true
	}

	if !dt.clean {
		dt.tris.SetLen(d.triangles.Len())
		dt.tris.Update(d.triangles)
		dt.clean = true
	}

	if d.picture == nil {
		dt.tris.Draw()
		return
	}

	switch d.cacheMode {
	case CacheModeUpdate:
		d.drawUpdate(dt, t)
	case CacheModePicture:
		d.drawPicture(dt, t)
	case CacheModePictureUpdate:
		d.drawPictureUpdate(dt, t)
	default:
		d.drawDefault(dt, t)
	}
}

// drawUpdate ensures the target's cached picture is created or updated before drawing triangles onto the target.
func (d *Drawer) drawUpdate(dt *drawerTarget, t ITarget) {
	if dt.pic == nil {
		dt.pic = t.MakePicture(d.picture)
	} else {
		dt.pic.Update(d.picture)
	}
	dt.pic.Draw(dt.tris)
}

// drawPicture retrieves or creates a cached picture from the target using the current IPicture.
// If the picture already exists, it updates it; otherwise, it creates a new one and caches it.
// Finally, it draws the triangles using the picture onto the target.
func (d *Drawer) drawPicture(dt *drawerTarget, t ITarget) {
	pic := dt.pics[d.picture]
	if pic == nil {
		pic = t.MakePicture(d.picture)
		dt.pics[d.picture] = pic
	} else {
		pic.Update(d.picture)
	}
	pic.Draw(dt.tris)
}

// drawPictureUpdate retrieves or creates a cached picture for the target, then draws triangles using that picture.
func (d *Drawer) drawPictureUpdate(dt *drawerTarget, t ITarget) {
	pic := dt.pics[d.picture]
	if pic == nil {
		pic = t.MakePicture(d.picture)
		dt.pics[d.picture] = pic
	}
	pic.Draw(dt.tris)
}

// drawDefault draws triangles onto the target using a newly created picture from the current IPicture without caching.
func (d *Drawer) drawDefault(dt *drawerTarget, t ITarget) {
	pic := t.MakePicture(d.picture)
	pic.Draw(dt.tris)
}
