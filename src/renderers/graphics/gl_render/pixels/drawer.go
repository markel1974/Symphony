package pixels

// CacheMode represents the cache behavior used for managing rendering resources internally.
type CacheMode int

// CacheModeUpdate represents the cache mode for updates only.
// CacheModePicture represents the cache mode for pictures only.
// CacheModePictureUpdate represents the cache mode for both pictures and updates.
const (
	CacheModeUpdate        CacheMode = 1
	CacheModePicture       CacheMode = 2
	CacheModePictureUpdate CacheMode = 3
)

// drawerTarget is a struct that associates ITargetTriangles and ITargetPictures for rendering on a specific ITarget.
// It maintains state for efficient drawing, including caching of pictures and tracking updates with a clean flag.
type drawerTarget struct {
	tris  ITargetTriangles
	pic   ITargetPicture
	pics  map[IPicture]ITargetPicture
	clean bool
}

// Drawer is responsible for managing the rendering state using ITriangles and IPicture objects.
// It supports various cache modes and manages drawing operations for multiple ITargets.
type Drawer struct {
	triangles ITriangles
	picture   IPicture
	cacheMode CacheMode
	targets   map[ITarget]*drawerTarget
}

// NewDrawer initializes and returns a new Drawer with the specified ITriangles, IPicture, and CacheMode.
func NewDrawer(triangles ITriangles, picture IPicture, cacheMode CacheMode) *Drawer {
	return &Drawer{
		triangles: triangles,
		picture:   picture,
		cacheMode: cacheMode,
		targets:   make(map[ITarget]*drawerTarget),
	}
}

// Triangles returns the ITriangles instance associated with the Drawer.
func (d *Drawer) Triangles() ITriangles {
	return d.triangles
}

// SetTriangles sets the ITriangles instance for the Drawer. This determines the triangles to be drawn.
func (d *Drawer) SetTriangles(triangles ITriangles) {
	d.triangles = triangles
}

// Picture returns the IPicture instance currently stored in the Drawer.
func (d *Drawer) Picture() IPicture {
	return d.picture
}

// SetPicture assigns a new IPicture to the Drawer, replacing the current picture.
func (d *Drawer) SetPicture(picture IPicture) {
	d.picture = picture
}

// CacheMode returns the current caching mode used by the Drawer.
func (d *Drawer) CacheMode() CacheMode {
	return d.cacheMode
}

// SetCacheMode sets the cache mode for the Drawer, determining how cached resources like pictures are managed.
func (d *Drawer) SetCacheMode(cacheMode CacheMode) {
	d.cacheMode = cacheMode
}

// Dirty marks all drawer targets as not clean, indicating they need to be updated before being drawn.
func (d *Drawer) Dirty() {
	for _, dt := range d.targets {
		dt.clean = false
	}
}

// Draw renders the triangles and picture onto the specified target using the current cache mode of the Drawer.
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

// drawUpdate updates or creates a cached ITargetPicture for the given ITarget and draws it using the supplied triangles.
func (d *Drawer) drawUpdate(dt *drawerTarget, t ITarget) {
	if dt.pic == nil {
		dt.pic = t.MakePicture(d.picture)
	} else {
		dt.pic.Update(d.picture)
	}
	dt.pic.Draw(dt.tris)
}

// drawPicture handles drawing operations when the `CacheModePicture` caching mode is enabled.
// It ensures that an `ITargetPicture` is retrieved or created for the current `IPicture` and updates it if necessary.
// The `ITargetPicture` is then used to draw the associated `ITargetTriangles` on the provided `ITarget`.
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

// drawPictureUpdate ensures the picture is retrieved from or added to the cache and draws it with the given triangles.
func (d *Drawer) drawPictureUpdate(dt *drawerTarget, t ITarget) {
	pic := dt.pics[d.picture]
	if pic == nil {
		pic = t.MakePicture(d.picture)
		dt.pics[d.picture] = pic
	}
	pic.Draw(dt.tris)
}

// drawDefault draws the supplied ITargetTriangles directly using a newly created ITargetPicture.
func (d *Drawer) drawDefault(dt *drawerTarget, t ITarget) {
	pic := t.MakePicture(d.picture)
	pic.Draw(dt.tris)
}
