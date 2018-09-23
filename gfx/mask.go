
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

    vert := NewShader("vert",vertexShader,gl.VERTEX_SHADER)
    if err = vert.CompileShader(); err != nil { log.Error("fail compile mask vertex shader: %s",err) }
    frag := NewShader("frag",fragmentShader,gl.FRAGMENT_SHADER)
    if err = frag.CompileShader(); err != nil { log.Error("fail compile mask frag shader: %s",err) }
    
    mask.program = NewProgram("mask")
    if err = mask.program.CreateProgram(vert,frag); err != nil { log.Error("fail to create program: %s",err) }
    
            

}

const vertexShader = `
attribute vec2 texcoord;
attribute vec3 vertex;
attribute vec4 color;

varying vec4 fragcolor;
varying vec2 fragcoord;

void main() {
    fragcolor = vec4( vertex, 1.0);
    fragcoord = texcoord;
    gl_Position = vec4(vertex,1);
}
`


const fragmentShader = `
varying vec4 fragcolor;
varying vec2 fragcoord;

float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.5) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}

void main() {
    vec4 col = vec4(0.0,0.0,0.0,0.0);
    vec2 pos = fragcoord;
    
    if ( grid(pos) ) { col = vec4(1.,1.,1.,0.5); }
    
//    if ( pos.y > 0.0 && pos.y < 1.0 && abs(pos.x) <= w ) { col = vec4(0.,1.,0.,1.); }
//    if ( pos.x > 0.0 && pos.x < 1.0 && abs(pos.y) <= w ) { col = vec4(1.,0.,0.,1.); }

    gl_FragColor = col;

}
`



