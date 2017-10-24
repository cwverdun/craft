package glw

import (
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//Program is the gl Program
type Program uint32

//ShaderSource the source code and the type of the shader
type ShaderSource struct {
	Type   uint32
	Source io.Reader
}

//NewProgram reads an compiles shaders and links new program
func NewProgram(shaderSources ...ShaderSource) (Program, error) {
	sources := make([]string, len(shaderSources))
	for i, s := range shaderSources {
		source, err := ioutil.ReadAll(s.Source)
		if err != nil {
			return 0, fmt.Errorf("cannot read shader %d: %s", i, err)
		}
		sources[i] = string(source) + "\x00"
	}
	shaders := make([]uint32, len(shaderSources))
	for i := 0; i < len(shaderSources); i++ {
		shader, err := compileShader(sources[i], shaderSources[i].Type)
		if err != nil {
			return 0, fmt.Errorf("shader %d compile error: %s", i, err)
		}
		shaders[i] = shader
	}
	program := gl.CreateProgram()
	for _, s := range shaders {
		gl.AttachShader(program, s)
	}

	gl.LinkProgram(program)
	if err := getShaderError(program, gl.LINK_STATUS); err != nil {
		return 0, err
	}
	gl.ValidateProgram(program)
	if err := getShaderError(program, gl.VALIDATE_STATUS); err != nil {
		return 0, err
	}

	for _, s := range shaders {
		gl.DeleteShader(s)
	}

	return Program(program), nil
}

func getShaderError(program uint32, pname uint32) error {
	var status int32
	gl.GetProgramiv(program, pname, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])
		return fmt.Errorf("failed to link program: %v", string(log))
	}
	return nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
		return 0, fmt.Errorf("failed to compile: %v", string(log))
	}

	return shader, nil
}

//UseProgram calls gl.UseProgram
func (p Program) UseProgram() {
	gl.UseProgram(uint32(p))
}

//GetUniformLocation call gl.GetUniformLocation
func (p Program) GetUniformLocation(name string) UniformLocation {
	return UniformLocation(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")))
}

//GetAttribLocation returns VertexAttrib with associated name
func (p Program) GetAttribLocation(name string) VertexAttrib {
	return VertexAttrib(gl.GetAttribLocation(uint32(p), gl.Str(name+"\x00")))
}

//BindFragDataLocation ...
func (p Program) BindFragDataLocation(color uint32, name string) {
	gl.BindFragDataLocation(uint32(p), color, gl.Str(name+"\x00"))
}

//Delete deletes the program
func (p Program) Delete() {
	gl.DeleteProgram(uint32(p))
}

//VertexAttrib ...
type VertexAttrib uint32

//EnableVertexAttribArray ...
func (va VertexAttrib) EnableVertexAttribArray() {
	gl.EnableVertexAttribArray(uint32(va))
}

//VertexAttribPointer ...
func (va VertexAttrib) VertexAttribPointer(size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	gl.VertexAttribPointer(uint32(va), size, xtype, normalized, stride, pointer)
}

//DisableVertexAttribArray ...
func (va VertexAttrib) DisableVertexAttribArray() {
	gl.DisableVertexAttribArray(uint32(va))
}

//UniformLocation ...
type UniformLocation int32

//UniformMatrix4fv ...
func (ul UniformLocation) UniformMatrix4fv(count int32, transpose bool, m mgl32.Mat4) {
	gl.UniformMatrix4fv(int32(ul), count, transpose, &m[0])
}

//Uniform1i ...
func (ul UniformLocation) Uniform1i(v int32) {
	gl.Uniform1i(int32(ul), v)
}
