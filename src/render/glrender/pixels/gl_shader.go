package pixels

import (
	"errors"
	executor2 "github.com/markel1974/c64emu/src/render/glrender/pixels/executor"

	"github.com/go-gl/mathgl/mgl32"
)

// canvasPosition is the iota constant representing the position data for a canvas.
// canvasColor is the iota constant representing the color data for a canvas.
// canvasTexCoords is the iota constant representing the texture coordinates for a canvas.
// canvasIntensity is the iota constant representing the intensity data for a canvas.
// canvasClip is the iota constant representing the clipping data for a canvas.
const (
	canvasPosition int = iota
	canvasColor
	canvasTexCoords
	canvasIntensity
	canvasClip
)

// defaultCanvasVertexFormat defines the default set of vertex attributes for canvas rendering, including position, color, texture coordinates, intensity, and clipping information.
var defaultCanvasVertexFormat = executor2.AttrFormat{
	canvasPosition:  executor2.Attr{Name: "aPosition", Type: executor2.Vec2},
	canvasColor:     executor2.Attr{Name: "aColor", Type: executor2.Vec4},
	canvasTexCoords: executor2.Attr{Name: "aTexCoords", Type: executor2.Vec2},
	canvasIntensity: executor2.Attr{Name: "aIntensity", Type: executor2.Float},
	canvasClip:      executor2.Attr{Name: "aClipRect", Type: executor2.Vec4},
}

// GLShader represents an OpenGL shader that encapsulates vertex and fragment shaders along with their configurations.
type GLShader struct {
	s      *executor2.Shader
	vf, uf executor2.AttrFormat
	vs, fs string

	uniforms []gsUniformAttr

	uniformDefaults struct {
		transform mgl32.Mat3
		colorMask mgl32.Vec4
		bounds    mgl32.Vec4
		texBounds mgl32.Vec4
		clipRect  mgl32.Vec4
	}
}

// NewGLShader creates a new GLShader instance using the provided fragment shader source code.
// It initializes default uniforms and updates the shader configuration before returning the instance.
func NewGLShader(fragmentShader string) *GLShader {
	gs := &GLShader{
		vf: defaultCanvasVertexFormat,
		vs: baseCanvasVertexShader,
		fs: fragmentShader,
	}

	gs.SetUniform("uTransform", &gs.uniformDefaults.transform)
	gs.SetUniform("uColorMask", &gs.uniformDefaults.colorMask)
	gs.SetUniform("uBounds", &gs.uniformDefaults.bounds)
	gs.SetUniform("uTexBounds", &gs.uniformDefaults.texBounds)

	gs.Update()

	return gs
}

// Update rebuilds the shader using the current vertex format, uniform format, and shader code, applying any changes.
func (gs *GLShader) Update() {
	gs.uf = make([]executor2.Attr, len(gs.uniforms))
	for idx := range gs.uniforms {
		gs.uf[idx] = executor2.Attr{
			Name: gs.uniforms[idx].Name,
			Type: gs.uniforms[idx].Type,
		}
	}

	var shader *executor2.Shader
	executor2.GraphicThread.Call(func() {
		var err error
		shader, err = executor2.NewShader(gs.vf, gs.uf, gs.vs, gs.fs)
		if err != nil {
			panic(errors.New("failed to create GLCanvas, there's a bug in the shader:" + err.Error()))
		}
	})

	gs.s = shader
}

// getUniform searches for a uniform variable by its name in the shader and returns its index or -1 if not found.
func (gs *GLShader) getUniform(Name string) int {
	for i, u := range gs.uniforms {
		if u.Name == Name {
			return i
		}
	}
	return -1
}

// SetUniform assigns a value to a uniform variable by name, updating it if it exists or creating a new entry if it does not.
func (gs *GLShader) SetUniform(name string, value interface{}) {
	t, p := getAttrType(value)
	if loc := gs.getUniform(name); loc > -1 {
		gs.uniforms[loc].Name = name
		gs.uniforms[loc].Type = t
		gs.uniforms[loc].isPointer = p
		gs.uniforms[loc].value = value
		return
	}
	gs.uniforms = append(gs.uniforms, gsUniformAttr{
		Name:      name,
		Type:      t,
		isPointer: p,
		value:     value,
	})
}

