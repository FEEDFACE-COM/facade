//go:build (linux && arm) || DARWIN_GUI
// +build linux,arm DARWIN_GUI

package gfx

import (
	"FEEDFACE.COM/facade/log"
	"fmt"
	gl "github.com/FEEDFACE-COM/piglet/gles2"


const DEBUG_MASK = false

type Mask struct {
	program *Program

	object *Object
	data   []float32

	width  float32
	height float32

	name string
}

func NewMask(name string, screen Size) *Mask {
	ret := &Mask{width: screen.W, height: screen.H}
	ret.name = name
	return ret
}

func (mask *Mask) ConfigureName(name string) {
	if DEBUG_MASK {
		log.Debug("%s configure name %s", mask.Desc(), name)
	}
	if mask.name != name {
		err := mask.program.Link("def", name)
		if err != nil {
			if DEBUG_MASK {
				log.Debug("%s fail to load mask shader %s: %s", mask.Desc(), name, err)
			}
			//            err = mask.program.LoadShaders( "def", mask.name )
		} else {
			mask.name = name
		}
	}
}

func (mask *Mask) Desc() string {
	return fmt.Sprintf("mask[%s]", mask.name)
}

func (mask *Mask) Render(debug bool) {

	if mask.name == "def" {
		return
	}

	mask.program.UseProgram(debug)
	mask.object.BindBuffer()

	mask.program.Uniform1f(SCREENRATIO, mask.width/mask.height)
	mask.program.VertexAttribPointer(VERTEX, 3, (3+2)*4, 0)
	mask.program.VertexAttribPointer(TEXCOORD, 2, (3+2)*4, 3*4)

	gl.DrawArrays(gl.TRIANGLES, 0, 3*2)

}

//func (mask *Mask) LoadShaders() error {
//    var err error
//
//    err = mask.program.GetCompileShaders("mask/","def",mask.name)
//    if err != nil {
//        return log.NewError("fail load mask shaders: %s",err)
//    }
//    err = mask.program.LinkProgram();
//    if err != nil {
//        return log.NewError("fail to link mask program: %s",err)
//    }
//    return nil
//}

func (mask *Mask) Init(programService *ProgramService) {
	w := mask.width
	h := mask.height

	v := h / h * h / 2.
	u := w / h * w / 2.
	mask.data = []float32{
		// x     //y          // tx //ty
		-w / 2., h / 2., 0.0, -u, v,
		-w / 2., -h / 2., 0.0, -u, -v,
		w / 2., -h / 2., 0.0, u, -v,
		w / 2., -h / 2., 0.0, u, -v,
		w / 2., h / 2., 0.0, u, v,
		-w / 2., h / 2., 0.0, -u, v,
	}

	mask.object = NewObject("mask")
	mask.object.Init()
	mask.object.BufferData(len(mask.data)*4, mask.data)

	mask.program = programService.GetProgram("mask", "mask/")
	mask.program.Link("def", mask.name)

	if DEBUG_MASK {
		log.Debug("%s init", mask.Desc())
	}

}
