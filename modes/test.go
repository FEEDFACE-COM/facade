
package modes

import(
	"github.com/go-gl/mathgl/mgl32"    
    conf "../conf"
    gfx "../gfx"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
)


type Test struct {
    
    camera *gfx.Camera

    vert map[string]*gfx.Shader
    frag map[string]*gfx.Shader

    program map[string]*gfx.Program

    object map[string]uint32

}



var vertexShader = map[string]string{

"ident":`
uniform mat4 projection;
uniform mat4 camera;
attribute vec3 vert;
attribute vec4 color;
varying vec4 vertColor;

void main() {
    vertColor = color;
    gl_Position = projection * camera * vec4(vert,1);
}
`,




}

var fragmentShader = map[string]string{

"ident":`
varying vec4 vertColor;
void main() {
    gl_FragColor = vertColor;
}
`,



}


func (test *Test) RenderAxis() {
    program := test.program["test"]
  
  
    program.Use()

    object := test.object["axis"]

    test.camera.Uniform(program)    
  
    const a = 1.0
    var axis []float32 = []float32{ 
        0.0, 0.0, 0.0,     1.0, 0.0, 0.0, 1.0,
          a, 0.0, 0.0,     1.0, 0.0, 0.0, 1.0,
        0.0, 0.0, 0.0,     0.0, 1.0, 0.0, 1.0,
        0.0,   a, 0.0,     0.0, 1.0, 0.0, 1.0,
        0.0, 0.0, 0.0,     0.0, 0.0, 1.0, 1.0,
        0.0, 0.0,   a,     0.0, 0.0, 1.0, 1.0,  
    }
    gl.BindBuffer(gl.ARRAY_BUFFER,object) 
    gl.BufferData(gl.ARRAY_BUFFER, len(axis)*4, gl.Ptr(axis), gl.STATIC_DRAW)
    
    
    program.VertexAttribPointer(gfx.VERTEX, 3, (3+4)*4, 0 )
    program.VertexAttribPointer(gfx.COLOR, 4, (3+4)*4, 3*4 )
    
    
    model := mgl32.Ident4()
    //	model = mgl32.Scale3D(0.25,0.25,0.25)
    
    program.UniformMatrix4fv(gfx.MODEL,1,&model[0])
  
    gl.LineWidth(4.0)    
    gl.BindBuffer(gl.ARRAY_BUFFER,object) 
    gl.DrawArrays(gl.LINES, 0, 3*2 )
    
    
}





func (test *Test) Render() {
    gl.ClearColor(.5,.5,.5,1.)
    
    if true { test.RenderAxis() }
        
}


func (test *Test) Init(camera *gfx.Camera) {

    test.camera = camera

    
    test.program = map[string]*gfx.Program{}
    
    test.object = map[string]uint32{}
    for _,name := range []string{ "axis", "quad" } {
        tmp := uint32(0x0)
        gl.GenBuffers(1,&tmp)
        test.object[name] = tmp
    }
    
    test.vert =  map[string]*gfx.Shader{}
    for name,src := range vertexShader {
        test.vert[name] = gfx.NewShader(name,src,gl.VERTEX_SHADER)
        if err := test.vert[name].Compile() ; err != nil {
            log.Error("fail compile vertex shader %s: %s",name,err)
        }
    }
    
    test.frag =  map[string]*gfx.Shader{}
    for name,src := range fragmentShader {
        test.frag[name] = gfx.NewShader(name,src,gl.FRAGMENT_SHADER)
        if err := test.frag[name].Compile() ; err != nil {
            log.Error("fail compile fragment shader %s: %s",name,err)
        }
    }
    
    var err error
    test.program["test"] = gfx.NewProgram("test");
    err = test.program["test"].Create(test.vert["ident"],test.frag["ident"])
    if err != nil { log.Error("fail to create test: %s",err) }

    
    
}



func (test *Test) Queue(text string) {
    log.Debug("test %s",text);    
}





func (test *Test) Configure(config *conf.TestConfig) {}

func NewTest(config *conf.TestConfig) *Test {
    ret := &Test{}
    return ret
}


func (test *Test) Desc() string { return "test[]" }
func (test *Test) Dump() string { return "test[]" }

