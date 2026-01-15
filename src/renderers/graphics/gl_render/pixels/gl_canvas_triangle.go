package pixels

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/markel1974/symphony/src/renderers/graphics/gl_render/pixels/executor"
)

// GLCanvasTriangles represents a drawable set of triangles specific to a GLCanvas context.
// It embeds GLTriangles and includes functionality for integration with the associated GLCanvas.
type GLCanvasTriangles struct {
	*GLTriangles
	dst *GLCanvas
}

// NewGLCanvasTriangles creates and returns a new GLCanvasTriangles instance using the given GLTriangles and GLCanvas.
func NewGLCanvasTriangles(t *GLTriangles, dst *GLCanvas) *GLCanvasTriangles {
	return &GLCanvasTriangles{
		GLTriangles: t,
		dst:         dst,
	}
}

// draw renders GLCanvasTriangles onto the frame using the provided texture and bounds within a graphics thread context.
func (ct *GLCanvasTriangles) draw(tex *executor.Texture, bounds Rect) {
	ct.dst.gf.Dirty()

	// save the current state vars to avoid race condition
	cmp := ct.dst.cmp
	smt := ct.dst.smooth
	mat := ct.dst.mat
	col := ct.dst.col

	executor.GraphicThread.Post(func() {
		ct.dst.setGlBounds()
		setBlendFunc(cmp)

		frame := ct.dst.gf.Frame()
		shader := ct.shader.s

		frame.Begin()
		shader.Begin()

		ct.shader.uniformDefaults.transform = mat
		ct.shader.uniformDefaults.colorMask = col
		dstBounds := ct.dst.Bounds()
		ct.shader.uniformDefaults.bounds = mgl32.Vec4{
			float32(dstBounds.Min.X),
			float32(dstBounds.Min.Y),
			float32(dstBounds.W()),
			float32(dstBounds.H()),
		}

		bx, by, bw, bh := intBounds(bounds)
		ct.shader.uniformDefaults.texBounds = mgl32.Vec4{
			float32(bx),
			float32(by),
			float32(bw),
			float32(bh),
		}

		for loc, u := range ct.shader.uniforms {
			_, _ = ct.shader.s.SetUniformAttr(loc, u.Value())
		}

		if tex == nil {
			ct.vs.Begin()
			ct.vs.Draw()
			ct.vs.End()
		} else {
			tex.Begin()
			if tex.Smooth() != smt {
				tex.SetSmooth(smt)
			}
			ct.vs.Begin()
			ct.vs.Draw()
			ct.vs.End()
			tex.End()
		}

		shader.End()
		frame.End()
	})
}

// Draw draws the triangles in the associated GLCanvasTriangles object onto the destination canvas.
func (ct *GLCanvasTriangles) Draw() {
	ct.draw(nil, Rect{})
}
