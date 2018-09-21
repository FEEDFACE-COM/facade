
// +build linux,arm

package gfx

import (
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"    
)

type Camera struct {
    projection mgl32.Mat4
    projectionUniform int32
    
    camera mgl32.Mat4
    cameraUniform int32
    
        
}


func NewCamera(width,height float32) *Camera {
    ret := &Camera{}
    w := float32(width)
    h := float32(height)
    ratio := h/w
    
    l := float32(-1.)
    r := float32(1.)
    b := float32(-ratio)
    t := float32(ratio)
    n := float32(-4096.)
    f := float32(4096.)
    
    ortho := mgl32.Ortho( l,r,b,t,n,f )
    persp := mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.1, 10.0)
    ret.projection = ortho
    log.Debug("orth\n%v",ortho)
    log.Debug("persp\n%v",persp)
    log.Debug("now\n%v",ret.projection)
    
    
    
    
    ret.camera = mgl32.LookAtV(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
    return ret
}

func (camera *Camera) Uniform(program uint32) {
	camera.projectionUniform = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(camera.projectionUniform, 1, false, &camera.projection[0])

	camera.cameraUniform = gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(camera.cameraUniform, 1, false, &camera.camera[0])
    
}