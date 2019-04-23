
package gfx


import (
    "image"
    log "../log"
    
)

type Texture struct {
    Name string
    Size struct{Width float32; Height float32}
    
    rgba *image.RGBA
    texture uint32 
    textureUniform int32
}


func NewTexture(name string) *Texture {
    return &Texture{Name: name}    
}


func (texture *Texture) Close() { }
func (texture *Texture) LoadEmpty() error { return log.NewError("TEXTURE NOT AVAILABLE") }
func (texture *Texture) LoadRGBA(rgba *image.RGBA) error { return log.NewError("TEXTURE NOT AVAILABLE") }
func (texture *Texture) TexImage() error { return log.NewError("TEXTURE NOT AVAILABLE") }