// gsUniformAttr represents a single uniform attribute in a graphics shader.
// Name specifies the name of the uniform attribute.
// Type defines the data type of the attribute, based on executor.AttrType.
// value stores the actual data of the uniform attribute.
// isPointer indicates whether the value is a pointer type.
type gsUniformAttr struct {
	Name      string
	Type      executor2.AttrType
	value     interface{}
	isPointer bool
}

// Value returns the stored value of the gsUniformAttr instance, resolving pointer dereferences based on its type.
func (gu *gsUniformAttr) Value() interface{} {
	if !gu.isPointer {
		return gu.value
	}
	switch gu.Type {
	case executor2.Vec2:
		return *gu.value.(*mgl32.Vec2)
	case executor2.Vec3:
		return *gu.value.(*mgl32.Vec3)
	case executor2.Vec4:
		return *gu.value.(*mgl32.Vec4)
	case executor2.Mat2:
		return *gu.value.(*mgl32.Mat2)
	case executor2.Mat23:
		return *gu.value.(*mgl32.Mat2x3)
	case executor2.Mat24:
		return *gu.value.(*mgl32.Mat2x4)
	case executor2.Mat3:
		return *gu.value.(*mgl32.Mat3)
	case executor2.Mat32:
		return *gu.value.(*mgl32.Mat3x2)
	case executor2.Mat34:
		return *gu.value.(*mgl32.Mat3x4)
	case executor2.Mat4:
		return *gu.value.(*mgl32.Mat4)
	case executor2.Mat42:
		return *gu.value.(*mgl32.Mat4x2)
	case executor2.Mat43:
		return *gu.value.(*mgl32.Mat4x3)
	case executor2.Int:
		return *gu.value.(*int32)
	case executor2.Float:
		return *gu.value.(*float32)
	default:
		panic("invalid attr type (Value)")
	}
}

// getAttrType determines the attribute type and pointer status of the given value for use in the shader execution context.
// It returns the corresponding executor.AttrType and a boolean indicating whether the value is a pointer type.
// Panics if the provided value is of an invalid or unsupported attribute type.
func getAttrType(v interface{}) (executor2.AttrType, bool) {
	switch v.(type) {
	case int32:
		return executor2.Int, false
	case float32:
		return executor2.Float, false
	case mgl32.Vec2:
		return executor2.Vec2, false
	case mgl32.Vec3:
		return executor2.Vec3, false
	case mgl32.Vec4:
		return executor2.Vec4, false
	case mgl32.Mat2:
		return executor2.Mat2, false
	case mgl32.Mat2x3:
		return executor2.Mat23, false
	case mgl32.Mat2x4:
		return executor2.Mat24, false
	case mgl32.Mat3:
		return executor2.Mat3, false
	case mgl32.Mat3x2:
		return executor2.Mat32, false
	case mgl32.Mat3x4:
		return executor2.Mat34, false
	case mgl32.Mat4:
		return executor2.Mat4, false
	case mgl32.Mat4x2:
		return executor2.Mat42, false
	case mgl32.Mat4x3:
		return executor2.Mat43, false
	case *mgl32.Vec2:
		return executor2.Vec2, true
	case *mgl32.Vec3:
		return executor2.Vec3, true
	case *mgl32.Vec4:
		return executor2.Vec4, true
	case *mgl32.Mat2:
		return executor2.Mat2, true
	case *mgl32.Mat2x3:
		return executor2.Mat23, true
	case *mgl32.Mat2x4:
		return executor2.Mat24, true
	case *mgl32.Mat3:
		return executor2.Mat3, true
	case *mgl32.Mat3x2:
		return executor2.Mat32, true
	case *mgl32.Mat3x4:
		return executor2.Mat34, true
	case *mgl32.Mat4:
		return executor2.Mat4, true
	case *mgl32.Mat4x2:
		return executor2.Mat42, true
	case *mgl32.Mat4x3:
		return executor2.Mat43, true
	case *int32:
		return executor2.Int, true
	case *float32:
		return executor2.Float, true
	default:
		panic("invalid attr type (getAttrType)")
	}
}
