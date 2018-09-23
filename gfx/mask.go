
// +build linux,arm

package gfx

import(
    conf "../conf"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Mask struct {

    Mask bool

    program uint32
    object uint32

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

    gl.UseProgram(mask.program)
    gl.BindBuffer(gl.ARRAY_BUFFER,mask.object)
    
    vertAttrib := uint32(gl.GetAttribLocation(mask.program, gl.Str("vert\x00")))
    gl.EnableVertexAttribArray(vertAttrib) 
    gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, (3+2)*4, gl.PtrOffset(0))

    texCoordAttrib := uint32(gl.GetAttribLocation(mask.program, gl.Str("vertTexCoord\x00")))
    gl.EnableVertexAttribArray(texCoordAttrib) 
    gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, (3+2)*4, gl.PtrOffset(3*4))

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

    gl.GenBuffers(1,&mask.object)
    gl.BindBuffer(gl.ARRAY_BUFFER,mask.object)
    gl.BufferData(gl.ARRAY_BUFFER, len(mask.data)*4, gl.Ptr(mask.data), gl.STATIC_DRAW )
    
    var err error
    vert := NewShader("vert",vertexShader,gl.VERTEX_SHADER)
    if err = vert.Compile(); err != nil { log.Error("fail compile mask vertex shader: %s",err) }
    
    frag := NewShader("frag",fragmentShader,gl.FRAGMENT_SHADER)
    if err = frag.Compile(); err != nil { log.Error("fail compile mask frag shader: %s",err) }
    
    mask.program, err = CreateProgram(vert,frag)
    if err != nil { log.Error("fail to creat program: %s",err) }
    
            

}

const vertexShader = `
attribute vec2 vertTexCoord;
attribute vec3 vert;
attribute vec4 color;

varying vec4 vertColor;
varying vec2 fragTexCoord;

void main() {
    vertColor = vec4( vert, 1.0);
    fragTexCoord = vertTexCoord;
    gl_Position = vec4(vert,1);
}
`


const fragmentShader = `
varying vec4 vertColor;
varying vec2 fragTexCoord;

float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.25) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}

void main() {
    vec4 col = vec4(0.0,0.0,0.0,0.0);
    vec2 pos = fragTexCoord;
    
    if ( grid(pos) ) { col = vec4(1.,1.,1.,0.5); }
    
//    if ( pos.y > 0.0 && pos.y < 1.0 && abs(pos.x) <= w ) { col = vec4(0.,1.,0.,1.); }
//    if ( pos.x > 0.0 && pos.x < 1.0 && abs(pos.y) <= w ) { col = vec4(1.,0.,0.,1.); }

    gl_FragColor = col;

}
`

const foo = `
bool test(vec2 pos) {
    return false;
}


varying vec4 vertColor;
varying vec2 fragTexCoord;
void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    vec2 pos = fragTexCoord;
    if (pos.x >= 0.0 && pos.y >= 0.0) { col.g = 1.0; col.b = 1.0; }
    if (pos.x >= 0.0 && pos.y <= 0.0) { col.g = 1.0;              }
    if (pos.x <= 0.0 && pos.y >= 0.0) {              col.b = 1.0; }
    if (pos.x <= 0.0 && pos.y <= 0.0) { col.g = 0.5; col.b = 0.5; }
    float a = 0.5;
    if ( abs(pos.x) < a && abs(pos.y) < a ) { col.r = 1.0; }
    gl_FragColor = vec4(col,1.0);
}
`


