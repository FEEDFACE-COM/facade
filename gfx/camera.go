
// +build linux,arm

package gfx

import (
    log "../log"
//    log "../log"
//    gl "src.feedface.com/gfx/piglet/gles2"
	"github.com/go-gl/mathgl/mgl32"    
)

type Camera struct {
    
    projection mgl32.Mat4
    projectionUniform int32
    
    view mgl32.Mat4
    viewUniform int32

    size Size    
    
    state CameraState
}


func orthographic(width,height float32) mgl32.Mat4 {
    const F = 2.
    var w float32 = F * (width/height) /2.
    var h float32 = F * (height/height) /2.
    var l float32 = F * 1000. /2.
    return mgl32.Ortho( -w, w, -h, h, -l, l )
}

func perspective(width,height float32) mgl32.Mat4 {
    return mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.01, 10.0)
}

func (camera *Camera) Ratio() float32 { return camera.size.W / camera.size.H }
    


func NewCamera(config *CameraConfig, screen Size) *Camera {
    ret := &Camera{size: screen}
    ret.state.ApplyConfig(config)
    return ret
}
    
func (camera *Camera) Uniform(program *Program) {
	camera.projectionUniform,_ = program.UniformMatrix4fv(PROJECTION, 1, &camera.projection[0] )
	camera.viewUniform,_ = program.UniformMatrix4fv(VIEW, 1, &camera.view[0] )
}


func (camera *Camera) Init(config *CameraConfig) {
	camera.Configure(config)
} 


func (camera *Camera) Configure(config *CameraConfig) {
    if config == nil { return }

	if camera.state.ApplyConfig(config) {
		log.Debug("camera config %s",camera.Desc())
	}
	
	
    
//    const MAGIC = 2.5 / 1.05
    const MAGIC = 2.5
    position := mgl32.Vec3{0,0,MAGIC}
    camera.view = mgl32.Ident4()
    zoom := float32(camera.state.Zoom)
    if camera.state.Isometric {
        camera.projection = orthographic(camera.size.W, camera.size.H)
        camera.view = camera.view.Mul4( mgl32.LookAtV(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}) )
        camera.view = camera.view.Mul4( mgl32.Scale3D( zoom, zoom, zoom ) )
    } else {
        camera.projection = perspective(camera.size.W, camera.size.H)
        camera.view = camera.view.Mul4( mgl32.LookAtV(position, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}) )
        camera.view = camera.view.Mul4( mgl32.Scale3D( zoom, zoom, zoom ) )
    }
    log.Debug("cam config %s",camera.Desc())

}


func (camera *Camera) Desc() string { return camera.state.Desc() }
    

