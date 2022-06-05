// +build RENDERER

package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Draft struct {
	config DraftConfig

	camera *gfx.Camera

	vert map[string]*gfx.Shader
	frag map[string]*gfx.Shader

	program map[string]*gfx.Program

	object map[string]uint32
}

func (draft *Draft) RenderAxis(debug bool) {
	program := draft.program["axis"]

	program.UseProgram(debug)

	object := draft.object["axis"]

	draft.camera.Uniform(program)

	const a = 1.0
	var axis []float32 = []float32{
		0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		a, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0,
		0.0, a, 0.0, 0.0, 1.0, 0.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, a, 0.0, 0.0, 1.0, 1.0,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, object)
	gl.BufferData(gl.ARRAY_BUFFER, len(axis)*4, gl.Ptr(axis), gl.STATIC_DRAW)

	program.VertexAttribPointer(gfx.VERTEX, 3, (3+4)*4, 0)
	program.VertexAttribPointer(gfx.COLOR, 4, (3+4)*4, 3*4)

	model := mgl32.Ident4()
	//	model = mgl32.Scale3D(0.25,0.25,0.25)

	program.UniformMatrix4fv(gfx.MODEL, 1, &model[0])

	gl.LineWidth(4.0)
	gl.BindBuffer(gl.ARRAY_BUFFER, object)
	gl.DrawArrays(gl.LINES, 0, 3*2)

}

func (draft *Draft) Render(camera *gfx.Camera, debug, verbose bool) {
	gl.ClearColor(.5, .5, .5, 1.)

	if true {
		draft.RenderAxis(debug)
	}

}

func (draft *Draft) Init(camera *gfx.Camera, font *gfx.Font, programService *gfx.ProgramService) {

	draft.camera = camera

	draft.program = map[string]*gfx.Program{}

	draft.object = map[string]uint32{}
	for _, name := range []string{"axis", "quad"} {
		tmp := uint32(0x0)
		gl.GenBuffers(1, &tmp)
		draft.object[name] = tmp
	}

	var err error
	{
		draft.program["axis"] = programService.GetProgram("axis", "")

		err = draft.program["axis"].Link("color", "color")
		if err != nil {
			log.Error("%s fail linking %s color shaders: %s", "axis", draft.Desc(), err)
		}

		//        draft.program["axis"] = gfx.NewProgram("axis","axis",programService)

		//        err = draft.program["axis"].GetCompileShaders("","color","color")
		//        if err != nil { log.Error("fail loading %s color shaders: %s","axis",err) }

		//        err = draft.program["axis"].LinkProgram()
		//        if err != nil { log.Error("fail linking %s color shaders: %s","axis",err) }
		//
		//        if draft.program["axis"] == nil {
		//            log.Error("fail Init!!")
		//        }

	}

}

func (draft *Draft) Configure(config *DraftConfig) {
	if config == nil {
		return
	}
}

func NewDraft(config *DraftConfig) *Draft {
	ret := &Draft{}
	return ret
}

func (draft *Draft) Desc() string { return draft.config.Desc() }
func (draft *Draft) Dump() string { return draft.config.Desc() }
