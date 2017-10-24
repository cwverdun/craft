package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/microo8/meh/opengl/craft/glw"
	"github.com/microo8/meh/opengl/craft/noise"
)

const (
	chunkSize         = 16
	chunkRenderRadius = 8
	chunkDeleteRadius = 12
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func round(val float32) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

//Chunk ...
type Chunk struct {
	P            int
	Q            int
	vertexBuffer glw.Buffer
	uvBuffer     glw.Buffer
	faces        int

	m [chunkSize][256][chunkSize]*Block
}

//NewChunk creates new chunk
func NewChunk(p, q int) *Chunk {
	chunk := &Chunk{P: p, Q: q}
	for dx := 0; dx < chunkSize; dx++ {
		for dz := 0; dz < chunkSize; dz++ {
			x := chunk.P*chunkSize + dx
			z := chunk.Q*chunkSize + dz
			f := noise.Simplex2(float64(x)*0.01, float64(z)*0.01, 4, 0.5, 2)
			g := noise.Simplex2(float64(x)*0.01, float64(z)*0.01, 2, 0.9, 2)
			mh := g*32 + 16
			h := f * mh
			w := DirtItem
			t := 12
			if h < float64(t) {
				h = float64(t - 1)
				w = SandItem
			}
			for y := 0; y < 256; y++ {
				if y < int(h) {
					chunk.m[dx][y][dz] = &Block{t: w}
				} else {
					chunk.m[dx][y][dz] = &Block{t: EmptyItem}
				}
			}
		}
	}
	chunk.genBuffers()
	return chunk
}

func (chunk *Chunk) genBuffers() {
	chunk.faces = 0
	for i := 0; i < chunkSize; i++ {
		for j := 0; j < 256; j++ {
			for k := 0; k < chunkSize; k++ {
				chunk.faces += int(chunk.m[i][j][k].CountExposedFaces(chunk, i, j, k))
			}
		}
	}
	if chunk.faces == 0 {
		return
	}
	vertexBuffer := make([]float32, chunk.faces*6*3)
	uvBuffer := make([]float32, chunk.faces*6*2)
	var vbOffset, uvOffset int
	for x := 0; x < chunkSize; x++ {
		for y := 0; y < 256; y++ {
			for z := 0; z < chunkSize; z++ {
				b := chunk.m[x][y][z]
				if b.faces == 0 {
					continue
				}
				b.MakeCube(chunk, x, y, z, vertexBuffer[vbOffset:], uvBuffer[uvOffset:])
				vbOffset += int(b.faces) * 6 * 3
				uvOffset += int(b.faces) * 6 * 2
			}
		}
	}
	chunk.vertexBuffer = glw.NewBuffer(vertexBuffer)
	chunk.uvBuffer = glw.NewBuffer(uvBuffer)
	return
}

func (chunk *Chunk) Draw(vert, uv glw.VertexAttrib) {
	chunk.vertexBuffer.BindBuffer()
	vert.EnableVertexAttribArray()
	vert.VertexAttribPointer(3, gl.FLOAT, false, 0, nil)

	chunk.uvBuffer.BindBuffer()
	uv.EnableVertexAttribArray()
	uv.VertexAttribPointer(2, gl.FLOAT, false, 0, nil)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(6*2*3*chunk.faces))
}

//Delete deletes the buffers of the chunk
func (chunk *Chunk) Delete() {
	chunk.vertexBuffer.Delete()
	chunk.uvBuffer.Delete()
}

type Block struct {
	t     itemType
	faces int16
	f     [6]bool
}

func (b *Block) CountExposedFaces(chunk *Chunk, x, y, z int) (faces int16) {
	if b.t == EmptyItem {
		b.faces = 0
		return
	}
	if b.f[0] = y == 0 || chunk.m[x][y-1][z].t == EmptyItem; b.f[0] { //Bottom
		faces++
	}
	if b.f[1] = y == 255 || chunk.m[x][y+1][z].t == EmptyItem; b.f[1] { //Top
		faces++
	}
	if b.f[2] = z == chunkSize-1 || chunk.m[x][y][z+1].t == EmptyItem; b.f[2] { //Front
		faces++
	}
	if b.f[3] = z == 0 || chunk.m[x][y][z-1].t == EmptyItem; b.f[3] { //Back
		faces++
	}
	if b.f[4] = x == 0 || chunk.m[x-1][y][z].t == EmptyItem; b.f[4] { //Left
		faces++
	}
	if b.f[5] = x == chunkSize-1 || chunk.m[x+1][y][z].t == EmptyItem; b.f[5] { //Right
		faces++
	}
	b.faces = faces
	return
}

