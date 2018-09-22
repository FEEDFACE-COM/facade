
// +build linux,arm

package gfx

import (
    "fmt"
    conf "../conf"
//    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"    
)

type Camera struct {
    projection mgl32.Mat4
    projectionUniform int32
    
    camera mgl32.Mat4
    cameraUniform int32
    
    Width float32
    Height float32
    Isometric bool
}


func NewCamera(config *conf.CameraConfig, width,height float32) *Camera {
    ret := &Camera{Width: width, Height: height}
    ret.Isometric = config.Isometric

    ratio := height/width
    
    l := float32(-1.)
    r := float32(1.)
    b := float32(-ratio)
    t := float32(ratio)
    n := float32(-4096.)
    f := float32(4096.)
    
    orthographic := mgl32.Ortho( l,r,b,t,n,f )
    perspective := mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.1, 10.0)
    
    if ret.Isometric {
        ret.projection = orthographic
    } else {
        ret.projection = perspective
    }
        
    
    ret.camera = mgl32.LookAtV(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
    return ret
}

func (camera *Camera) Uniform(program uint32) {
	camera.projectionUniform = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(camera.projectionUniform, 1, false, &camera.projection[0])

	camera.cameraUniform = gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(camera.cameraUniform, 1, false, &camera.camera[0])
    
}


func (camera *Camera) Configure(config *conf.CameraConfig) {
    
}


func (camera *Camera) Desc() string {
    tmp := ""
    if camera.Isometric { tmp = " iso" }
    return fmt.Sprintf("camera[%5.0fx%5.0f%s]",camera.Width,camera.Height,tmp) 
}