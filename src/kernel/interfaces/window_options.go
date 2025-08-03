package interfaces

// WindowOptions represents configurable parameters for a task, including offsets, scaling, and associated command line.
type WindowOptions struct {
	offsetY int
	offsetX int
	scale   float64
	caption string
}

func NewWindowOptions(offsetX int, offsetY int, scale float64) *WindowOptions {
	return &WindowOptions{
		offsetY: offsetY,
		offsetX: offsetX,
		scale:   scale,
		caption: "",
	}
}

// SetWindowOption updates the task's X, Y offsets or Scale based on the given option ('x', 'y', or 'z') and value.
func (t *WindowOptions) SetWindowOption(option rune, value float64) {
	switch option {
	case 'y':
		t.SetOffsetY(t.OffsetY() + int(value))
	case 'x':
		t.SetOffsetX(t.OffsetX() + int(value))
	case 'z':
		if scale := t.Scale() + value; scale >= 0.2 && scale <= 1 {
			t.SetScale(scale)
		}
	}
}

// OffsetX returns the X-axis offset value for the Process.
func (t *WindowOptions) OffsetX() int {
	return t.offsetX
}

// SetOffsetX sets the horizontal offset (offsetX) of the task to the specified value x.
func (t *WindowOptions) SetOffsetX(x int) {
	t.offsetX = x
}

// OffsetY returns the current vertical offset value for the task.
func (t *WindowOptions) OffsetY() int {
	return t.offsetY
}

// SetOffsetY sets the vertical offset value for the task.
func (t *WindowOptions) SetOffsetY(y int) {
	t.offsetY = y
}

// Scale returns the current scale factor of the task. It determines the zoom level or relative size adjustment.
func (t *WindowOptions) Scale() float64 {
	return t.scale
}

// SetScale sets the scale factor for the Process object to the specified value.
func (t *WindowOptions) SetScale(scale float64) {
	t.scale = scale
}

// SetWindowOptions updates the Process's properties using the provided TaskOptions. It returns immediately if options is nil.
func (t *WindowOptions) SetWindowOptions(options *WindowOptions) {
	if options == nil {
		return
	}
	t.scale = options.scale
	t.offsetX = options.offsetX
	t.offsetY = options.offsetY
}

// SetCaption updates the task's caption using a provided string and task ID, returning true to indicate successful update.
func (t *WindowOptions) SetCaption(caption string) bool {
	t.caption = caption
	return true
}

// Caption returns the current caption text for the window options.
func (t *WindowOptions) Caption() string {
	return t.caption
}
