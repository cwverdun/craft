package main

import (
	"bytes"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/microo8/craft/glw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const (
	windowWidth  = 1920
	windowHeight = 1024
)

func main() {
	cam := &glw.Camera{
		Pos:         mgl32.Vec3{100, 50, 100},
		Rotation:    mgl32.Vec3{0, 0, 0},
		Up:          mgl32.Vec3{0, 1, 0},
		Model:       mgl32.Ident4(),
		Perspective: mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 3000.0),
	}
	cam.Update()
	app, err := glw.NewApp("Cube", windowWidth, windowHeight, cam)
	if err != nil {
		panic(err)
	}
	defer app.Terminate()

	rr := newWorld()
	app.Renderers = append(app.Renderers, rr)

	app.SetCameraCursor(0.001, 5)

	if err = app.Run(); err != nil {
		panic(err)
	}
}

type world struct {
	p       glw.Program
	texture glw.Texture

	mvp  glw.UniformLocation
	vert glw.VertexAttrib
	uv   glw.VertexAttrib

	chunks []*Chunk
	player *Player
}

func newWorld() *world {
	w := new(world)
	w.player = &Player{}
	vs := glw.ShaderSource{Type: gl.VERTEX_SHADER, Source: bytes.NewBufferString(vertexShaderSrc)}
	fs := glw.ShaderSource{Type: gl.FRAGMENT_SHADER, Source: bytes.NewBufferString(fragmentShaderSrc)}
	p, err := glw.NewProgram(vs, fs)
	if err != nil {
		panic(err)
	}
	w.p = p
	w.p.UseProgram()

	texture, err := glw.NewTextureFromFile("assets/textures/texture.png")
	if err != nil {
		log.Fatalln(err)
	}
	w.texture = texture

	// Configure the vertex data
	glw.NewVertexArray().BindVertexArray()

	w.mvp = w.p.GetUniformLocation("mvp")
	w.vert = w.p.GetAttribLocation("vert")
	w.uv = w.p.GetAttribLocation("uv")

	return w
}

//Render rotates the camera around the cube
func (w *world) Render(projView mgl32.Mat4, elapsed float64) {
	p := round(w.player.X / chunkSize)
	q := round(w.player.Y / chunkSize)
	//delete far chunks
	var deleteIndices []int
	for i, chunk := range w.chunks {
		dp := chunk.P - p
		dq := chunk.Q - q
		if abs(dp) >= chunkDeleteRadius && abs(dq) >= chunkDeleteRadius {
			deleteIndices = append(deleteIndices, i)
		}
	}
	for _, i := range deleteIndices {
		w.chunks[i].Delete()
		w.chunks = append(w.chunks[:i], w.chunks[i+1:]...)
	}
	//create new chunks in render radius
	for i := -chunkRenderRadius; i <= chunkRenderRadius; i++ {
		for j := -chunkRenderRadius; j <= chunkRenderRadius; j++ {
			if w.getChunk(i, j) != nil {
				continue
			}
			w.chunks = append(w.chunks, NewChunk(i, j))
		}
	}
	//Render
	gl.ClearColor(0.53, 0.81, 0.92, 1.00)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	w.mvp.UniformMatrix4fv(1, false, projView)

	for _, chunk := range w.chunks {
		dp := chunk.P - p
		dq := chunk.Q - q
		if abs(dp) <= chunkRenderRadius && abs(dq) <= chunkRenderRadius {
			chunk.Draw(w.vert, w.uv)
		}
	}
}

func (w *world) getChunk(p, q int) *Chunk {
	for _, chunk := range w.chunks {
		if chunk.P == p && chunk.Q == q {
			return chunk
		}
	}
	return nil
}

type Player struct {
	X, Y, Z float32
}

var vertexShaderSrc = `
#version 330
uniform mat4 mvp;
in vec3 vert;
in vec2 uv;
out vec2 fragUV;
void main() {
    fragUV = uv;
    gl_Position = mvp * vec4(vert, 1);
}
`

var fragmentShaderSrc = `
#version 330
uniform sampler2D tex;
in vec2 fragUV;
out vec4 color;
void main() {
    color = texture(tex, fragUV);
}
`
