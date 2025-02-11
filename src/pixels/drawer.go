package pixels

// CacheMode defines the mode for managing caching behavior in rendering systems.
type CacheMode int

// CacheModeUpdate represents the cache mode for updating purposes.
// CacheModePicture represents the cache mode for picture storage.
// CacheModePictureUpdate represents the cache mode for updating and storing pictures.
const (
	CacheModeUpdate        CacheMode = 1
	CacheModePicture       CacheMode = 2
	CacheModePictureUpdate CacheMode = 3
)

// Drawer is a structure responsible for managing rendering operations using triangles, pictures, and cache modes.
type Drawer struct {
	triangles ITriangles
	picture   IPicture
	cacheMode CacheMode
	targets   map[ITarget]*drawerTarget
}

// drawerTarget is a struct that encapsulates target-specific representations of triangles and pictures for drawing.
// It includes an ITargetTriangles instance, an ITargetPicture instance, a map for cached ITargetPictures by IPicture keys, and a clean flag.
type drawerTarget struct {
	tris  ITargetTriangles
	pic   ITargetPicture
	pics  map[IPicture]ITargetPicture
	clean bool
}

// NewDrawer creates a new Drawer instance with the provided triangles, picture, and cache mode settings.
func NewDrawer(triangles ITriangles, picture IPicture, cacheMode CacheMode) *Drawer {
	return &Drawer{
		triangles: triangles,
		picture:   picture,
		cacheMode: cacheMode,
		targets:   make(map[ITarget]*drawerTarget),
	}
}

// Triangles retrieves the current ITriangles instance associated with the Drawer.
func (d *Drawer) Triangles() ITriangles {
	return d.triangles
}

// SetTriangles assigns the provided ITriangles instance to the Drawer for use in rendering operations.
func (d *Drawer) SetTriangles(triangles ITriangles) {
	d.triangles = triangles
}

// Picture returns the IPicture associated with the Drawer. It represents the current drawable image context.
func (d *Drawer) Picture() IPicture {
	return d.picture
}

// SetPicture sets the IPicture instance to be used by the Drawer for rendering operations.
func (d *Drawer) SetPicture(picture IPicture) {
	d.picture = picture
}

// CacheMode returns the current cache mode used by the Drawer.
func (d *Drawer) CacheMode() CacheMode {
	return d.cacheMode
}

// SetCacheMode sets the cache mode for the Drawer, determining how resources are managed during rendering.
func (d *Drawer) SetCacheMode(cacheMode CacheMode) {
	d.cacheMode = cacheMode
}

// Dirty marks all targets in the Drawer as needing an update by setting their clean state to false.
func (d *Drawer) Dirty() {
	for _, dt := range d.targets {
		dt.clean = false
	}
}

// Draw renders the stored triangles and picture onto the given ITarget, adapting behavior based on the cache mode.
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

// drawUpdate manages the update and drawing process of a cached picture for the specified drawer target on the given target.
func (d *Drawer) drawUpdate(dt *drawerTarget, t ITarget) {
	if dt.pic == nil {
		dt.pic = t.MakePicture(d.picture)
	} else {
		dt.pic.Update(d.picture)
	}
	dt.pic.Draw(dt.tris)
}

// drawPicture initializes or updates a cached picture for the given drawer target and draws it with specified triangles.
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

// drawPictureUpdate ensures an ITargetPicture is created or reused from the cache and draws it with the given triangles.
func (d *Drawer) drawPictureUpdate(dt *drawerTarget, t ITarget) {
	pic := dt.pics[d.picture]
	if pic == nil {
		pic = t.MakePicture(d.picture)
		dt.pics[d.picture] = pic
	}
	pic.Draw(dt.tris)
}

// drawDefault creates a picture from the provided IPicture and draws it using the associated target triangles.
func (d *Drawer) drawDefault(dt *drawerTarget, t ITarget) {
	pic := t.MakePicture(d.picture)
	pic.Draw(dt.tris)
}
