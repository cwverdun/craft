package glw

import "github.com/go-gl/gl/v4.5-core/gl"

//VertexArray ...
type VertexArray uint32

//NewVertexArray ...
func NewVertexArray() VertexArray {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	return VertexArray(vao)
}

//BindVertexArray ...
func (va VertexArray) BindVertexArray() {
	gl.BindVertexArray(uint32(va))
}

//Delete deletes the VertexArray
func (va VertexArray) Delete() {
	vao := uint32(va)
	gl.DeleteVertexArrays(1, &vao)
}
