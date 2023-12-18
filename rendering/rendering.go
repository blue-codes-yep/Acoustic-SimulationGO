// rendering/rendering.go
package rendering

import (
	"log"
	"math"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 1600
	windowHeight = 900
)

const (
	numCirclePoints = 100 // Number of points to represent the circle
)

var (
	circleVAO uint32
	circleVBO uint32
)

func checkGLError(operation string) {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		log.Printf("OpenGL error after %s: %v\n", operation, errCode)
	}
}

func InitCircle() {
	gl.GenVertexArrays(1, &circleVAO)
	gl.GenBuffers(1, &circleVBO)
	gl.BindVertexArray(circleVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, circleVBO)

	var vertices [numCirclePoints * 3]float32 // 3 coordinates per point

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices[:]), gl.DYNAMIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
}

func DrawWave(startTime time.Time, waveCenter mgl32.Vec2, waveSpeed float32) {
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

func RenderWave(shaderProgram uint32, currentTime float64) {
	maxX, maxY, maxZ := 10.0, 20.0, 40.0
	dx, dz := 1.0, 1.0 // Increments for quads
	amplitude := 1.0
	frequency := 150.0
	wavelength := 1.0

	gl.UseProgram(shaderProgram)

	for x := 0.0; x < maxX-dx; x += dx {
		for z := 0.0; z < maxZ-dz; z += dz {
			// Calculate base Y position and neighboring offsets
			y := float32(amplitude * math.Sin(2.0*math.Pi*frequency*(x+currentTime)/wavelength))
			y2 := float32(amplitude * math.Sin(2.0*math.Pi*frequency*(x+dx+currentTime)/wavelength))
			y3 := float32(amplitude * math.Sin(2.0*math.Pi*frequency*(x+dx+dz+currentTime)/wavelength))
			y4 := float32(amplitude * math.Sin(2.0*math.Pi*frequency*(x+dz+currentTime)/wavelength))

			// Clamp Y values to limits
			clampY(y, float32(maxY))
			clampY(y2, float32(maxY))
			clampY(y3, float32(maxY))
			clampY(y4, float32(maxY))

			// Create four vertices for the quad
			vertices := []float32{
				float32(x), y, float32(z),
				float32(x + dx), y2, float32(z),
				float32(x + dx), y3, float32(z + dz),
				float32(x), y4, float32(z + dz),
			}

			// Bind and write vertex data
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))

			// Draw the triangle formed by the first three vertices
			gl.DrawArrays(gl.TRIANGLES, 0, 3)

			// Draw the triangle formed by the last three vertices (connects back to the beginning)
			gl.DrawArrays(gl.TRIANGLES, 3, 3)
		}
	}

	checkGLError("RenderWave")
}

func clampY(y float32, maxY float32) {
	if y < 0.0 {
		y = 0.0
	}
}

func vertexColor(y float32, maxY float32) []float32 {
	// Normalize Y position to range between 0 and 1
	normalizedY := y / maxY
	return []float32{0.0, 0.0, normalizedY, 1.0} // Blue to white gradient
}
