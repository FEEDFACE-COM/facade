// +build !linux !arm

package gfx

import (
	"FEEDFACE.COM/facade/log"
	"image"
)

type Texture struct{ Name string }

func NewTexture(name string) *Texture {
	return &Texture{Name: name}
}

func (texture *Texture) Close()                          {}
func (texture *Texture) LoadEmpty() error                { return log.NewError("TEXTURE NOT AVAILABLE") }
func (texture *Texture) LoadRGBA(rgba *image.RGBA) error { return log.NewError("TEXTURE NOT AVAILABLE") }
func (texture *Texture) TexImage() error                 { return log.NewError("TEXTURE NOT AVAILABLE") }

func (texture *Texture) Desc() string { return "texture[]" }
