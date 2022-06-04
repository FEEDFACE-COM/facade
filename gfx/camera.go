// +build darwin,amd64 darwin,arm64

package gfx

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	projection        mgl32.Mat4
	projectionUniform int32

	view        mgl32.Mat4
	viewUniform int32

	size Size

	zoom      float32
	isometric bool
}

func orthographic(width, height float32) mgl32.Mat4 {
	const F = 2.
	var w float32 = F * (width / height) / 2.
	var h float32 = F * (height / height) / 2.
	var l float32 = F * 1000. / 2.
	return mgl32.Ortho(-w, w, -h, h, -l, l)
}

func perspective(width, height float32) mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.01, 10.0)
}

func (camera *Camera) Ratio() float32 { return camera.size.W / camera.size.H }

func NewCamera(zoom float32, iso bool, screen Size) *Camera {
	ret := &Camera{zoom: zoom, isometric: iso, size: screen}
	return ret
}

func (camera *Camera) Projection() *float32 { return &camera.projection[0] }
func (camera *Camera) View() *float32       { return &camera.view[0] }

//rem remove:
func (camera *Camera) Uniform(program *Program) {
	camera.projectionUniform, _ = program.UniformMatrix4fv(PROJECTION, 1, &camera.projection[0])
	camera.viewUniform, _ = program.UniformMatrix4fv(VIEW, 1, &camera.view[0])
}

func (camera *Camera) Init() {
	camera.calculate()
}

func (camera *Camera) ConfigureZoom(zoom float32) {
	if camera.zoom != zoom {
		camera.zoom = zoom
		camera.calculate()
	}
}

func (camera *Camera) ConfigureIsometric(iso bool) {
	if camera.isometric != iso {
		camera.isometric = iso
		camera.calculate()
	}
}

func (camera *Camera) calculate() {
	//    const MAGIC = 2.5 / 1.05
	const MAGIC = 2.5
	position := mgl32.Vec3{0, 0, MAGIC}
	camera.view = mgl32.Ident4()
	zoom := camera.zoom
	if camera.isometric {
		camera.projection = orthographic(camera.size.W, camera.size.H)
		camera.view = camera.view.Mul4(mgl32.LookAtV(mgl32.Vec3{1, 1, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}))
		camera.view = camera.view.Mul4(mgl32.Scale3D(zoom, zoom, zoom))
	} else {
		camera.projection = perspective(camera.size.W, camera.size.H)
		camera.view = camera.view.Mul4(mgl32.LookAtV(position, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}))
		camera.view = camera.view.Mul4(mgl32.Scale3D(zoom, zoom, zoom))
	}

}

func (camera *Camera) Desc() string {
	ret := "camera["
	ret += fmt.Sprintf("%.1f", camera.zoom)
	if camera.isometric {
		ret += "i"
	} else {
		ret += "p"
	}
	ret += "]"
	return ret
}
