
// +build linux,arm

package gfx

import(
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Mask struct {

    config MaskConfig
    program *Program

    object *Object
    data []float32
    
    Width float32
    Height float32
    
}

func NewMask(config *MaskConfig, screen Size) *Mask {
    ret := &Mask{config: *config, Width: screen.W, Height: screen.H}
    return ret
}


func (mask *Mask) Configure(config *MaskConfig) {
    if config == nil { return }
    if *config == mask.config { return }
    
    log.Debug("config %s -> %s",mask.Desc(),config.Desc())
    mask.config = *config
    
    
}

func (mask *Mask) Desc() string { return mask.config.Desc() }


func (mask *Mask) Render() {

    if !mask.config.Mask {
        return
    }

    mask.program.UseProgram(false)
    mask.object.BindBuffer()
    
    mask.program.VertexAttribPointer(VERTEX, 3, (3+2)*4, 0 )
    mask.program.VertexAttribPointer(TEXCOORD, 2, (3+2)*4, 3*4)

    gl.DrawArrays(gl.TRIANGLES, 0, 3*2 )

}


func (mask *Mask) Init() {
    w := mask.Width  
    h := mask.Height 

    v := h/h * h/2. 
    u := w/h * w/2. 
    mask.data = []float32{
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
    mask.program = NewProgram("mask")

    err = mask.program.LoadShaders("mask","mask")
    if err != nil { log.Error("fail load mask shaders: %s",err) }

    err = mask.program.LinkProgram(); 
    if err != nil { log.Error("fail to link mask program: %s",err) }

}

