// +build linux,arm

package gfx

import (
	//    "strings"
	"fmt"
	//    log "../log"
	gl "github.com/FEEDFACE-COM/piglet/gles2"
)

type Object struct {
	Name   string
	object uint32
}

func NewObject(name string) *Object {
	ret := &Object{Name: name}
	return ret
}

func (object *Object) Init() error {
	//TODO: check for cleanup?
	gl.GenBuffers(1, &object.object)
	return nil
}

func (object *Object) BindBuffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, object.object)
}

func (object *Object) BufferData(size int, value []float32) {
    if size <= 0 {
        return
    }
	gl.BindBuffer(gl.ARRAY_BUFFER, object.object)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(value), gl.STATIC_DRAW)
}

func (object *Object) Desc() string {
	return fmt.Sprintf("object[%s]", object.Name)
}
