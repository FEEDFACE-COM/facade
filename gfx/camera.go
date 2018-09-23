
// +build linux,arm

package gfx

import (
    "fmt"
    conf "../conf"
//    log "../log"
//    gl "src.feedface.com/gfx/piglet/gles2"
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


func orthographic(width,height float32) mgl32.Mat4 {
    const F = 5.
    var w float32 = F * (width/height) /2.
    var h float32 = F * (height/height) /2.
    var l float32 = F * 1000. /2.
    return mgl32.Ortho( -w, w, -h, h, -l, l )
}

func perspective(width,height float32) mgl32.Mat4 {
    return mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.1, 10.0)
}

func NewCamera(config *conf.CameraConfig, width,height float32) *Camera {
    ret := &Camera{Width: width, Height: height}
    ret.Isometric = config.Isometric
    
    ret.Configure(config)
    return ret
}

func (camera *Camera) Uniform(program *Program) {
//	
//	
//	camera.projectionUniform = UniformMatrix4fv(program, PROJECTION, 1, &camera.projection[0])
//	camera.cameraUniform = UniformMatrix4fv(program, VIEW, 1, &camera.camera[0])
	
	
	camera.projectionUniform = program.UniformMatrix4fv(PROJECTION, 1, &camera.projection[0] )
	camera.cameraUniform = program.UniformMatrix4fv(VIEW, 1, &camera.camera[0] )
    
}


func (camera *Camera) Configure(config *conf.CameraConfig) {
    camera.Isometric = config.Isometric
    if camera.Isometric {
        camera.projection = orthographic(camera.Width, camera.Height)
        camera.camera = mgl32.LookAtV(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
    } else {
        camera.projection = perspective(camera.Width, camera.Height)
        camera.camera = mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
    }
}


func (camera *Camera) Desc() string {
    ret := ""
    tmp := ""
    if camera.Isometric { tmp = " iso" }
    ret += fmt.Sprintf("camera[%.0fx%.0f%s]@%p",camera.Width,camera.Height,tmp,&camera.camera) 
    ret += fmt.Sprintf("\n%v",camera.projection)
    return ret
}