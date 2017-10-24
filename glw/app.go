package glw

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//Camera the world camera
type Camera struct {
	Pos         mgl32.Vec3
	Rotation    mgl32.Vec3
	Up          mgl32.Vec3
	Perspective mgl32.Mat4
	Model       mgl32.Mat4
	ProjView    mgl32.Mat4
}

//Update recalculates MVP
func (cam *Camera) Update() {
	view := mgl32.LookAtV(cam.Pos, cam.Rotation, cam.Up)
	cam.ProjView = cam.Perspective.Mul4(view).Mul4(cam.Model)
}

//App wraps the window creation and gl initialization and the main loop
type App struct {
	window    *glfw.Window
	Renderers []Renderer
	Camera    *Camera

	horizontalAngle float64
	verticalAngle   float64
	mouseSpeed      float64
	speed           float64
	lastTime        float64
}

//NewApp creates window
func NewApp(title string, width, height int, cam *Camera) (*App, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %s", err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Cube", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating window: %s", err)
	}
	window.MakeContextCurrent()
	// Initialize Glow
	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("error initializing gl: %s", err)
	}
	return &App{window: window, Camera: cam}, nil
}

//Terminate terminates glfw
func (app *App) Terminate() {
	glfw.Terminate()
}

//Run runs the main loop every time it runs trough renderers
func (app *App) Run() error {
	// Configure global settings
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	previousTime := glfw.GetTime()

	for !app.window.ShouldClose() {

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		for _, r := range app.Renderers {
			r.Render(app.Camera.ProjView, elapsed)
		}

		// Maintenance
		app.window.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

func (app *App) SetCameraCursor(mouseSpeed, speed float64) {
	app.mouseSpeed = mouseSpeed
	app.speed = speed
	app.lastTime = glfw.GetTime()
	app.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	app.window.SetCursorPosCallback(app.cursor)
	app.window.SetKeyCallback(app.key)
}

func (app *App) cursor(_ *glfw.Window, xpos, ypos float64) {
	w, h := app.window.GetSize()
	app.window.SetCursorPos(float64(w)/2, float64(h)/2)
	app.horizontalAngle += app.mouseSpeed * (float64(w)/2 - xpos)
	app.verticalAngle += app.mouseSpeed * (float64(h)/2 - ypos)

	direction := mgl32.Vec3{
		float32(math.Cos(app.verticalAngle) * math.Sin(app.horizontalAngle)),
		float32(math.Sin(app.verticalAngle)),
		float32(math.Cos(app.verticalAngle) * math.Cos(app.horizontalAngle)),
	}

	right := mgl32.Vec3{
		float32(math.Sin(app.horizontalAngle - math.Pi/2)),
		0,
		float32(math.Cos(app.horizontalAngle - math.Pi/2)),
	}

	up := right.Cross(direction)

	app.Camera.Rotation = app.Camera.Pos.Add(direction)
	app.Camera.Up = up
	app.Camera.Update()
}

func (app *App) key(_ *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyW:
	case glfw.KeyA:
	case glfw.KeyS:
	case glfw.KeyD:
	default:
		return
	}
	currentTime := glfw.GetTime()
	deltaTime := currentTime - app.lastTime

	direction := mgl32.Vec3{
		float32(math.Cos(app.verticalAngle) * math.Sin(app.horizontalAngle)),
		float32(math.Sin(app.verticalAngle)),
		float32(math.Cos(app.verticalAngle) * math.Cos(app.horizontalAngle)),
	}

	right := mgl32.Vec3{
		float32(math.Sin(app.horizontalAngle - math.Pi/2)),
		0,
		float32(math.Cos(app.horizontalAngle - math.Pi/2)),
	}

	switch key {
	case glfw.KeyW:
		app.Camera.Pos = app.Camera.Pos.Add(direction.Mul(float32(deltaTime * app.speed)))
	case glfw.KeyS:
		app.Camera.Pos = app.Camera.Pos.Sub(direction.Mul(float32(deltaTime * app.speed)))
	case glfw.KeyA:
		app.Camera.Pos = app.Camera.Pos.Sub(right.Mul(float32(deltaTime * app.speed)))
	case glfw.KeyD:
		app.Camera.Pos = app.Camera.Pos.Add(right.Mul(float32(deltaTime * app.speed)))
	default:
		return
	}
	up := right.Cross(direction)
	app.Camera.Rotation = app.Camera.Pos.Add(direction)
	app.Camera.Up = up
	app.Camera.Update()

	app.lastTime = currentTime
}
