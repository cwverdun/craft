package glw

import "github.com/go-gl/mathgl/mgl32"

//Renderer interface is represents an renderer that will be used in the main loop of the application
type Renderer interface {
	Render(projView mgl32.Mat4, elapsed float64)
}
