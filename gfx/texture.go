
// +build linux,arm


package gfx

import (
    "os"
    "errors"
    "image"
    "image/draw"
    _ "image/png"
    gl "src.feedface.com/gfx/piglet/gles2"
    log "../log"
)

type Texture struct {
    path string
    rgba *image.RGBA
    texture uint32 
}


func NewTexture() *Texture {
    return &Texture{}    
}

func (texture *Texture) Close() {

    log.Debug("~tex") 
    if texture.texture != 0 {
        gl.DeleteTextures(1, &texture.texture)    
    }   
    texture.texture = 0
}



func (texture *Texture) LoadFile(path string) error {
    imgFile, err := os.Open(path)
    if err != nil {
        log.Error("fail open %s: %s",path,err)
        return err
    }
    
    img, _, err := image.Decode(imgFile)
    if err != nil {
        log.Error("fail decode %s: %s",path,err)
        return err
    }
    
    texture.rgba = image.NewRGBA(img.Bounds())
    if texture.rgba.Stride != texture.rgba.Rect.Size().X * 4 {
        log.Error("fail rgba, unsupported stride")
        return errors.New("unsupported stride")    
    }
    draw.Draw(texture.rgba, texture.rgba.Bounds(), img, image.Point{0,0}, draw.Src)
    
    return nil
}

func (texture *Texture) genTexture() error {
    gl.GenTextures(1, &texture.texture)
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
		gl.Ptr(texture.rgba.Pix) )
    return nil
}


