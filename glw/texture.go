package glw

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" //for decoding jpeg files
	_ "image/png"  //for decoding png files
	"io"
	"os"

	"github.com/go-gl/gl/v4.5-core/gl"
)

//Texture the gl texture
type Texture uint32

//NewTexture creates new texture from img file
func NewTexture(r io.Reader) (Texture, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return 0, fmt.Errorf("cannot decode image: %s", err)
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	switch trueim := img.(type) {
	case *image.RGBA:
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGBA,
			int32(trueim.Bounds().Dx()), int32(trueim.Bounds().Dy()),
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(trueim.Pix),
		)
	default:
		copy := image.NewRGBA(trueim.Bounds())
		draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGBA,
			int32(copy.Bounds().Dx()), int32(copy.Bounds().Dy()),
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(copy.Pix),
		)
	}

	return Texture(texture), nil
}

//NewTextureFromFile loads texture from filepath
func NewTextureFromFile(filepath string) (Texture, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	texture, err := NewTexture(f)
	if err != nil {
		return 0, err
	}
	return texture, nil
}

//BindTexture ...
func (t Texture) BindTexture() {
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}
