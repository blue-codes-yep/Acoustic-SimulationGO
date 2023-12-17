package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 1600
	windowHeight = 900
)

const (
	numCirclePoints = 75 // Number of points to represent the circle
)

var (
	circleVAO uint32
	circleVBO uint32
)

func loadShaderFile(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkGLError(operation string) {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		log.Printf("OpenGL error after %s: %v\n", operation, errCode)
	}
}

func main() {
	window := openWindow("Acoustic Simulation")
	defer glfw.Terminate()

	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	checkGLError("initializing OpenGL")
	/* Provide debug*/
	gl.DebugMessageCallback(func(
		source uint32,
		gltype uint32,
		id uint32,
		severity uint32,
		length int32,
		message string,
		userParam unsafe.Pointer,
	) {
		log.Println(message)
	}, nil)
	gl.Enable(gl.DEBUG_OUTPUT)

	initCircle()
	// Load shaders
	vertexShaderSource, err := loadShaderFile("shaders/vertexShader.glsl")
	if err != nil {
		log.Fatalf("failed to load vertex shader: %v", err)
	}
	fragmentShaderSource, err := loadShaderFile("shaders/fragmentShader.glsl")
	if err != nil {
		log.Fatalf("failed to load fragment shader: %v", err)
	}

	// Compile and link shaders
	shaderProgram := createProgram(vertexShaderSource, fragmentShaderSource)

	startTime := time.Now()
	waveSpeed := float32(100.0) // Adjust wave speed as needed

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Viewport(0, 0, windowWidth, windowHeight)

		gl.UseProgram(shaderProgram)
		log.Println("Using shader program")
		drawWave(startTime, mgl32.Vec2{0, 0}, waveSpeed)
		checkGLError("after drawing wave")
		// Debugging: Check for OpenGL errors
		if err := gl.GetError(); err != 0 {
			log.Printf("OpenGL error: %v\n", err)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}

	// Clean up
	gl.DeleteProgram(shaderProgram)
}

func openWindow(title string) *glfw.Window {
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

	return window
}

func initCircle() {
	gl.GenVertexArrays(1, &circleVAO)
	gl.GenBuffers(1, &circleVBO)
	gl.BindVertexArray(circleVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, circleVBO)

	var vertices [numCirclePoints * 3]float32 // 3 coordinates per point

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices[:]), gl.DYNAMIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
}

func drawWave(startTime time.Time, waveCenter mgl32.Vec2, waveSpeed float32) {
	elapsed := float32(time.Since(startTime).Seconds())
	waveRadius := elapsed * waveSpeed

	var vertices [numCirclePoints * 3]float32
	for i := 0; i < numCirclePoints; i++ {
		angle := 2 * math.Pi * float64(i) / float64(numCirclePoints)
		x := waveCenter.X() + waveRadius*float32(math.Cos(angle))
		y := waveCenter.Y() + waveRadius*float32(math.Sin(angle))

		ndcX := (x/windowWidth)*2 - 1
		ndcY := (y/windowHeight)*2 - 1

		vertices[i*3] = ndcX
		vertices[i*3+1] = ndcY
		vertices[i*3+2] = 0
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, circleVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices[:]))
	checkGLError("BufferSubData in drawWave")

	gl.BindVertexArray(circleVAO)
	gl.DrawArrays(gl.LINE_LOOP, 0, numCirclePoints)
	checkGLError("DrawArrays in drawWave")
}

func compileShader(source string, shaderType uint32) (uint32, error) {
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

func createProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
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
