package shader

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	id uint32
}

func New(kind uint32, r io.Reader) (*Shader, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read shader: %w", err)
	}
	b = terminateB(b)

	src, free := gl.Strs(string(b))
	defer free()

	id := gl.CreateShader(kind)
	gl.ShaderSource(id, 1, src, nil)

	var ok int32
	gl.CompileShader(id)
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &ok)
	if ok == gl.FALSE {
		var size int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &size)
		log := make([]byte, size)
		gl.GetShaderInfoLog(id, size, nil, &log[0])
		return nil, fmt.Errorf("compile shader: %s", log)
	}

	return &Shader{id: id}, nil
}

func NewFromFile(kind uint32, name string) (*Shader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	return New(kind, f)
}

func (s *Shader) Delete() { gl.DeleteShader(s.id) }

type Program struct {
	id uint32
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	id := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(id, shader.id)
	}

	var ok int32
	gl.LinkProgram(id)
	gl.GetProgramiv(id, gl.LINK_STATUS, &ok)
	if ok == gl.FALSE {
		var size int32
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &size)
		log := make([]byte, size)
		gl.GetProgramInfoLog(id, size, nil, &log[0])
		return nil, fmt.Errorf("link program: %s", log)
	}

	for _, shader := range shaders {
		shader.Delete()
	}

	return &Program{id: id}, nil
}

func NewProgramFromFile(vertexFile, fragmentFile string) (*Program, error) {
	vertex, err := NewFromFile(gl.VERTEX_SHADER, vertexFile)
	if err != nil {
		return nil, fmt.Errorf("create vertex: %w", err)
	}

	fragment, err := NewFromFile(gl.FRAGMENT_SHADER, fragmentFile)
	if err != nil {
		return nil, fmt.Errorf("create fragment: %w", err)
	}

	return NewProgram(vertex, fragment)
}

func (p *Program) Use()    { gl.UseProgram(p.id) }
func (p *Program) Delete() { gl.DeleteProgram(p.id) }

func (p *Program) SetUniformInt(name string, vs ...int32) error {
	location, err := p.location(name)
	if err != nil {
		return err
	}

	switch len(vs) {
	case 1:
		gl.ProgramUniform1i(p.id, location, vs[0])
	case 2:
		gl.ProgramUniform2i(p.id, location, vs[0], vs[1])
	case 3:
		gl.ProgramUniform3i(p.id, location, vs[0], vs[1], vs[2])
	case 4:
		gl.ProgramUniform4i(p.id, location, vs[0], vs[1], vs[2], vs[3])
	default:
		return errors.New("values count exceeds 4")
	}

	return nil
}

func (p *Program) SetUniformFloat(name string, vs ...float32) error {
	location, err := p.location(name)
	if err != nil {
		return err
	}

	switch len(vs) {
	case 1:
		gl.ProgramUniform1f(p.id, location, vs[0])
	case 2:
		gl.ProgramUniform2f(p.id, location, vs[0], vs[1])
	case 3:
		gl.ProgramUniform3f(p.id, location, vs[0], vs[1], vs[2])
	case 4:
		gl.ProgramUniform4f(p.id, location, vs[0], vs[1], vs[2], vs[3])
	default:
		return errors.New("values count exceeds 4")
	}

	return nil
}

// TODO: get locations once
func (p *Program) location(name string) (int32, error) {
	name = terminate(name)

	location := gl.GetUniformLocation(p.id, gl.Str(name))
	if location == -1 {
		return 0, fmt.Errorf("uniform %s not found", name)
	}

	return location, nil
}

func terminate(s string) string {
	return string(terminateB([]byte(s)))
}

func terminateB(b []byte) []byte {
	if last := len(b) - 1; b[last] == 0 {
		return b
	}
	return append(b, 0)
}
