// +build linux,arm

package gfx

import (
	log "../log"
    "fmt"
	"errors"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
)

type Texture struct {
	Name string
	Size struct {
		Width  float32
		Height float32
	}

	rgba           *image.RGBA
	texture        uint32
	textureUniform int32
}

func NewTexture(name string) *Texture {
	return &Texture{Name: name}
}


func (texture *Texture) Desc() string {
    ret := ""
    ret += fmt.Sprintf("texture['%s' %.0fx%.0f]",texture.Name,texture.Size.Width,texture.Size.Height)
    return ret
}

func (texture *Texture) LoadFile(path string) error {
	imgFile, err := os.Open(path)
	if err != nil {
		log.Error("fail open %s: %s", path, err)
		return err
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Error("fail decode %s: %s", path, err)
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	err = texture.LoadRGBA(rgba)
	//    log.Debug("read texture file %s",path)
	return err
}

func (texture *Texture) LoadEmpty() error {
	empty := image.NewRGBA(image.Rect(0, 0, 2, 2))
	err := texture.LoadRGBA(empty)
	return err
}

func (texture *Texture) LoadRGBA(rgba *image.RGBA) error {
	texture.rgba = rgba

	//should copy rgba, not point, into the 'loaded' picture here!

	if texture.rgba.Stride != texture.rgba.Rect.Size().X*4 {
		log.Error("fail rgba, unsupported stride")
		return errors.New("unsupported stride")
	}

	texture.Size.Width = float32(texture.rgba.Rect.Size().X)
	texture.Size.Height = float32(texture.rgba.Rect.Size().Y)

	return nil
}

func WhiteColor() *Texture { return ColorTexture("white", image.Black) }
func BlackColor() *Texture { return ColorTexture("black", image.Black) }
func GreyColor() *Texture  { return ColorTexture("gray", image.NewUniform(color.Gray16{0x8000})) }

func ColorTexture(name string, color *image.Uniform) *Texture {
	ret := NewTexture(name)
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 2))
	draw.Draw(rgba, rgba.Bounds(), color, image.ZP, draw.Src)
	ret.LoadRGBA(rgba)
	ret.TexImage()
	return ret
}

func (texture *Texture) Init() {
	gl.GenTextures(1, &texture.texture)
	gl.ActiveTexture(gl.TEXTURE0)
}

func (texture *Texture) BindTexture() {
	gl.BindTexture(gl.TEXTURE_2D, texture.texture)
}

func (texture *Texture) Close() {
	if texture.texture != 0 {
		gl.DeleteTextures(1, &texture.texture)
		texture.texture = 0
	}
}

func (texture *Texture) TexImage() error {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(texture.rgba.Rect.Size().X),
		int32(texture.rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(texture.rgba.Pix))
	//    log.Debug("+tex #%d",texture.texture)
	return nil
}

//rem remove?
func (texture *Texture) Uniform(program *Program) {
	texture.textureUniform, _ = program.Uniform1i(TEXTURE, 0) // REM, should be index not 0 ?
}
