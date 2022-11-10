package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/junk1tm/learn-opengl/color"
	"github.com/junk1tm/learn-opengl/object"
	"github.com/junk1tm/learn-opengl/shader"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return err
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	if err := gl.Init(); err != nil {
		return err
	}

	program, err := shader.NewProgramFromFile("vertex.glsl", "fragment.glsl")
	if err != nil {
		return err
	}
	defer program.Delete()

	program.Use()

	// indices := []uint32{
	// 	0, 1, 3, // first triangle
	// 	1, 2, 3, // second triangle
	// }
	obj := object.New(
		[]object.Vertex{
			// positions(3) + colors(3)
			{0.5, -0.5, 0, 1, 0, 0},
			{-0.5, -0.5, 0, 0, 1, 0},
			{0, 0.5, 0, 0, 0, 1},
		},
		object.WithAttribute(3, false), // color
		// WithIndices(indices),
	)

	var input inputHandler

	for !window.ShouldClose() {
		input.handle(window)
		color.Clear(color.DarkGreen)

		// time := glfw.GetTime()
		// green := math.Sin(time)/2 + 0.5
		if err := program.SetUniformFloat("shift", input.shift); err != nil {
			return err
		}

		obj.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}

	return nil
}

type inputHandler struct {
	shift float32
}

func (i *inputHandler) handle(window *glfw.Window) {
	switch {
	case window.GetKey(glfw.KeyEscape) == glfw.Press:
		window.SetShouldClose(true)

	case window.GetKey(glfw.Key1) == glfw.Press:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	case window.GetKey(glfw.Key2) == glfw.Press:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	case window.GetKey(glfw.KeyRight) == glfw.Press:
		i.shift += 0.001

	case window.GetKey(glfw.KeyLeft) == glfw.Press:
		i.shift -= 0.001
	}
}