func (b *Block) MakeCube(chunk *Chunk, x, y, z int, vertexBuffer, uvBuffer []float32) {
	var vbOffset, uvOffset int
	for f := 0; f < 6; f++ {
		if !b.f[f] {
			continue
		}
		for i := 0; i < 6; i++ {
			vertexBuffer[vbOffset+i*3+0] = cubeVertices[f*6*3+i*3+0] + float32(x) + (float32(chunk.P) * chunkSize)
			vertexBuffer[vbOffset+i*3+1] = cubeVertices[f*6*3+i*3+1] + float32(y)
			vertexBuffer[vbOffset+i*3+2] = cubeVertices[f*6*3+i*3+2] + float32(z) + (float32(chunk.Q) * chunkSize)

			uvBuffer[uvOffset+i*2+0] = uvs[f*6*2+i*2+0] + (texWidth * float32(b.t-1))
			uvBuffer[uvOffset+i*2+1] = uvs[f*6*2+i*2+1]
		}
		vbOffset += 6 * 3
		uvOffset += 6 * 2
	}
}

type itemType int16

const (
	EmptyItem itemType = iota
	DirtItem
	SandItem
	StoneItem
	BrickItem
)

const (
	itemsCount = 8
	texWidth   = 1 / float32(itemsCount)
	texHeight  = 1 / float32(3)
)

var uvs = []float32{
	//Bottom
	0.0, 2 * texHeight,
	texWidth, 2 * texHeight,
	0.0, 1.0,
	texWidth, 2 * texHeight,
	texWidth, 1.0,
	0.0, 1.0,
	//Top
	0.0, 0.0,
	0.0, texHeight,
	texWidth, 0.0,
	texWidth, 0.0,
	0.0, texHeight,
	texWidth, texHeight,
	//Front
	0.0, 2 * texHeight,
	texWidth, 2 * texHeight,
	0.0, texHeight,
	texWidth, 2 * texHeight,
	texWidth, texHeight,
	0.0, texHeight,
	//Back
	0.0, 2 * texHeight,
	0.0, texHeight,
	texWidth, 2 * texHeight,
	texWidth, 2 * texHeight,
	0.0, texHeight,
	texWidth, texHeight,
	//Left
	texWidth, 2 * texHeight,
	0.0, texHeight,
	0.0, 2 * texHeight,
	texWidth, 2 * texHeight,
	texWidth, texHeight,
	0.0, texHeight,
	//Right
	0.0, 2 * texHeight,
	texWidth, 2 * texHeight,
	texWidth, texHeight,
	0.0, 2 * texHeight,
	texWidth, texHeight,
	0.0, texHeight,
}

var cubeVertices = []float32{
	// Bottom
	0.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	0.0, 0.0, 1.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 1.0,
	0.0, 0.0, 1.0,
	// Top
	0.0, 1.0, 0.0,
	0.0, 1.0, 1.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	0.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	// Front
	0.0, 0.0, 1.0,
	1.0, 0.0, 1.0,
	0.0, 1.0, 1.0,
	1.0, 0.0, 1.0,
	1.0, 1.0, 1.0,
	0.0, 1.0, 1.0,
	// Back
	0.0, 0.0, 0.0,
	0.0, 1.0, 0.0,
	1.0, 0.0, 0.0,
	1.0, 0.0, 0.0,
	0.0, 1.0, 0.0,
	1.0, 1.0, 0.0,
	// Left
	0.0, 0.0, 1.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 0.0,
	0.0, 0.0, 1.0,
	0.0, 1.0, 1.0,
	0.0, 1.0, 0.0,
	// Right
	1.0, 0.0, 1.0,
	1.0, 0.0, 0.0,
	1.0, 1.0, 0.0,
	1.0, 0.0, 1.0,
	1.0, 1.0, 0.0,
	1.0, 1.0, 1.0,
}
