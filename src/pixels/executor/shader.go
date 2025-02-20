package executor

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Shader represents a compiled, linked, and usable GPU shader program with associated vertex and uniform attribute formats.
type Shader struct {
	program    Binder
	vertexFmt  AttrFormat
	uniformFmt AttrFormat
	uniformLoc []int32
}

// NewShader creates a new Shader object using specified vertex and uniform formats and given GLSL shader code.
// It compiles the vertex and fragment shaders, links them into a shader program, and initializes uniform locations.
// Returns a pointer to the created Shader and an error if shader compilation or linking fails.
func NewShader(vertexFmt, uniformFmt AttrFormat, vertexShader, fragmentShader string) (*Shader, error) {
	shader := &Shader{
		program: Binder{
			restoreLoc: gl.CURRENT_PROGRAM,
			bindFunc: func(obj uint32) {
				gl.UseProgram(obj)
			},
		},
		vertexFmt:  vertexFmt,
		uniformFmt: uniformFmt,
		uniformLoc: make([]int32, len(uniformFmt)),
	}

	var vShader, fShader uint32

	{
		vShader = gl.CreateShader(gl.VERTEX_SHADER)
		src, free := gl.Strs(vertexShader)
		defer free()
		length := int32(len(vertexShader))
		gl.ShaderSource(vShader, 1, src, &length)
		gl.CompileShader(vShader)

		var success int32
		gl.GetShaderiv(vShader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(vShader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(vShader, logLen, nil, &infoLog[0])
			return nil, fmt.Errorf("error compiling vertex shader: %s", string(infoLog))
		}

		defer gl.DeleteShader(vShader)
	}

	{
		fShader = gl.CreateShader(gl.FRAGMENT_SHADER)
		src, free := gl.Strs(fragmentShader)
		defer free()
		length := int32(len(fragmentShader))
		gl.ShaderSource(fShader, 1, src, &length)
		gl.CompileShader(fShader)

		var success int32
		gl.GetShaderiv(fShader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(fShader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(fShader, logLen, nil, &infoLog[0])
			return nil, fmt.Errorf("error compiling fragment shader: %s", string(infoLog))
		}

		defer gl.DeleteShader(fShader)
	}

	{
		shader.program.obj = gl.CreateProgram()
		gl.AttachShader(shader.program.obj, vShader)
		gl.AttachShader(shader.program.obj, fShader)
		gl.LinkProgram(shader.program.obj)

		var success int32
		gl.GetProgramiv(shader.program.obj, gl.LINK_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetProgramiv(shader.program.obj, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetProgramInfoLog(shader.program.obj, logLen, nil, &infoLog[0])
			return nil, fmt.Errorf("error linking shader program: %s", string(infoLog))
		}
	}

	for i, uniform := range uniformFmt {
		loc := gl.GetUniformLocation(shader.program.obj, gl.Str(uniform.Name+"\x00"))
		shader.uniformLoc[i] = loc
	}

	runtime.SetFinalizer(shader, (*Shader).delete)

	return shader, nil
}

// delete releases the OpenGL shader program associated with the shader object using the graphic thread.
func (s *Shader) delete() {
	GraphicThread.Post(func() {
		gl.DeleteProgram(s.program.obj)
	})
}

// ID returns the OpenGL program ID associated with the Shader.
func (s *Shader) ID() uint32 {
	return s.program.obj
}

// VertexFormat returns the vertex attribute format (AttrFormat) used by the Shader.
func (s *Shader) VertexFormat() AttrFormat {
	return s.vertexFmt
}

// UniformFormat returns the format of the uniform attributes used by the shader.
func (s *Shader) UniformFormat() AttrFormat {
	return s.uniformFmt
}

// SetUniformAttr assigns a value to a specified uniform attribute in the shader program.
// It accepts the uniform index and a value of the appropriate type.
// Returns true if the uniform was successfully set, or false and an error if it fails.
func (s *Shader) SetUniformAttr(uniform int, value interface{}) (bool, error) {
	if s.uniformLoc[uniform] < 0 {
		return false, nil
	}

	switch s.uniformFmt[uniform].Type {
	case Int:
		value := value.(int32)
		gl.Uniform1iv(s.uniformLoc[uniform], 1, &value)
	case Float:
		value := value.(float32)
		gl.Uniform1fv(s.uniformLoc[uniform], 1, &value)
	case Vec2:
		value := value.(mgl32.Vec2)
		gl.Uniform2fv(s.uniformLoc[uniform], 1, &value[0])
	case Vec3:
		value := value.(mgl32.Vec3)
		gl.Uniform3fv(s.uniformLoc[uniform], 1, &value[0])
	case Vec4:
		value := value.(mgl32.Vec4)
		gl.Uniform4fv(s.uniformLoc[uniform], 1, &value[0])
	case Mat2:
		value := value.(mgl32.Mat2)
		gl.UniformMatrix2fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat23:
		value := value.(mgl32.Mat2x3)
		gl.UniformMatrix2x3fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat24:
		value := value.(mgl32.Mat2x4)
		gl.UniformMatrix2x4fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat3:
		value := value.(mgl32.Mat3)
		gl.UniformMatrix3fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat32:
		value := value.(mgl32.Mat3x2)
		gl.UniformMatrix3x2fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat34:
		value := value.(mgl32.Mat3x4)
		gl.UniformMatrix3x4fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat4:
		value := value.(mgl32.Mat4)
		gl.UniformMatrix4fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat42:
		value := value.(mgl32.Mat4x2)
		gl.UniformMatrix4x2fv(s.uniformLoc[uniform], 1, false, &value[0])
	case Mat43:
		value := value.(mgl32.Mat4x3)
		gl.UniformMatrix4x3fv(s.uniformLoc[uniform], 1, false, &value[0])
	default:
		return false, errors.New("set uniform attr: invalid attribute type")
	}
	return true, nil
}

// Begin activates the shader program by binding it to the current OpenGL context.
func (s *Shader) Begin() {
	s.program.Bind()
}

// End restores the previous OpenGL state that was active before the shader was bound.
func (s *Shader) End() {
	s.program.Restore()
}
