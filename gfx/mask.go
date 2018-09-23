
// +build linux,arm

package gfx

import(
    conf "../conf"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Mask struct {

    Mask bool

    program *Program

    object *Object
    data []float32
    
    Width float32
    Height float32
    
}

func NewMask(config *conf.MaskConfig, width, height float32) *Mask {
    ret := &Mask{Width: width, Height: height}
    ret.Configure(config)
    return ret
}


func (mask *Mask) Configure(config *conf.MaskConfig) {
    mask.Mask = config.Mask
}

func (mask *Mask) Desc() string { return "mask[]" }


func (mask *Mask) Render() {

    if !mask.Mask {
        return
    }

    mask.program.UseProgram()
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

    vert := NewShader("vert",VertexShader["mask"],gl.VERTEX_SHADER)
    if err = vert.CompileShader(); err != nil { log.Error("fail compile mask vertex shader: %s",err) }
    frag := NewShader("frag",FragmentShader["mask"],gl.FRAGMENT_SHADER)
    if err = frag.CompileShader(); err != nil { log.Error("fail compile mask frag shader: %s",err) }
    
    mask.program = NewProgram("mask")
    if err = mask.program.CreateProgram(vert,frag); err != nil { log.Error("fail to create program: %s",err) }
    
            

}

