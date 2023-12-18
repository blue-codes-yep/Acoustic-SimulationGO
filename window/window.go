// window/window.go
package window

import "github.com/go-gl/glfw/v3.3/glfw"

const (
	windowWidth  = 1600
	windowHeight = 900
)

func OpenWindow(title string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	println("Window created")
	return window
}
