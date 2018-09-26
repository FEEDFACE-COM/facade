
package modes

import(
	"github.com/go-gl/mathgl/mgl32"    
    conf "../conf"
    gfx "../gfx"
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
)


type Test struct {
    
    config conf.TestConfig
    
    camera *gfx.Camera

    vert map[string]*gfx.Shader
    frag map[string]*gfx.Shader

    program map[string]*gfx.Program

    object map[string]uint32

}



func (test *Test) RenderAxis(debug bool) {
    program := test.program["axis"]
  
    program.UseProgram(debug)

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
    program.VertexAttribPointer(gfx.COLOR,  4, (3+4)*4, 3*4 )
    
    
    model := mgl32.Ident4()
    //	model = mgl32.Scale3D(0.25,0.25,0.25)
    
    program.UniformMatrix4fv(gfx.MODEL,1,&model[0])

  
    gl.LineWidth(4.0)    
    gl.BindBuffer(gl.ARRAY_BUFFER,object) 
    gl.DrawArrays(gl.LINES, 0, 3*2 )
    
    
}





func (test *Test) Render(camera *gfx.Camera, debug, verbose bool) {
    gl.ClearColor(.5,.5,.5,1.)
    
    if true { test.RenderAxis(debug) }
        
}


func (test *Test) Init(camera *gfx.Camera, font *gfx.Font) {

    test.camera = camera

    
    test.program = map[string]*gfx.Program{}
    
    test.object = map[string]uint32{}
    for _,name := range []string{ "axis", "quad" } {
        tmp := uint32(0x0)
        gl.GenBuffers(1,&tmp)
        test.object[name] = tmp
    }
    
    var err error
    {
        test.program["axis"] = gfx.NewProgram("axis")
    
        err = test.program["axis"].LoadShaders("color", "color")
        if err != nil { log.Error("fail loading %s color shaders: %s","axis",err) }

        err = test.program["axis"].LinkProgram()
        if err != nil { log.Error("fail linking %s color shaders: %s","axis",err) }
        
        if test.program["axis"] == nil {
            log.Error("fail Init!!")    
        }
            
    }
    
}



func (test *Test) Queue(text string) {
//    log.Debug("test %s",text);    
}





func (test *Test) Configure(config *conf.TestConfig) {
    if config == nil { return }
}

func NewTest(config *conf.TestConfig) *Test {
    ret := &Test{}
    return ret
}


func (test *Test) Desc() string { return test.config.Desc() }
func (test *Test) Dump() string { return test.config.Desc() }

