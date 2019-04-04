
// +build linux,arm

package gfx

import(
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Mask struct {

    program *Program

    object *Object
    data []float32
    
    width float32
    height float32
    
    state MaskState
}

func NewMask(config *MaskConfig, screen Size) *Mask {
    ret := &Mask{width: screen.W, height: screen.H}
    ret.state.ApplyConfig(config)
    return ret
}


func (mask *Mask) Configure(config *MaskConfig) {
    if config == nil { return }

	if mask.state.ApplyConfig(config) {
        log.Debug("mask config %s",mask.Desc())
        mask.LoadShaders()	
    }
	
}


func (mask *Mask) Desc() string { return mask.state.Desc() }

func (mask *Mask) Render(debug bool) {

    if mask.state.Mask == "def" {
        return
    }

    mask.program.UseProgram(debug)
    mask.object.BindBuffer()
    
    
    mask.program.Uniform1f(SCREENRATIO, mask.width / mask.height)
    mask.program.VertexAttribPointer(VERTEX, 3, (3+2)*4, 0 )
    mask.program.VertexAttribPointer(TEXCOORD, 2, (3+2)*4, 3*4)

    gl.DrawArrays(gl.TRIANGLES, 0, 3*2 )

}

func (mask *Mask) LoadShaders() error {
    var err error
    err = mask.program.GetCompileShaders("mask/","def",mask.state.Mask)
    if err != nil { log.Error("fail load mask shaders: %s",err) }
    err = mask.program.LinkProgram(); 
    if err != nil { log.Error("fail to link mask program: %s",err) }
    return nil
}

func (mask *Mask) Init() {
    w := mask.width  
    h := mask.height 

    v := h/h * h/2. 
    u := w/h * w/2. 
    mask.data = []float32{
        // x     //y          // tx //ty
        -w/2.,  h/2., 0.0,    -u,  v,
        -w/2., -h/2., 0.0,    -u, -v,
         w/2., -h/2., 0.0,     u, -v,
         w/2., -h/2., 0.0,     u, -v,
         w/2.,  h/2., 0.0,     u,  v,
        -w/2.,  h/2., 0.0,    -u,  v,
    }

    mask.object = NewObject("mask")
    mask.object.Init()
    mask.object.BufferData(len(mask.data)*4, mask.data)

    mask.program = GetProgram("mask")
    mask.LoadShaders()


}

