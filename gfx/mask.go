
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
    mask bool
}

func NewMask(config *MaskConfig, screen Size) *Mask {
    ret := &Mask{width: screen.W, height: screen.H}
    return ret
}


func (mask *Mask) Configure(config *MaskConfig) {
    if config == nil { return }
    
    // TODO: if no change, return

    if val,ok := config.Mask(); ok {
        mask.mask = val
        log.Debug("config %s",mask.Desc())
    }    
}

func (mask *Mask) Mask() bool { return mask.mask }

func (mask *Mask) Desc() string {
    ret := "mask["
    if mask.mask {
        ret += "✓"
    }
    ret += "]"
    return ret
}


func (mask *Mask) Render(debug bool) {

    if !mask.mask {
        return
    }

    mask.program.UseProgram(debug)
    mask.object.BindBuffer()
    
    
    mask.program.Uniform1f(SCREENRATIO, mask.width / mask.height)
    mask.program.VertexAttribPointer(VERTEX, 3, (3+2)*4, 0 )
    mask.program.VertexAttribPointer(TEXCOORD, 2, (3+2)*4, 3*4)

    gl.DrawArrays(gl.TRIANGLES, 0, 3*2 )

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

    var err error    
    mask.program = GetProgram("mask")

    err = mask.program.GetCompileShaders("mask/","null","null")
    if err != nil { log.Error("fail load mask shaders: %s",err) }

    err = mask.program.LinkProgram(); 
    if err != nil { log.Error("fail to link mask program: %s",err) }

}

