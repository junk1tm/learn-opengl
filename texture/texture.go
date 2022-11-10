package texture

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	id uint32
}

func New(r io.Reader) (*Texture, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Src)

	// coords := []float32{
	// 	0.0, 0.0, // lower-left corner
	// 	1.0, 0.0, // lower-right corner
	// 	0.5, 1.0, // top-center corner
	// }

	var texture Texture
	gl.GenTextures(1, &texture.id)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	width, height := int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy())
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(rgba))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return &texture, nil
}

func NewFromFile(name string) (*Texture, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	return New(f)
}
