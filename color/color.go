package color

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Color struct {
	R, G, B, A float32
}

var DarkGreen = Color{R: 0.2, G: 0.3, B: 0.3, A: 1.0}

func Clear(c Color) {
	gl.ClearColor(c.R, c.G, c.B, c.A)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
