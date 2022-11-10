package object

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	uint32size  = 4
	float32size = 4
)

type Vertex []float32

type Object struct {
	vao uint32
	vbo uint32
	ebo uint32
	// attributes:
	index  uint32
	stride int32
	offset int
	// length:
	verticesLen int32
	indicesLen  int
}

func New(vertices []Vertex, opts ...Option) *Object {
	const coordsSize = 3

	obj := Object{
		index:       0,
		stride:      int32(len(vertices[0])) * float32size,
		offset:      0,
		verticesLen: int32(len(vertices)),
		indicesLen:  0,
	}
	gl.GenVertexArrays(1, &obj.vao)
	gl.BindVertexArray(obj.vao)
	defer gl.BindVertexArray(0)

	var data []float32
	for _, v := range vertices {
		data = append(data, v...)
	}

	gl.GenBuffers(1, &obj.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, obj.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*float32size, gl.Ptr(data), gl.STATIC_DRAW)

	// coordinates:
	WithAttribute(coordsSize, false)(&obj)
	for _, opt := range opts {
		opt(&obj)
	}

	return &obj
}

type Option func(*Object)

func WithIndices(indices []uint32) Option {
	return func(obj *Object) {
		obj.indicesLen = len(indices)
		gl.GenBuffers(1, &obj.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, obj.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*uint32size, gl.Ptr(indices), gl.STATIC_DRAW)
	}
}

func WithAttribute(size int, normalized bool) Option {
	return func(obj *Object) {
		gl.EnableVertexAttribArray(obj.index)
		gl.VertexAttribPointer(obj.index, int32(size), gl.FLOAT, normalized, obj.stride, gl.PtrOffset(obj.offset*float32size))
		obj.index++
		obj.offset += size
	}
}

func (o *Object) Draw() {
	gl.BindVertexArray(o.vao)
	defer gl.BindVertexArray(0)

	if o.ebo != 0 {
		gl.DrawElements(gl.TRIANGLES, int32(o.indicesLen), gl.UNSIGNED_INT, gl.PtrOffset(0))
		return
	}

	gl.DrawArrays(gl.TRIANGLES, 0, o.verticesLen)
}

func (o *Object) Delete() {
	gl.DeleteVertexArrays(1, &o.vao)
	gl.DeleteBuffers(1, &o.vbo)
	gl.DeleteBuffers(1, &o.ebo)
}
