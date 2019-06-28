
// +build linux,arm

package gfx

import(
    "github.com/go-gl/mathgl/mgl32"    
    log "../log"
    gl "src.feedface.com/gfx/piglet/gles2"
    
)

type Axis struct {
    
    program *Program
    object *Object
    
    data []float32
}


func NewAxis() *Axis { return &Axis{} }


func (axis *Axis) Render(camera *Camera, debug bool) {

    axis.program.UseProgram(debug)
    axis.object.BindBuffer()
    camera.Uniform(axis.program)
    
    axis.program.VertexAttribPointer(VERTEX, 3, (3+4)*4, 0 )
    axis.program.VertexAttribPointer(COLOR,  4, (3+4)*4, 3*4 )

    model := mgl32.Ident4()
    axis.program.UniformMatrix4fv(MODEL,1,&model[0])

    gl.LineWidth(4.0)    
    gl.DrawArrays(gl.LINES, 0, 3*2 )
    
    
}
    

func (axis *Axis) Init(shaderService *ShaderService) {
    
    const a = 1.0
    axis.data = []float32{ 
        0.0, 0.0, 0.0,     1.0, 0.0, 0.0, 1.0,
          a, 0.0, 0.0,     1.0, 0.0, 0.0, 1.0,
        0.0, 0.0, 0.0,     0.0, 1.0, 0.0, 1.0,
        0.0,   a, 0.0,     0.0, 1.0, 0.0, 1.0,
        0.0, 0.0, 0.0,     0.0, 0.0, 1.0, 1.0,
        0.0, 0.0,   a,     0.0, 0.0, 1.0, 1.0,  
    }
    

    axis.object = NewObject("axis")
    axis.object.Init()
    axis.object.BufferData(len(axis.data)*4, axis.data)

    var err error    
    axis.program = NewProgram("axis",shaderService)

    err = axis.program.GetCompileShaders("","color","color")
    if err != nil { 
        log.Error("%s fail load axis shaders: %s",axis.Desc(),err) 
        return
    }
    
    err = axis.program.LinkProgram(); 
    if err != nil { log.Error("%s fail link axis program: %s",axis.Desc(),err) }

//    log.Debug("%s init",axis.Desc())
}

func (axis *Axis) Desc() string { return "axis[]" }
