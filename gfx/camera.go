
// +build linux,arm

package gfx

import (
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
    ret.projection = mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.1, 10.0)
    ret.camera = mgl32.LookAtV(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
    return ret
}

func (camera *Camera) Uniform(program uint32) {
	camera.projectionUniform = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(camera.projectionUniform, 1, false, &camera.projection[0])

	camera.cameraUniform = gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(camera.cameraUniform, 1, false, &camera.camera[0])
    
}