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
