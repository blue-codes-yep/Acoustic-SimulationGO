package main

import (
	"fmt"
	"go_physics_engine/rendering"
	"go_physics_engine/window"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func checkGLError(operation string) {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		log.Printf("OpenGL error after %s: %v\n", operation, errCode)
	}
}
func LoadShaderFile(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read shader file %s: %v\n", filePath, err)
		return "", err
	}
	return string(fileData), nil
}

func main() {
	println("Main")
	win := window.OpenWindow("Acoustic Simulation")
	println("Window opened")
	defer glfw.Terminate()

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		log.Printf("OpenGL error after initialization: %v\n", errCode)
	}

	rendering.InitCircle()

	// Load shaders
	vertexShaderSource, err := LoadShaderFile("shader/vertexShader.glsl")
	if err != nil {
		log.Fatalf("failed to load vertex shader: %v", err)
	}
	fragmentShaderSource, err := LoadShaderFile("shader/fragmentShader.glsl")
	if err != nil {
		log.Fatalf("failed to load fragment shader: %v", err)
	}
	println("Shaders loaded")
	shaderProgram := CreateProgram(vertexShaderSource, fragmentShaderSource)
	println("Shader program created")
	startTime := time.Now()
	waveSpeed := float32(100.0) // Adjust wave speed as needed

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Get the current size of the window
		width, height := win.GetSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		gl.UseProgram(shaderProgram)
		rendering.DrawWave(startTime, mgl32.Vec2{0, 0}, waveSpeed)

		win.SwapBuffers()
		glfw.PollEvents()
	}

	// Clean up
	gl.DeleteProgram(shaderProgram)
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	// Convert the Go string into a null-terminated C string
	csource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csource, nil)
	free() // This frees the C string memory

	gl.CompileShader(shader)
	checkGLError("compiling shader")
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		shaderTypeName := "Unknown"
		if shaderType == gl.VERTEX_SHADER {
			shaderTypeName = "Vertex Shader"
		} else if shaderType == gl.FRAGMENT_SHADER {
			shaderTypeName = "Fragment Shader"
		}

		return 0, fmt.Errorf("%s compilation error: %v\nShader Source:\n%s", shaderTypeName, log, source)
	}

	return shader, nil
}

func CreateProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	checkGLError("linking shader program")
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to link program: %v", log))

	}

	return program
}
