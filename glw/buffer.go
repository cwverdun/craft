package glw

import "github.com/go-gl/gl/v4.5-core/gl"

//Buffer the gl data buffer
type Buffer uint32

//NewBuffer creates new data buffer
func NewBuffer(data []float32) Buffer {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	b := Buffer(vbo)
	b.BindBuffer()
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(&data[0]), gl.STATIC_DRAW)
	return b
}

//BindBuffer ...
func (b Buffer) BindBuffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(b))
}

//Delete deletes the buffer
func (b Buffer) Delete() {
	if b == 0 {
		return
	}
	vbo := uint32(b)
	gl.DeleteBuffers(1, &vbo)
}
